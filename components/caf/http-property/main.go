package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	caf "github.com/cascades-fbp/cascades-caf"
	httputils "github.com/cascades-fbp/cascades-http/utils"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	zmq "github.com/pebbe/zmq4"
)

var (
	intervalEndpoint = flag.String("port.int", "10s", "Component's input port endpoint")
	requestEndpoint  = flag.String("port.req", "", "Component's input port endpoint")
	templateEndpoint = flag.String("port.tmpl", "", "Component's input port endpoint")
	propertyEndpoint = flag.String("port.prop", "", "Component's output port endpoint")
	responseEndpoint = flag.String("port.resp", "", "Component's output port endpoint")
	bodyEndpoint     = flag.String("port.body", "", "Component's output port endpoint")
	errorEndpoint    = flag.String("port.err", "", "Component's error port endpoint")
	jsonFlag         = flag.Bool("json", false, "Print component documentation in JSON")
	debug            = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	intPort, reqPort, tmplPort            *zmq.Socket
	propPort, respPort, bodyPort, errPort *zmq.Socket
	outCh                                 chan bool
	exitCh                                chan os.Signal
	err                                   error
)

type requestIP struct {
	URL         string              `json:"url"`
	Method      string              `json:"method"`
	ContentType string              `json:"content-type"`
	Headers     map[string][]string `json:"headers"`
}

func assertError(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

func validateArgs() {
	if *requestEndpoint == "" || *templateEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *propertyEndpoint == "" && *responseEndpoint == "" && *bodyEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func openPorts() {
	var err error

	// Input sockets
	// Interval socket
	intPort, err = utils.CreateInputPort("http-property.int", *intervalEndpoint, nil)
	utils.AssertError(err)

	// Request socket
	reqPort, err = utils.CreateInputPort("http-property.req", *requestEndpoint, nil)
	utils.AssertError(err)

	// Property template socket
	tmplPort, err = utils.CreateInputPort("http-property.tmpl", *templateEndpoint, nil)
	utils.AssertError(err)

	// Output sockets
	// Property socket
	if *propertyEndpoint != "" {
		propPort, err = utils.CreateOutputPort("http-property.prop", *propertyEndpoint, outCh)
		utils.AssertError(err)
	}

	// Response socket
	if *responseEndpoint != "" {
		respPort, err = utils.CreateOutputPort("http-property.resp", *responseEndpoint, outCh)
		utils.AssertError(err)
	}

	// Response body socket
	if *bodyEndpoint != "" {
		bodyPort, err = utils.CreateOutputPort("http-property.body", *bodyEndpoint, outCh)
		utils.AssertError(err)
	}

	// Error socket
	if *errorEndpoint != "" {
		errPort, err = utils.CreateOutputPort("http-property.err", *errorEndpoint, outCh)
		utils.AssertError(err)

	}
}

func closePorts() {
	intPort.Close()
	reqPort.Close()
	tmplPort.Close()
	propPort.Close()
	respPort.Close()
	bodyPort.Close()
	errPort.Close()

	zmq.Term()
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
	outCh = make(chan bool)
	exitCh = make(chan os.Signal, 1)

	// Start the communication & processing logic
	go mainLoop()

	// Wait for the end...
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh

	log.Println("Done")
}

// mainLoop initiates all ports and handles the traffic
func mainLoop() {
	openPorts()
	defer closePorts()

	waitCh := make(chan bool)
	go func() {
		total := 0
		for {
			v := <-outCh
			if !v {
				log.Println("An OUT port is closed. Interrupting execution")
				exitCh <- syscall.SIGTERM
				break
			} else {
				total++
			}
			// At least one output ports are opened
			if total >= 1 && waitCh != nil {
				waitCh <- true
			}
		}
	}()

	log.Println("Waiting for port connections to establish... ")
	select {
	case <-waitCh:
		log.Println("Ports connected")
		waitCh = nil
	case <-time.Tick(30 * time.Second):
		log.Println("Timeout: port connections were not established within provided interval")
		os.Exit(1)
	}

	// Setup socket poll items
	poller := zmq.NewPoller()
	poller.Add(intPort, zmq.POLLIN)
	poller.Add(reqPort, zmq.POLLIN)
	poller.Add(tmplPort, zmq.POLLIN)

	// This is obviously dangerous but we need it to deal with our custom CA's
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 30 * time.Second

	var (
		interval     time.Duration
		ip           [][]byte
		request      *requestIP
		propTemplate *caf.PropertyTemplate
		httpRequest  *http.Request
	)

	for {
		sockets, err := poller.Poll(-1)
		if err != nil {
			log.Println("Error polling ports:", err.Error())
			continue
		}
		for _, socket := range sockets {
			if socket.Socket == nil {
				log.Println("ERROR: could not find socket in polling items array")
				continue
			}
			ip, err = socket.Socket.RecvMessageBytes(0)
			if err != nil {
				log.Println("Error receiving message:", err.Error())
				continue
			}
			if !runtime.IsValidIP(ip) || !runtime.IsPacket(ip) {
				log.Println("Invalid IP:", ip)
				continue
			}
			switch socket.Socket {
			case intPort:
				interval, err = time.ParseDuration(string(ip[1]))
				log.Println("Interval specified:", interval)
			case reqPort:
				err = json.Unmarshal(ip[1], &request)
				if err != nil {
					log.Println("ERROR: failed to unmarshal request:", err.Error())
					continue
				}
				log.Println("Request specified:", request)
			case tmplPort:
				err = json.Unmarshal(ip[1], &propTemplate)
				if err != nil {
					log.Println("ERROR: failed to unmarshal template:", err.Error())
					continue
				}
				log.Printf("Template specified: %+v", propTemplate)

			default:
				log.Println("ERROR: IP from unhandled socket received!")
				continue
			}
		}
		if interval > 0 && request != nil && propTemplate != nil {
			log.Println("Component configured. Moving on...")
			break
		}
	}

	log.Println("Started...")
	ticker := time.NewTicker(interval)
	for _ = range ticker.C {
		httpRequest, err = http.NewRequest(request.Method, request.URL, nil)
		utils.AssertError(err)

		// Set the accepted Content-Type
		if request.ContentType != "" {
			httpRequest.Header.Add("Content-Type", request.ContentType)
		}

		// Set any additional headers if provided
		for k, v := range request.Headers {
			httpRequest.Header.Add(k, v[0])
		}

		response, err := client.Do(httpRequest)
		if err != nil {
			log.Printf("ERROR performing HTTP %s %s: %s\n", request.Method, request.URL, err.Error())
			if errPort != nil {
				errPort.SendMessageDontwait(runtime.NewPacket([]byte(err.Error())))
			}
			continue
		}

		resp, err := httputils.Response2Response(response)
		if err != nil {
			log.Println("ERROR converting response to reply:", err.Error())
			if errPort != nil {
				errPort.SendMessageDontwait(runtime.NewPacket([]byte(err.Error())))
			}
			continue
		}

		// Property output socket
		if propPort != nil {
			var data interface{}
			if strings.HasSuffix(request.ContentType, "json") {
				err = json.Unmarshal(resp.Body, &data)
				if err != nil {
					log.Println("ERROR unmarshaling the JSON response:", err.Error())
					continue
				}
			} else {
				// TODO: support other content-types
				log.Printf("WARNING processing of %s is not supported", request.ContentType)
				continue
			}

			prop, err := propTemplate.Fill(data)
			if err != nil {
				log.Println("ERROR filling template with data: ", err.Error())
				continue
			}

			out, _ := json.Marshal(prop)
			propPort.SendMessage(runtime.NewPacket(out))
		}

		// Extra output sockets (e.g., for debugging)
		if respPort != nil {
			ip, err = httputils.Response2IP(resp)
			if err != nil {
				log.Println("ERROR converting reply to IP:", err.Error())
				if errPort != nil {
					errPort.SendMessageDontwait(runtime.NewPacket([]byte(err.Error())))
				}
			} else {
				respPort.SendMessage(ip)
			}
		}
		if bodyPort != nil {
			bodyPort.SendMessage(runtime.NewPacket(resp.Body))
		}
	}
}
