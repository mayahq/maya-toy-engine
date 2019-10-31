package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	caf "github.com/cascades-fbp/cascades-caf"
	helper "github.com/cascades-fbp/cascades-mqtt/lib"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	zmq "github.com/pebbe/zmq4"
)

var (
	// Flags
	optionsEndpoint  = flag.String("port.options", "", "Component's input port endpoint")
	templateEndpoint = flag.String("port.tmpl", "", "Component's input port endpoint")
	propertyEndpoint = flag.String("port.prop", "", "Component's output port endpoint")
	errorEndpoint    = flag.String("port.err", "", "Component's output port endpoint")
	jsonFlag         = flag.Bool("json", false, "Print component documentation in JSON")
	debug            = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	optsPort, tmplPort, propPort, errPort *zmq.Socket
	propTemplate                          *caf.PropertyTemplate
	options                               optionsIP
	outCh                                 chan bool
	exitCh                                chan os.Signal
	err                                   error
)

type optionsIP struct {
	OptionsURI  string `json:"optionsURI"`
	ContentType string `json:"content-type"`
}

func validateArgs() {
	if *optionsEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *templateEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *propertyEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func openPorts() {
	optsPort, err = utils.CreateInputPort("mqtt-property.options", *optionsEndpoint, nil)
	utils.AssertError(err)

	tmplPort, err = utils.CreateInputPort("mqtt-property.template", *templateEndpoint, nil)
	utils.AssertError(err)

	propPort, err = utils.CreateOutputPort("mqtt-property.property", *propertyEndpoint, outCh)
	utils.AssertError(err)

	if *errorEndpoint != "" {
		errPort, err = utils.CreateOutputPort("mqtt-property.err", *errorEndpoint, nil)
		utils.AssertError(err)
	}
}

func closePorts() {
	optsPort.Close()
	tmplPort.Close()
	propPort.Close()
	if errPort != nil {
		errPort.Close()
	}
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

	log.Println("Waiting for options to arrive...")
	var (
		ip            [][]byte
		clientOptions *mqtt.ClientOptions
		client        *mqtt.MqttClient
		defaultTopic  string
		qos           mqtt.QoS
	)

	// Setup socket poll items
	poller := zmq.NewPoller()
	poller.Add(optsPort, zmq.POLLIN)
	poller.Add(tmplPort, zmq.POLLIN)

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
			case optsPort:
				err = json.Unmarshal(ip[1], &options)
				if err != nil {
					log.Println("ERROR: failed to unmarshal options:", err.Error())
					continue
				}
				clientOptions, defaultTopic, qos, err = helper.ParseOptionsURI(options.OptionsURI)
				if err != nil {
					log.Printf("Failed to parse connection uri. Error: %s", err.Error())
					continue
				}
				log.Println("Options specified:", options)
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

			if defaultTopic != "" && clientOptions != nil && propTemplate != nil {
				log.Println("Component configured. Moving on...")
				break
			}
		}

		client = mqtt.NewClient(clientOptions)
		if _, err = client.Start(); err != nil {
			log.Printf("Failed to create MQTT client. Error: %s", err.Error())
			continue
		}

		defer client.Disconnect(1e6)

		optsPort.Close()
		tmplPort.Close()
		break
	}

	log.Println("Started...")
	topicFilter, err := mqtt.NewTopicFilter(defaultTopic, byte(qos))
	utils.AssertError(err)
	_, err = client.StartSubscription(messageHandler, topicFilter)
	utils.AssertError(err)

	ticker := time.Tick(1 * time.Second)
	for _ = range ticker {
	}
}

func messageHandler(client *mqtt.MqttClient, message mqtt.Message) {
	var data interface{}
	if strings.HasSuffix(options.ContentType, "json") {
		err = json.Unmarshal(message.Payload(), &data)
		if err != nil {
			log.Println("ERROR unmarshaling the JSON message:", err.Error())
			return
		}
	} else {
		// TODO: support other content-types
		log.Printf("WARNING processing of %s is not supported", options.ContentType)
		return
	}
	prop, err := propTemplate.Fill(data)

	if err != nil {
		log.Println("ERROR filling template with data: ", err.Error())
		return
	}
	out, _ := json.Marshal(prop)
	propPort.SendMessage(runtime.NewPacket(out))
}
