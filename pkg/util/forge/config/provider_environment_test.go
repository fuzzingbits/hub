package config

import (
	"os"
	"testing"
)

func TestProviderEnvironment(t *testing.T) {
	var stringPointerExample = new(string)
	*stringPointerExample = "foobar2"

	startingConfig := TestConfigStruct{
		Name: "Aaron",
		Age:  22,
	}

	targetConfig := TestConfigStruct{
		Name: "George",
		Age:  42,
		Foo: TestConfigStructChild{
			Bar:       true,
			PtrString: stringPointerExample,
		},
	}

	os.Setenv("FORGE_CONFIG_TEST_NAME", "George")
	os.Setenv("FORGE_CONFIG_TEST_AGE", "42")
	os.Setenv("FORGE_CONFIG_TEST_BAR", "true")
	os.Setenv("FORGE_CONFIG_TEST_PTRSTRING", *stringPointerExample)

	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, &startingConfig, &targetConfig, nil, false)
}

func TestProviderEnvironmentInvalidInt(t *testing.T) {
	startingConfig := TestConfigStruct{
		Name: "Aaron",
		Age:  22,
	}

	os.Setenv("FORGE_CONFIG_TEST_BAR", "not a valid bool")

	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, &startingConfig, &startingConfig, nil, true)
}

func TestProviderEnvironmentInvalidBool(t *testing.T) {
	startingConfig := TestConfigStruct{
		Name: "Aaron",
		Age:  22,
	}

	os.Setenv("FORGE_CONFIG_TEST_AGE", "not a valid int")

	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, &startingConfig, &startingConfig, nil, true)
}

func TestProviderEnvironmentErrUnexportedField(t *testing.T) {
	startingConfig := TestConfigUnexported{
		name: "Aaron",
	}

	os.Setenv("FORGE_CONFIG_TEST_NAME", "George")

	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, &startingConfig, &startingConfig, ErrUnexportedField, false)
}

func TestProviderEnvironmentErrUnsupportedType(t *testing.T) {
	startingConfig := TestConfigStruct{
		Skills: []string{"go"},
	}

	os.Setenv("FORGE_CONFIG_TEST_SKILLS", "go")

	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, &startingConfig, &startingConfig, ErrUnsupportedType, false)
}

func TestProviderEnvironmentPointerSetError(t *testing.T) {
	var intPointerExample = new(int)
	*intPointerExample = 22

	startingConfig := TestConfigStruct{
		Name: "Aaron",
		Age:  22,
	}

	targetConfig := TestConfigStruct{
		Name: "Aaron",
		Age:  22,
	}

	os.Setenv("FORGE_CONFIG_TEST_PTRINT", "not an int")

	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, &startingConfig, &targetConfig, nil, true)
}

func TestProviderEnvironmentNotPointer(t *testing.T) {
	config := Config{
		Providers: []Provider{
			ProviderEnvironment{},
		},
	}

	configTestHelper(t, config, TestConfigStruct{}, TestConfigStruct{}, nil, true)
}
