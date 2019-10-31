package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	zmq "github.com/alecthomas/gozmq"
	httputils "github.com/cascades-fbp/cascades-http/utils"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
)

var (
	// Flags
	optionsEndpoint = flag.String("port.options", "", "Component's options port endpoint")
	inputEndpoint   = flag.String("port.in", "", "Component's input port endpoint")
	outputEndpoint  = flag.String("port.out", "", "Component's output port endpoint")
	jsonFlag        = flag.Bool("json", false, "Print component documentation in JSON")
	debug           = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	context                      *zmq.Context
	optionsPort, inPort, outPort *zmq.Socket
	err                          error
)

func validateArgs() {
	if *optionsEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *inputEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *outputEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func openPorts() {
	context, err = zmq.NewContext()
	utils.AssertError(err)

	optionsPort, err = utils.CreateInputPort(context, *optionsEndpoint)
	utils.AssertError(err)

	inPort, err = utils.CreateInputPort(context, *inputEndpoint)
	utils.AssertError(err)
}

func closePorts() {
	optionsPort.Close()
	if inPort != nil {
		inPort.Close()
	}
	if outPort != nil {
		outPort.Close()
	}
	context.Close()
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

	openPorts()
	defer closePorts()

	exitCh := utils.HandleInterruption()

	// Wait for the configuration on the options port
	var bindAddr string
	for {
		log.Println("Waiting for configuration...")
		ip, err := optionsPort.RecvMultipart(0)
		if err != nil {
			log.Println("Error receiving IP:", err.Error())
			continue
		}
		if !runtime.IsValidIP(ip) || !runtime.IsPacket(ip) {
			continue
		}
		bindAddr = string(ip[1])
		break
	}
	optionsPort.Close()

	// Data from http handler and data to http handler
	inCh := make(chan httputils.HTTPResponse)
	outCh := make(chan HandlerRequest)

	go func(ctx *zmq.Context, endpoint string) {
		outPort, err = utils.CreateOutputPort(context, endpoint)
		utils.AssertError(err)

		// Map of uuid to requests
		dataMap := make(map[string]chan httputils.HTTPResponse)

		// Start listening in/out channels
		for {
			select {
			case data := <-outCh:
				dataMap[data.Request.Id] = data.ResponseCh
				ip, _ := httputils.Request2IP(data.Request)
				outPort.SendMultipart(ip, 0)
			case resp := <-inCh:
				if respCh, ok := dataMap[resp.Id]; ok {
					log.Println("Resolved channel for response", resp.Id)
					respCh <- resp
					delete(dataMap, resp.Id)
					continue
				}
				log.Println("Didn't find request handler mapping for a given ID", resp.Id)
			}
		}
	}(context, *outputEndpoint)

	// Web server goroutine
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", Handler(outCh))

		s := &http.Server{
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		ln, err := net.Listen("tcp", bindAddr)
		if err != nil {
			log.Println(err.Error())
			exitCh <- syscall.SIGTERM
			return
		}

		log.Printf("Starting listening %v", bindAddr)
		err = s.Serve(ln)
		if err != nil {
			log.Println(err.Error())
			exitCh <- syscall.SIGTERM
			return
		}
	}()

	// Process incoming message forever
	for {
		ip, err := inPort.RecvMultipart(0)
		if err != nil {
			log.Println("Error receiving message:", err.Error())
			continue
		}
		if !runtime.IsValidIP(ip) {
			log.Println("Received invalid IP")
			continue
		}

		resp, err := httputils.IP2Response(ip)
		if err != nil {
			log.Printf("Error converting IP to response: %s", err.Error())
			continue
		}
		inCh <- *resp
	}
}
