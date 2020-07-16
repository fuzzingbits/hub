package config

import (
	"testing"
)

func TestProviderJSON(t *testing.T) {
	startingConfig := TestConfigStruct{
		Name: "Aaron",
		Age:  22,
	}

	targetConfig := TestConfigStruct{
		Name: "George",
		Age:  22,
	}

	config := Config{
		Providers: []Provider{
			ProviderJSON{
				FileLocations: []string{
					"./provider_json_test_files",
				},
				Filename: "config.json",
			},
		},
	}

	configTestHelper(t, config, &startingConfig, &targetConfig, nil, false)
}

func TestProviderJSONNotFound(t *testing.T) {
	config := Config{
		Providers: []Provider{
			ProviderJSON{
				FileLocations: []string{
					"./provider_json_test_files",
				},
				Filename: "file_not_found.json",
			},
		},
	}

	configTestHelper(t, config, &TestConfigStruct{}, &TestConfigStruct{}, nil, false)
}
func TestProviderJSONInvalid(t *testing.T) {
	config := Config{
		Providers: []Provider{
			ProviderJSON{
				FileLocations: []string{
					"./provider_json_test_files",
				},
				Filename: "invalid_json.json",
			},
		},
	}

	configTestHelper(t, config, &TestConfigStruct{}, &TestConfigStruct{}, nil, true)
}
