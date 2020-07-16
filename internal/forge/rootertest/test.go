package rootertest

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestCase to be tested
type TestCase struct {
	// Test Name
	Name string
	// Request
	Method     string
	URL        string
	Body       io.Reader
	RequestMod func(req *http.Request)
	// Response Checks
	TargetStatusCode       int
	SkipResponseBytesCheck bool
	TargetResponseBytes    []byte
	ResponseChecker        func(response *http.Response) error
}

// Test all the provided test cases
func Test(t *testing.T, handler http.Handler, testCases []TestCase) {
	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req, err := http.NewRequest(testCase.Method, ts.URL+testCase.URL, testCase.Body)
			if err != nil {
				log.Fatal(err)
			}

			if testCase.RequestMod != nil {
				testCase.RequestMod(req)
			}

			response, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			responseBytes, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if response.StatusCode != testCase.TargetStatusCode {
				t.Errorf(
					"%s return %d instead of %d",
					testCase.URL,
					response.StatusCode,
					testCase.TargetStatusCode,
				)
			}

			if !testCase.SkipResponseBytesCheck {
				if !reflect.DeepEqual(responseBytes, testCase.TargetResponseBytes) {
					t.Errorf(
						"%s returned: %s expected: %s",
						testCase.URL,
						string(responseBytes),
						string(testCase.TargetResponseBytes),
					)
				}
			}

			if testCase.ResponseChecker != nil {
				if err := testCase.ResponseChecker(response); err != nil {
					log.Fatal(err)
				}
			}
		})
	}
}
