package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	httputils "github.com/cascades-fbp/cascades-http/utils"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	zmq "github.com/pebbe/zmq4"
)

var (
	// Flags
	routesEndpoint  = flag.String("port.routes", "", "Component's input port endpoint")
	requestEndpoint = flag.String("port.request", "", "Component's input port endpoint")
	successEndpoint = flag.String("port.success", "", "Component's output port endpoint")
	failEndpoint    = flag.String("port.fail", "", "Component's output port endpoint")
	jsonFlag        = flag.Bool("json", false, "Print component documentation in JSON")
	debug           = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	routesPort, requestPort, failPort *zmq.Socket
	successPorts                      map[string]*zmq.Socket
	err                               error
)

func validateArgs() {
	if *routesEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *requestEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *successEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *failEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	successes := strings.Split(*successEndpoint, ",")
	if len(successes) == 0 {
		fmt.Println("ERROR: invalid length of the SUCCESS array port!")
		flag.Usage()
		os.Exit(1)
	}
}

func openPorts() {
	requestPort, err = utils.CreateInputPort(*requestEndpoint)
	utils.AssertError(err)

	failPort, err = utils.CreateOutputPort(*failEndpoint)
	utils.AssertError(err)

	successes := strings.Split(*successEndpoint, ",")
	successPorts = make(map[string]*zmq.Socket, len(successes))

	var port *zmq.Socket

	for i, endpoint := range patterns {
		port, err = utils.CreateOutputPort(strings.TrimSpace(successes[i]))
		utils.AssertError(err)
		successPorts
	}
}

func closePorts() {
	requestPort.Close()
	failPort.Close()
	for _, p := range patternPorts {
		p.Close()
	}
	for _, p := range successPorts {
		p.Close()
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

	openPorts()
	defer closePorts()

	poller.Add(requestPort, zmq.POLLIN)

	exitCh := utils.HandleInterruption()
	err = runtime.SetupShutdownByDisconnect(requestPort, "http-router.in", exitCh)
	utils.AssertError(err)

	// Main loop
	var (
		index       int = -1
		outputIndex int = -1
		params      url.Values
		ip          [][]byte
		pLength     int     = len(pollItems)
		router      *Router = NewRouter()
	)
	for {
		// Poll sockets
		log.Println("Polling sockets...")

		sockets, err := poller.Poll(-1)
		if err != nil {
			log.Println("Error polling ports:", err.Error())
			os.Exit(1)
		}

		for index, s := range sockets {
			ip, err = s.Socket.RecvMessageBytes(0)
			if err != nil {
				log.Printf("Failed to receive data. Error: %s", err.Error())
				continue
			}
			if !runtime.IsValidIP(ip) {
				log.Println("Received invalid IP")
				continue
			}
		}

		// Pattern arrived

		if index < pLength-1 {
			// Close pattern socket
			port = pollItems[index].Socket
			port.Close()

			// Resolve corresponding output socket index
			outputIndex = -1
			for i, s := range patternPorts {
				if s == port {
					outputIndex = i
				}
			}
			if outputIndex == -1 {
				log.Printf("Failed to resolve output socket index")
				continue
			}

			// Remove closed socket from polling items
			pollItems = append(pollItems[:index], pollItems[index+1:]...)
			pLength -= 1

			// Add pattern to router
			parts := strings.Split(string(ip[1]), " ")
			method := strings.ToUpper(strings.TrimSpace(parts[0]))
			pattern := strings.TrimSpace(parts[1])
			switch method {
			case "GET":
				router.Get(pattern, outputIndex)
			case "POST":
				router.Post(pattern, outputIndex)
			case "PUT":
				router.Put(pattern, outputIndex)
			case "DELETE":
				router.Del(pattern, outputIndex)
			case "HEAD":
				router.Head(pattern, outputIndex)
			case "OPTIONS":
				router.Options(pattern, outputIndex)
			default:
				log.Printf("Unsupported HTTP method %s in pattern %s", method, pattern)
			}
			continue
		}

		// Request arrive
		req, err := httputils.IP2Request(ip)
		if err != nil {
			log.Printf("Failed to convert IP to request. Error: %s", err.Error())
			continue
		}

		outputIndex, params = router.Route(req.Method, req.URI)
		log.Printf("Output index for %s %s: %v (params=%#v)", req.Method, req.URI, outputIndex, params)

		switch outputIndex {
		case NotFound:
			log.Println("Sending Not Found response to FAIL output")
			resp := &httputils.HTTPResponse{
				Id:         req.Id,
				StatusCode: http.StatusNotFound,
			}
			ip, _ = httputils.Response2IP(resp)
			failPort.SendMultipart(ip, 0)
		case MethodNotAllowed:
			log.Println("Sending Method Not Allowed response to FAIL output")
			resp := &httputils.HTTPResponse{
				Id:         req.Id,
				StatusCode: http.StatusMethodNotAllowed,
			}
			ip, _ = httputils.Response2IP(resp)
			failPort.SendMultipart(ip, 0)
		default:
			for k, values := range params {
				req.Form[k] = values
			}
			ip, _ = httputils.Request2IP(req)
			successPorts[outputIndex].SendMultipart(ip, 0)
		}

		index = -1
	}
}
