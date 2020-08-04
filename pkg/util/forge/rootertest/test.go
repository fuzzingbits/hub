package rootertest

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// TestCase to be tested
type TestCase struct {
	// Test Name
	Name string
	// Request
	Request    *http.Request
	Method     string
	URL        string
	Body       io.Reader
	RequestMod func(req *http.Request)
	// Response Checks
	TargetStatusCode       int
	TargetResponseBytes    []byte
	SkipResponseBytesCheck bool
	CustomResponseChecker  func(t *testing.T, response *http.Response)
}

// Test all the provided test cases
func Test(t *testing.T, handler http.Handler, testCases []TestCase) {
	var testCaseNameMap = map[string]bool{}

	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, testCase := range testCases {
		if strings.Contains(testCase.Name, " ") {
			t.Fatalf("Test Case Name Can Not Contain Spaces: '%s'", testCase.Name)
		}

		_, alreadyExists := testCaseNameMap[testCase.Name]
		if alreadyExists {
			t.Fatalf("Test Case Names Must Be Unique: '%s'", testCase.Name)
		}
		testCaseNameMap[testCase.Name] = true

		t.Run(testCase.Name, func(t *testing.T) {
			// Build the request if one is not set
			if testCase.Request == nil {
				var err error
				testCase.Request, err = http.NewRequest(
					testCase.Method,
					testCase.URL,
					testCase.Body,
				)
				if err != nil {
					log.Fatal(err)
				}
			}

			tsURL, _ := url.Parse(ts.URL)
			testCase.Request.URL.Host = tsURL.Host
			testCase.Request.URL.Scheme = tsURL.Scheme

			// Modify the request if RequestMod func is set
			if testCase.RequestMod != nil {
				testCase.RequestMod(testCase.Request)
			}

			// Make the request
			response, err := http.DefaultClient.Do(testCase.Request)
			if err != nil {
				log.Fatal(err)
			}

			// Always confirm the status code
			if response.StatusCode != testCase.TargetStatusCode {
				t.Errorf(
					"%s return %d instead of %d",
					testCase.URL,
					response.StatusCode,
					testCase.TargetStatusCode,
				)
			}

			// Compare the response bytes
			if !testCase.SkipResponseBytesCheck {
				// Read out all the bytes
				responseBytes, err := ioutil.ReadAll(response.Body)
				response.Body.Close()
				if err != nil {
					log.Fatal(err)
				}

				if !reflect.DeepEqual(responseBytes, testCase.TargetResponseBytes) {
					t.Fatalf(
						"%s returned: %s expected: %s",
						testCase.URL,
						string(responseBytes),
						string(testCase.TargetResponseBytes),
					)
				}
			}

			// Use the custom response checker
			if testCase.CustomResponseChecker != nil {
				testCase.CustomResponseChecker(t, response)
			}
		})
	}
}
