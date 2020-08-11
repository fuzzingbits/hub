package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/fuzzingbits/hub/pkg/util/forge/typescript"
)

// TypeScriptInterfaces written to `ui/assets/types.ts`
var TypeScriptInterfaces = []interface{}{
	rooter.Response{},
	entity.CreateUserRequest{},
	entity.ServerStatus{},
	entity.UserContext{},
	entity.User{},
	entity.UserSettings{},
	entity.CreateUserRequest{},
	entity.DeleteUserRequest{},
	entity.UserLoginRequest{},
}

func main() {
	tsTypesTargetFile := "ui/assets/types.ts"
	newContents := generateTypesContents()

	if err := ioutil.WriteFile(tsTypesTargetFile, newContents, 0644); err != nil {
		log.Fatalf("Error building types: %s", err.Error())

	}

	if err := buildClientFile(); err != nil {
		log.Fatalf("Error building API Client: %s", err.Error())
	}
}

func generateTypesContents() []byte {
	// Generate the interfaces
	interfaces := typescript.Generate(TypeScriptInterfaces)

	// Build a list of the interfaces as strings
	interfaceStrings := []string{}
	for _, tsInterface := range interfaces {
		interfaceStrings = append(interfaceStrings, tsInterface.String())
	}

	return []byte(strings.Join(interfaceStrings, "\n\n") + "\n")
}
