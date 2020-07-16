package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is foobar
type Client struct {
	HTTPClient *http.Client
	Service    Service
}

// Service is foobar
type Service interface {
	ModRequest(request *http.Request) error
	ErrorCheck(responseBytes []byte) error
}

// CurlSimple is foobar
func (client Client) CurlSimple(method string, endpoint *url.URL, payload interface{}, target interface{}) error {
	// Write the payload as json to a buffer
	var payloadBuffer io.Reader
	if payload != nil {
		payloadBytes, _ := json.Marshal(payload)
		payloadBuffer = bytes.NewBuffer(payloadBytes)
	}

	// Build a new request
	request, err := http.NewRequest(method, endpoint.String(), payloadBuffer)
	if err != nil {
		return err
	}

	return client.Curl(request, target)
}

// Curl makes the request using the service interface functions and unmarshals the response into the target
func (client Client) Curl(request *http.Request, target interface{}) error {
	// Apply any request modifications via the service
	if client.Service != nil {
		if err := client.Service.ModRequest(request); err != nil {
			return err
		}
	}

	// Make the call to get the response bytes
	responseBytes, err := client.Call(request)
	if err != nil {
		return err
	}

	// If there is a target try to unmarshal the response bytes to the target
	if target != nil {
		if err := json.Unmarshal(responseBytes, target); err != nil {
			return err
		}
	}

	// Check for additional errors via the service
	if client.Service != nil {
		if err := client.Service.ErrorCheck(responseBytes); err != nil {
			return err
		}
	}

	return nil
}

// Call makes the request using the client and returns the response bytes
func (client Client) Call(request *http.Request) ([]byte, error) {
	// Make the network request
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// Read the body of the response
	body, _ := ioutil.ReadAll(response.Body)

	// Check for bad status codes
	if response.StatusCode >= 400 {
		return body, fmt.Errorf("Curl Error Code: [%d]", response.StatusCode)
	}

	return body, nil
}
