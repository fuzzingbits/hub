package rootertest

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestCase to be tested
type TestCase struct {
	Name          string
	URL           string
	StatusCode    int
	ResponseBytes []byte
	Checker       func(response *http.Response)
}

// Test all the provided test cases
func Test(t *testing.T, handler http.Handler, testCases []TestCase) {
	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			response, err := http.Get(ts.URL + testCase.URL)
			if err != nil {
				log.Fatal(err)
			}

			responseBytes, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if response.StatusCode != testCase.StatusCode {
				t.Errorf(
					"%s return %d instead of %d",
					testCase.URL,
					response.StatusCode,
					testCase.StatusCode,
				)
			}

			if !reflect.DeepEqual(responseBytes, testCase.ResponseBytes) {
				t.Errorf(
					"%s returned: %s expected: %s",
					testCase.URL,
					string(responseBytes),
					string(testCase.ResponseBytes),
				)
			}

			if testCase.Checker != nil {
				testCase.Checker(response)
			}
		})
	}
}
