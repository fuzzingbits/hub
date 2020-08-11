package main

import (
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/fuzzingbits/hub/pkg/api"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
)

// ClientFileInput is the input struct to the client template
type ClientFileInput struct {
	Endpoints []Endpoint
}

// Endpoint is a endpoint in the client
type Endpoint struct {
	FunctionName string
	URL          string
	Method       string
	ReturnType   string
	PayloadType  string
}

func buildClientFile() error {
	clientTemplate := template.Must(template.ParseFiles("ops/builder/client.gotemplate"))
	clientFile, err := os.OpenFile("./ui/assets/api.ts", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err

	}
	defer clientFile.Close()
	a := &api.App{}
	if err := clientTemplate.Execute(clientFile, ClientFileInput{
		Endpoints: convertTypes(a.GetRoutes()),
	}); err != nil {
		return err
	}

	return nil
}

func convertTypes(routes []rooter.Route) []Endpoint {
	endpoints := []Endpoint{}

	for _, route := range routes {
		method := "get"
		payloadType := ""
		if route.Payload != nil {
			method = "post"
			payloadType = reflect.TypeOf(route.Payload).Name()
		}
		returnType := ""
		if route.Response != nil {
			returnType = reflect.TypeOf(route.Response).Name()
		}
		endpoints = append(endpoints, Endpoint{
			FunctionName: convertRouteName(route.Path),
			URL:          route.Path,
			Method:       method,
			PayloadType:  payloadType,
			ReturnType:   returnType,
		})
	}

	return endpoints
}

func convertRouteName(path string) string {
	path = strings.TrimPrefix(path, "/api")
	path = strings.ReplaceAll(path, "/", " ")
	path = strings.Title(path)
	path = strings.ReplaceAll(path, " ", "")
	path = strings.ToLower(path[0:1]) + path[1:]

	return path
}
