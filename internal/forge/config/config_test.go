package config

import (
	"os"
	"reflect"
	"testing"
)

type TestConfigStructChild struct {
	Bar       bool    `env:"FORGE_CONFIG_TEST_BAR"`
	PtrString *string `env:"FORGE_CONFIG_TEST_PTRSTRING"`
}

type TestConfigStruct struct {
	Name   string   `env:"FORGE_CONFIG_TEST_NAME" json:"name"`
	Skills []string `env:"FORGE_CONFIG_TEST_SKILLS"`
	Age    int      `env:"FORGE_CONFIG_TEST_AGE" json:"age"`
	PtrInt *int     `env:"FORGE_CONFIG_TEST_PTRINT"`
	Foo    TestConfigStructChild
}

type TestConfigUnexported struct {
	name string `env:"FORGE_CONFIG_TEST_NAME"`
}

func resetTest() {
	os.Unsetenv("FORGE_CONFIG_TEST_BAR")
	os.Unsetenv("FORGE_CONFIG_TEST_PTRSTRING")
	os.Unsetenv("FORGE_CONFIG_TEST_NAME")
	os.Unsetenv("FORGE_CONFIG_TEST_SKILLS")
	os.Unsetenv("FORGE_CONFIG_TEST_AGE")
	os.Unsetenv("FORGE_CONFIG_TEST_PTRINT")
}

func configTestHelper(t *testing.T, config Config, startingConfig interface{}, targetConfig interface{}, targetErr error, justLookForAnyError bool) {
	defer resetTest()

	err := config.Unmarshal(startingConfig)
	if justLookForAnyError {
		if err == nil {
			t.Errorf("No error was found but one was expected")
			return
		}
	} else {
		if err != targetErr {
			t.Errorf("error was not correct, got: \"%v\", want: \"%v\"", err, targetErr)
			return
		}
	}

	if !reflect.DeepEqual(startingConfig, targetConfig) {
		t.Errorf("target config did not match the starting config. got: %+v, want: %+v", startingConfig, targetConfig)
		return
	}
}
