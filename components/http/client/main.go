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

	httputils "github.com/cascades-fbp/cascades-http/utils"
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	zmq "github.com/pebbe/zmq4"
)

var (
	// Flags
	requestEndpoint  = flag.String("port.req", "", "Component's input port endpoint")
	responseEndpoint = flag.String("port.resp", "", "Component's output port endpoint")
	bodyEndpoint     = flag.String("port.body", "", "Component's output port endpoint")
	errorEndpoint    = flag.String("port.err", "", "Component's error port endpoint")
	jsonFlag         = flag.Bool("json", false, "Print component documentation in JSON")
	debug            = flag.Bool("debug", false, "Enable debug mode")

	// Internal
	// Internal
	reqPort, respPort, bodyPort, errPort *zmq.Socket
	reqCh, respCh, bodyCh, errCh         chan bool
	exitCh                               chan os.Signal
	err                                  error
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
	reqCh = make(chan bool)
	bodyCh = make(chan bool)
	respCh = make(chan bool)
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
	if bodyPort != nil {
		ports++
	}
	if respPort != nil {
		ports++
	}
	if errPort != nil {
		ports++
	}

	waitCh := make(chan bool)
	reqExitCh := make(chan bool, 1)
	go func(num int) {
		total := 0
		for {
			select {
			case v := <-reqCh:
				if v {
					total++
				} else {
					reqExitCh <- true
				}
			case v := <-bodyCh:
				if !v {
					log.Println("BODY port is closed. Interrupting execution")
					exitCh <- syscall.SIGTERM
					break
				} else {
					total++
				}
			case v := <-respCh:
				if !v {
					log.Println("RESP port is closed. Interrupting execution")
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

	// This is obviously dangerous but we need it to deal with our custom CA's
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 30 * time.Second

	// Main loop
	var (
		ip            [][]byte
		clientOptions *httputils.HTTPClientOptions
		request       *http.Request
	)

	log.Println("Started")

	for {
		ip, err = reqPort.RecvMessageBytes(zmq.DONTWAIT)
		if err != nil {
			select {
			case <-reqExitCh:
				log.Println("REQ port is closed. Interrupting execution")
				exitCh <- syscall.SIGTERM
				break
			default:
				// IN port is still open
			}
			time.Sleep(2 * time.Second)
			continue
		}
		if !runtime.IsValidIP(ip) {
			log.Println("Invalid IP:", ip)
			continue
		}

		err = json.Unmarshal(ip[1], &clientOptions)
		if err != nil {
			log.Println("ERROR: failed to unmarshal request options:", err.Error())
			continue
		}
		if clientOptions == nil {
			log.Println("ERROR: received nil request options")
			continue
		}

		if clientOptions.Form != nil {
			request, err = http.NewRequest(clientOptions.Method, clientOptions.URL, strings.NewReader(clientOptions.Form.Encode()))
		} else {
			request, err = http.NewRequest(clientOptions.Method, clientOptions.URL, nil)
		}
		utils.AssertError(err)

		if clientOptions.ContentType != "" {
			request.Header.Add("Content-Type", clientOptions.ContentType)
		}

		for k, v := range clientOptions.Headers {
			request.Header.Add(k, v[0])
		}

		response, err := client.Do(request)
		if err != nil {
			log.Printf("ERROR performing HTTP %s %s: %s", request.Method, request.URL, err.Error())
			if errPort != nil {
				errPort.SendMessageDontwait(runtime.NewPacket([]byte(err.Error())))
			}
			clientOptions = nil
			continue
		}
		resp, err := httputils.Response2Response(response)
		if err != nil {
			log.Printf("ERROR converting response to reply: %s", err.Error())
			if errPort != nil {
				errPort.SendMessageDontwait(runtime.NewPacket([]byte(err.Error())))
			}
			clientOptions = nil
			continue
		}
		ip, err = httputils.Response2IP(resp)
		if err != nil {
			log.Printf("ERROR converting reply to IP: %s", err.Error())
			if errPort != nil {
				errPort.SendMessageDontwait(runtime.NewPacket([]byte(err.Error())))
			}
			clientOptions = nil
			continue
		}

		if respPort != nil {
			respPort.SendMessage(ip)
		}
		if bodyPort != nil {
			bodyPort.SendMessage(runtime.NewPacket(resp.Body))
		}

		select {
		case <-reqCh:
			log.Println("REQ port is closed. Interrupting execution")
			exitCh <- syscall.SIGTERM
			break
		default:
			// file port is still open
		}

		clientOptions = nil
		continue
	}
}

// validateArgs checks all required flags
func validateArgs() {
	if *requestEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *responseEndpoint == "" && *bodyEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

// openPorts create ZMQ sockets and start socket monitoring loops
func openPorts() {
	reqPort, err = utils.CreateInputPort("http/client.req", *requestEndpoint, reqCh)
	utils.AssertError(err)

	if *responseEndpoint != "" {
		respPort, err = utils.CreateOutputPort("http/client.resp", *responseEndpoint, respCh)
		utils.AssertError(err)
	}
	if *bodyEndpoint != "" {
		bodyPort, err = utils.CreateOutputPort("http/client.body", *bodyEndpoint, bodyCh)
		utils.AssertError(err)
	}
	if *errorEndpoint != "" {
		errPort, err = utils.CreateOutputPort("http/client.err", *errorEndpoint, errCh)
		utils.AssertError(err)
	}
}

// closePorts closes all active ports and terminates ZMQ context
func closePorts() {
	log.Println("Closing ports...")
	reqPort.Close()
	if bodyPort != nil {
		bodyPort.Close()
	}
	if respPort != nil {
		respPort.Close()
	}
	if errPort != nil {
		errPort.Close()
	}
	zmq.Term()
}
