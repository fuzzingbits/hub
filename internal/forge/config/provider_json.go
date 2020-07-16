package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ProviderJSON is foobar
type ProviderJSON struct {
	FileLocations []string
	Filename      string
}

// Unmarshal is foobar
func (p ProviderJSON) Unmarshal(target interface{}) error {
	for _, fileLocation := range p.FileLocations {
		filename := p.getFilename(fileLocation)
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(fileBytes, target); err != nil {
			return err
		}

		return nil
	}

	// TODO: return no file found error?
	return nil
}

func (p ProviderJSON) getFilename(fileLocation string) string {
	return fmt.Sprintf("%s/%s", fileLocation, p.Filename)
}
