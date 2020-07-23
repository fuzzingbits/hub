package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockableTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m mockableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

type mockableService struct {
	ErrorCheckError error
	ModRequestError error
}

func (m mockableService) ModRequest(request *http.Request) error {
	return m.ModRequestError
}

func (m mockableService) ErrorCheck(responseBytes []byte) error {
	return m.ErrorCheckError
}

type testJSONResponse struct {
	Ok bool `json:"OK"`
}

type testJSONPayload struct {
	Name string `json:"name"`
}

func TestClient_CurlSimple(t *testing.T) {
	service := &mockableService{}

	transport := &mockableTransport{}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := Client{
		HTTPClient: httpClient,
		Service:    service,
	}

	testTarget := testJSONResponse{}
	testPayload := testJSONPayload{Name: "Aaron"}

	type args struct {
		method   string
		endpoint *url.URL
		payload  interface{}
		target   interface{}
	}

	tests := []struct {
		name         string
		fields       Client
		args         args
		targetTarget interface{}
		wantErr      bool
		setup        func()
	}{
		{
			name:   "Primary test track",
			fields: client,
			args: args{
				method:   http.MethodGet,
				endpoint: &url.URL{},
				target:   &testTarget,
				payload:  testPayload,
			},
			targetTarget: &testJSONResponse{Ok: true},
			wantErr:      false,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					if err := testPayloadHelper(req, testPayload); err != nil {
						return nil, err
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"OK": true}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:   "Bad status code test",
			fields: client,
			args: args{
				endpoint: &url.URL{},
			},
			wantErr: true,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"OK": false}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:   "Malformed JSON Response",
			fields: client,
			args: args{
				endpoint: &url.URL{},
				target:   &testTarget,
			},
			targetTarget: &testJSONResponse{Ok: false},
			wantErr:      true,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"OK" false}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:   "Transport Error",
			fields: client,
			args: args{
				endpoint: &url.URL{},
			},
			wantErr: true,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("random error")
				}
			},
		},
		{
			name:   "Service ErrorCheck Error",
			fields: client,
			args: args{
				endpoint: &url.URL{},
			},
			wantErr: true,
			setup: func() {
				service.ErrorCheckError = errors.New("random error")
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"OK" false}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:   "Service ModRequest Error",
			fields: client,
			args: args{
				endpoint: &url.URL{},
			},
			wantErr: true,
			setup: func() {
				service.ModRequestError = errors.New("random error")
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"OK" false}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:   "Bad Method Error",
			fields: client,
			args: args{
				method:   "!@#$%^&*()_+",
				endpoint: &url.URL{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			// Check for an error
			if err := tt.fields.CurlSimple(tt.args.method, tt.args.endpoint, tt.args.payload, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Client.CurlSimple() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check the state of the target
			if !reflect.DeepEqual(tt.targetTarget, tt.args.target) {
				t.Errorf("Client.CurlSimple() target = %v, want %v", tt.args.target, tt.targetTarget)
			}

			{ // Cleanup/Reset
				transport.RoundTripFunc = nil
				service.ErrorCheckError = nil
				service.ModRequestError = nil
				testTarget.Ok = false
			}
		})
	}
}

func testPayloadHelper(req *http.Request, payload interface{}) error {
	if req.Body == nil {
		if payload != nil {
			return errors.New("no body sent but one was expected")
		}

		return nil // No body and no expected body
	}

	requestBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(payloadBytes, requestBytes) {
		return fmt.Errorf("Request body did not match. Want: %s, Got: %s",
			string(payloadBytes),
			string(requestBytes))
	}

	return nil
}
