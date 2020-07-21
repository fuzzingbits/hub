package main

import (
	"fmt"
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
	entity.UserSession{},
	entity.User{},
	entity.UserSettings{},
}

func main() {
	tsTypesTargetFile := "ui/assets/types.ts"
	newContents := generateTypesContents()

	if err := ioutil.WriteFile(tsTypesTargetFile, newContents, 0644); err != nil {
		log.Print(fmt.Errorf("Could not write file [%s]: %w", tsTypesTargetFile, err))
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
