package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	caf "github.com/cascades-fbp/cascades-caf"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	uuid "github.com/nu7hatch/gouuid"
	zmq "github.com/pebbe/zmq4"
)

var (
	dataEndpoint     = flag.String("port.data", "", "Component's input (array-)port endpoint")
	templateEndpoint = flag.String("port.tmpl", "", "Component's input port endpoint")
	updatedEndpoint  = flag.String("port.update", "", "Component's output port endpoint")
	matchedEndpoint  = flag.String("port.match", "", "Component's output port endpoint")
	errorEndpoint    = flag.String("port.err", "", "Component's error port endpoint")
	jsonFlag         = flag.Bool("json", false, "Print component documentation in JSON")
	debug            = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	dataPortsArray                        []*zmq.Socket
	poller                                *zmq.Poller
	tmplPort, updPort, matchPort, errPort *zmq.Socket
	err                                   error
	inCh                                  chan bool
	outCh                                 chan bool
	exitCh                                chan os.Signal
	ctxTemplate                           caf.ContextTemplate
	contexts                              map[string]*caf.Context
)

type rawData struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Group   string        `json:"group"`
	Entries []interface{} `json"entries"` // Determines whether it is a Context or a Property
}

func validateArgs() {
	if *dataEndpoint == "" || *templateEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *updatedEndpoint == "" && *matchedEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func openPorts() {
	// Data
	dataports := strings.Split(*dataEndpoint, ",")
	if len(dataports) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	dataPortsArray = []*zmq.Socket{}
	poller = zmq.NewPoller()

	for i, endpoint := range dataports {
		endpoint = strings.TrimSpace(endpoint)
		log.Printf("Connecting DATA[%v]=%s", i, endpoint)
		port, err := utils.CreateInputPort(fmt.Sprintf("context.data-%s", i), endpoint, inCh)
		utils.AssertError(err)

		dataPortsArray = append(dataPortsArray, port)
		poller.Add(port, zmq.POLLIN)
	}

	// Template
	tmplPort, err = utils.CreateInputPort("context.tmpl", *templateEndpoint, nil)
	utils.AssertError(err)

	// Update
	if *updatedEndpoint != "" {
		updPort, err = utils.CreateOutputPort("context.update", *updatedEndpoint, outCh)
		utils.AssertError(err)
	}

	// Match
	if *matchedEndpoint != "" {
		matchPort, err = utils.CreateOutputPort("context.match", *matchedEndpoint, outCh)
		utils.AssertError(err)
	}

	// Error
	if *errorEndpoint != "" {
		errPort, err = utils.CreateOutputPort("context.err", *errorEndpoint, nil)
		utils.AssertError(err)
	}
}

func closePorts() {
	for _, port := range dataPortsArray {
		port.Close()
	}

	tmplPort.Close()
	updPort.Close()
	matchPort.Close()
	errPort.Close()

	zmq.Term()
}

func getPortIndex(port *zmq.Socket) (int, error) {
	i := -1
	for k, p := range dataPortsArray {
		if port == p {
			i = k
			break
		}
	}
	if i < 0 {
		return i, errors.New("Socket not found in the dataPortsArray")
	}
	return i, nil
}

func updateContext(i int, data []byte) (*caf.Context, error) {
	var msg rawData
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	ctx, ok := contexts[msg.Group]
	if !ok {
		log.Println("Received data with a new group:", msg.Group)
		u4, _ := uuid.NewV4()
		ctx = &caf.Context{
			ID:          u4.String(),
			Name:        ctxTemplate.Name,
			Group:       msg.Group,
			Description: ctxTemplate.Description,
			Condition:   ctxTemplate.Condition,
			Timestamp:   time.Now().Unix(),
			Entries:     make([]interface{}, len(dataPortsArray)),
		}
		contexts[msg.Group] = ctx
	} else {
		ctx.Timestamp = time.Now().Unix()
	}

	// Unmarshal JSON in the corresponding entries slice element
	if len(msg.Entries) > 0 {
		// Context
		var tmp caf.Context
		_ = json.Unmarshal(data, &tmp)
		ctx.Entries[i] = tmp
	} else {
		// Property
		var tmp caf.Property
		_ = json.Unmarshal(data, &tmp)
		ctx.Entries[i] = tmp
	}

	// Once determined, context never gets undetermined again
	if !ctx.Determined {
		allValues := true
		for _, v := range ctx.Entries {
			if v == nil {
				allValues = false
			}
		}
		if allValues {
			ctx.Determined = true
		}
	}
	return ctx, nil
}

func main() {
	flag.Parse()

	if *jsonFlag {
		doc, _ := registryEntry.JSON()
		fmt.Println(string(doc))
		os.Exit(0)
	}

	log.SetFlags(0)
	if *debug {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	validateArgs()

	// Communication channels
	inCh = make(chan bool)
	outCh = make(chan bool)
	exitCh = make(chan os.Signal, 1)

	// Start the communication & processing logic
	go mainLoop()

	// Wait for the end...
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh

	log.Println("Done")
}

func mainLoop() {
	openPorts()
	defer closePorts()

	waitCh := make(chan bool)
	go func() {
		totalIn := 0
		totalOut := 0
		for {
			select {
			case v := <-inCh:
				if !v {
					log.Println("IN port is closed. Interrupting execution")
					exitCh <- syscall.SIGTERM
					break
				} else {
					totalIn++
				}
			case v := <-outCh:
				if !v {
					log.Println("OUT port is closed. Interrupting execution")
					exitCh <- syscall.SIGTERM
					break
				} else {
					totalOut++
				}
			}

			if totalIn >= 1 && totalOut >= 1 && waitCh != nil {
				waitCh <- true
			}
		}
	}()

	log.Println("Waiting for port connections to establish... ")
	select {
	case <-waitCh:
		log.Println("One of the output ports is connected")
		waitCh = nil
	case <-time.Tick(30 * time.Second):
		log.Println("Timeout: port connections were not established within provided interval")
		os.Exit(1)
	}

	contexts = make(map[string]*caf.Context)
	var ip [][]byte
	for {
		ip, err = tmplPort.RecvMessageBytes(0)
		if err != nil {
			log.Println("Error receiving message:", err.Error())
			continue
		}

		if !runtime.IsValidIP(ip) || !runtime.IsPacket(ip) {
			log.Println("Received invalid IP")
			continue
		}

		err = json.Unmarshal(ip[1], &ctxTemplate)
		if err != nil {
			log.Println("ERROR: failed to unmarshal template:", err.Error())
			continue
		}
		log.Printf("Template specified: %+v", ctxTemplate)
		break
	}
	log.Println("Started...")

	for {
		results, err := poller.Poll(-1)
		if err != nil {
			log.Println("Error polling ports:", err.Error())
			continue
		}
		for _, r := range results {
			if r.Socket == nil {
				log.Println("ERROR: could not find socket in the polling results")
				continue
			}
			ip, err = r.Socket.RecvMessageBytes(0)
			if err != nil {
				log.Println("Error receiving message:", err.Error())
				continue
			}
			if !runtime.IsValidIP(ip) || !runtime.IsPacket(ip) {
				log.Println("Received invalid IP")
				continue
			}

			// Find the port index
			index, err := getPortIndex(r.Socket)
			if err != nil {
				log.Printf("Error processing IP: %v\n", err.Error())
				continue
			}

			// Update the corresponding context
			ctx, err := updateContext(index, ip[1])
			if err != nil {
				log.Printf("Error updating context: %v\n", err.Error())
				continue
			}

			if ctx.Determined {
				// Evaluate the context
				m1 := ctx.Matching
				m2, err := ctx.Evaluate()
				if err != nil {
					log.Printf("Error evaluating context %s: %v\n", ctx.ID, err.Error())
				} else {
					// Notify if matching status has changed
					if m1 != m2 {
						if matchPort != nil {
							out, _ := json.Marshal(ctx)
							matchPort.SendMessage(runtime.NewPacket(out))
						}
					}
				}
			}

			// Notify about the update (always)
			if updPort != nil {
				out, _ := json.Marshal(ctx)
				updPort.SendMessage(runtime.NewPacket(out))
			}
		}
	}
}
