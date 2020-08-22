package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/fuzzingbits/hub/pkg/api"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/fuzzingbits/hub/pkg/util/forge/typescript"
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
	a := &api.App{}

	fileName := "./ui/assets/api.ts"
	apiFilePlaceholderFinder := regexp.MustCompile(`\t\/\/ ---- Auto Generated Functions BEGIN ---- \/\/\n([\S\s]*)\t\/\/ ---- Auto Generated Functions END ---- \/\/\n`)

	clientTemplate := template.Must(template.ParseFiles("ops/builder/client.gotemplate"))
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	tsFunctions := bytes.NewBuffer([]byte{})
	if err := clientTemplate.Execute(tsFunctions, ClientFileInput{
		Endpoints: convertTypes(a.GetRoutes()),
	}); err != nil {
		return err
	}

	newFileBytes := apiFilePlaceholderFinder.ReplaceAll(fileBytes, tsFunctions.Bytes())
	if err := ioutil.WriteFile(fileName, newFileBytes, 0644); err != nil {
		return err
	}

	return nil
}

func convertTypes(routes []rooter.Route) []Endpoint {
	endpoints := []Endpoint{}

	for _, route := range routes {
		if route.ExcludeFromTypeScript {
			continue
		}

		method := "get"

		payloadType := convertType(route.Payload)
		if payloadType != "" {
			method = "post"
		}

		endpoints = append(endpoints, Endpoint{
			FunctionName: convertRouteName(route.Path),
			URL:          route.Path,
			Method:       method,
			PayloadType:  payloadType,
			ReturnType:   convertType(route.Response),
		})
	}

	return endpoints
}

func convertType(v interface{}) string {
	if v == nil {
		return ""
	}

	rawTypeString := reflect.ValueOf(v).Type().String()

	return typescript.TranslateReflectTypeString(rawTypeString)
}

func convertRouteName(path string) string {
	path = strings.TrimPrefix(path, "/api")
	path = strings.ReplaceAll(path, "/", " ")
	path = strings.Title(path)
	path = strings.ReplaceAll(path, " ", "")
	path = strings.ToLower(path[0:1]) + path[1:]

	return path
}
