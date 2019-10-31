package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	helper "github.com/cascades-fbp/cascades-mqtt/lib"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	zmq "github.com/pebbe/zmq4"
)

var (
	// Flags
	inputEndpoint   = flag.String("port.in", "", "Component's input port endpoint")
	optionsEndpoint = flag.String("port.options", "", "Component's input port endpoint")
	errorEndpoint   = flag.String("port.err", "", "Component's output port endpoint")
	jsonFlag        = flag.Bool("json", false, "Print component documentation in JSON")
	debug           = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	inPort, optionsPort, errPort *zmq.Socket
	inCh, errCh                  chan bool
	exitCh                       chan os.Signal
	err                          error
)

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
	errCh = make(chan bool)
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

	ports := 1
	if errPort != nil {
		ports++
	}

	waitCh := make(chan bool)
	go func(num int) {
		total := 0
		for {
			select {
			case v := <-inCh:
				if !v {
					log.Println("IN port is closed. Interrupting execution")
					exitCh <- syscall.SIGTERM
					break
				} else {
					total++
				}
			case v := <-errCh:
				if !v {
					log.Println("ERR port is closed. Interrupting execution")
					exitCh <- syscall.SIGTERM
					break
				} else {
					total++
				}
			}
			if total >= num && waitCh != nil {
				waitCh <- true
			}
		}
	}(ports)

	log.Println("Waiting for port connections to establish... ")
	select {
	case <-waitCh:
		log.Println("Ports connected")
		waitCh = nil
	case <-time.Tick(30 * time.Second):
		log.Println("Timeout: port connections were not established within provided interval")
		exitCh <- syscall.SIGTERM
		return
	}

	log.Println("Waiting for options to arrive...")
	var (
		ip            [][]byte
		clientOptions *mqtt.ClientOptions
		client        *mqtt.MqttClient
		defaultTopic  string
		qos           mqtt.QoS
	)
	for {
		ip, err = optionsPort.RecvMessageBytes(0)
		if err != nil {
			continue
		}
		if !runtime.IsValidIP(ip) || !runtime.IsPacket(ip) {
			continue
		}

		clientOptions, defaultTopic, qos, err = helper.ParseOptionsURI(string(ip[1]))
		if err != nil {
			log.Printf("Failed to parse connection uri. Error: %s", err.Error())
			continue
		}

		client = mqtt.NewClient(clientOptions)
		if _, err = client.Start(); err != nil {
			log.Printf("Failed to create MQTT client. Error: %s", err.Error())
			continue
		}
		break
	}
	defer client.Disconnect(1e6)
	optionsPort.Close()

	log.Println("Started...")
	var customTopic string
	for {
		ip, err = inPort.RecvMessageBytes(0)
		if err != nil {
			log.Println("Error receiving message:", err.Error())
			continue
		}
		if !runtime.IsValidIP(ip) {
			continue
		}
		if runtime.IsOpenBracket(ip) {
			for {
				ip, err = inPort.RecvMessageBytes(0)
				if err != nil {
					log.Println("Error receiving message:", err.Error())
					continue
				}
				if !runtime.IsValidIP(ip) {
					continue
				}
				if runtime.IsCloseBracket(ip) {
					customTopic = ""
					break
				}
				if customTopic == "" {
					customTopic = string(ip[1])
					if customTopic == "" {
						customTopic = defaultTopic
					}
					continue
				}
				client.Publish(qos, customTopic, ip[1])
			}
			continue
		}
		client.Publish(qos, defaultTopic, ip[1])
	}
}

func validateArgs() {
	if *optionsEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *inputEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func openPorts() {
	optionsPort, err = utils.CreateInputPort("mqtt/pub.options", *optionsEndpoint, nil)
	utils.AssertError(err)

	inPort, err = utils.CreateInputPort("mqtt/pub.in", *inputEndpoint, inCh)
	utils.AssertError(err)

	if *errorEndpoint != "" {
		errPort, err = utils.CreateOutputPort("mqtt/pub.err", *errorEndpoint, errCh)
		utils.AssertError(err)
	}
}

func closePorts() {
	log.Println("Closing ports...")
	optionsPort.Close()
	inPort.Close()
	if errPort != nil {
		errPort.Close()
	}
	zmq.Term()
}
