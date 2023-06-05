package autotakeover

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestTakeover(t *testing.T) {
	const fakeToken string = "FAKETOKEN"

	type testcase struct {
		name      string
		inputPath string
		handlers  map[string]http.HandlerFunc
		err       bool
		expected  string
	}

	var cases []testcase = []testcase{
		{
			name:      "input does not exist",
			inputPath: "invalid_input_file.txt",
			handlers: map[string]http.HandlerFunc{
				"/": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			err:      true,
			expected: "failed to load input file - might not exist:",
		},
		{
			name:      "happy (fake) case",
			inputPath: "test.txt",
			handlers: map[string]http.HandlerFunc{
				"/": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			err:      false,
			expected: "something",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// setup mock server
			mux := http.NewServeMux()
			for path, handler := range tc.handlers {
				mux.HandleFunc(path, handler)
			}
			srv := httptest.NewServer(mux)
			defer srv.Close()

			// override baseUrl package-level variable
			baseUrl = srv.URL

			// capture log package output
			outReader, outWriter, _ := os.Pipe()
			log.SetOutput(outWriter)
			defer func() {
				outWriter.Close()
				log.SetOutput(os.Stderr)
			}()

			Takeover(tc.inputPath, fakeToken)

			bytes, err := io.ReadAll(outReader)
			if err != nil {
				t.Fatalf("failed reading stdout written by log pkg: %e\n", err)
			}

			output := strings.TrimSpace(string(bytes))
			if tc.expected != output {
				t.Fatalf("expected output from log to be %s and got: %s\n", tc.expected, output)
			}
		})
	}
}
