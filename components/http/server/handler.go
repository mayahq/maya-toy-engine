package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	httputils "github.com/cascades-fbp/cascades-http/utils"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	timeout = time.Duration(15) * time.Second
)

type HandlerRequest struct {
	ResponseCh chan httputils.HTTPResponse
	Request    *httputils.HTTPRequest
}

func respondWithTimeout(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(rw, "Couldn't process request in a given time")
}

func Handler(out chan HandlerRequest) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		log.Println("Handler:", req.Method, req.RequestURI)

		id, _ := uuid.NewV4()
		r := httputils.Request2Request(req)
		r.Id = id.String()

		hr := &HandlerRequest{
			ResponseCh: make(chan httputils.HTTPResponse),
			Request:    r,
		}

		// Send request to OUT port
		log.Println("Sending request to out channel (for OUTPUT port)")
		select {
		case out <- *hr:
		case <-time.Tick(timeout):
			respondWithTimeout(rw)
			return
		}

		// Wait for response from IN port
		log.Println("Waiting for response from a channel port (from INPUT port)")

		var resp httputils.HTTPResponse
		select {
		case resp = <-hr.ResponseCh:
		case <-time.Tick(timeout):
			respondWithTimeout(rw)
			return
		}

		log.Println("Data arrived. Responding to HTTP response...")
		for name, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(name, value)
			}
		}
		rw.WriteHeader(resp.StatusCode)
		fmt.Fprint(rw, string(resp.Body))
	}
}
