package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cascades-fbp/cascades/runtime"
)

// HTTPClientOptions describe options IP for client configuration
type HTTPClientOptions struct {
	URL         string              `json:"url"`
	Method      string              `json:"method"`
	ContentType string              `json:"content-type"`
	Headers     map[string][]string `json:"headers"`
	Form        url.Values          `json:"form"`
}

//
// HTTPRequest data structure for IP
//
type HTTPRequest struct {
	ID     string              `json:"id"`      // Assigned by server component
	Method string              `json:"method"`  // GET/POST/PUT/etc
	URI    string              `json:"uri"`     // Full URL that hit the server
	Header map[string][]string `json:"headers"` // Map of headers
	Form   map[string][]string `json:"form"`    // Map of GET/POST/PUT values
}

//
// HTTPResponse data structure for IP
//
type HTTPResponse struct {
	ID         string              `json:"id"`      // Retrieved from request structure
	StatusCode int                 `json:"status"`  // Response HTTP status code
	Header     map[string][]string `json:"headers"` // Map of headers
	Body       []byte              `json:"body"`    // Body of the response
}

// Request2Request create our internal request structure based on the standard one
func Request2Request(request *http.Request) *HTTPRequest {
	// Parse GET/POST/PUT params into request.Form
	request.ParseForm()
	// Create data structure
	res := &HTTPRequest{
		Method: request.Method,
		URI:    request.RequestURI,
		Header: request.Header,
		Form:   request.Form,
	}
	return res
}

// Response2Response create our internal response structure based on the standard one
func Response2Response(response *http.Response) (*HTTPResponse, error) {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	rep := &HTTPResponse{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Body:       body,
	}
	return rep, nil
}

// Request2IP converts a given request to IP
func Request2IP(request *HTTPRequest) ([][]byte, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return runtime.NewPacket(payload), nil
}

// Response2IP сonverts a given response to IP
func Response2IP(response *HTTPResponse) ([][]byte, error) {
	payload, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return runtime.NewPacket(payload), nil
}

// IP2Request сonverts a given IP to request structure
func IP2Request(ip [][]byte) (*HTTPRequest, error) {
	var req *HTTPRequest
	err := json.Unmarshal(ip[1], &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// IP2Response сonverts a given IP to response structure
func IP2Response(ip [][]byte) (*HTTPResponse, error) {
	var res *HTTPResponse
	err := json.Unmarshal(ip[1], &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
