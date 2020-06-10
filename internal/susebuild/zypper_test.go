// Copyright (c) 2020 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package susebuild

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mssola/capture"
)

func TestParseStdinSuccessful(t *testing.T) {
	content := []byte("RESOLVEURL\nkey:value1\nanother:value2")
	tmp, err := ioutil.TempFile("", "file")
	if err != nil {
		t.Fatalf("Initialization error: %v", err)
	}

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if _, err := tmp.Write(content); err != nil {
		t.Fatalf("Initialization error: %v", err)
	}

	if _, err := tmp.Seek(0, 0); err != nil {
		t.Fatalf("Initialization error: %v", err)
	}

	// stdin switcheroo
	old := os.Stdin
	defer func() {
		os.Stdin = old
	}()
	os.Stdin = tmp

	params, err := ParseStdin()
	if err != nil {
		t.Fatalf("ParseStdin returned an error: %v", err)
	}
	if len(params) != 2 {
		t.Fatalf("There should be two entries, %v instead", len(params))
	}
	if params["key"] != "value1" {
		t.Fatalf("Expected 'key' to contain 'value1', got '%v' instead", params["key"])
	}
	if params["another"] != "value2" {
		t.Fatalf("Expected 'key' to contain 'value2', got '%v' instead", params["another"])
	}
}

func TestParseStdinBadInput(t *testing.T) {
	old := os.Stdin
	defer func() {
		os.Stdin = old
	}()
	os.Stdin = nil

	_, err := ParseStdin()
	if err.Error() != "invalid argument" {
		t.Fatalf("Expected there to be an invalid argument error, got '%v' instead", err)
	}
}

func TestPrintResponseBadResponseFromServer(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{
			bootstrapped: make(chan bool, 1),
			response:     SuseBuildConfig{},
		}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		params := map[string]string{}

		err := PrintResponse(params)
		if err == nil || err.Error() != "Empty response from the server" {
			t.Fatalf("Expecting an error from ReadConfigFromServer, got %v", err)
		}
	})
}

func TestPrintResponseNoCredentials(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{
			bootstrapped: make(chan bool, 1),
			response: SuseBuildConfig{
				InstanceData: "instance data",
			},
		}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		params := map[string]string{}

		err := PrintResponse(params)
		if err == nil || err.Error() != "No credentials given" {
			t.Fatalf("Expecting a 'No credentials given' error, got %v", err)
		}
	})
}

func TestPrintFromConfiguration(t *testing.T) {
	withSuppressedLog(func() {
		res := capture.All(func() {
			printFromConfiguration("/path", &SuseBuildConfig{
				InstanceData: "instance data",
				ServerFqdn:   "test.fqdn.com",
				ServerIp:     "1.1.1.1",
				Username:     "banjo",
				Password:     "kazooie",
				Ca:           "ca",
			})
		})

		lines := strings.Split(string(res.Stdout), "\n")
		expected := []string{
			"RESOLVEDURL",
			"X-Instance-Data:instance data",
			"",
			"https://banjo:kazooie@test.fqdn.com/path\x00",
		}
		if len(lines) != len(expected) {
			t.Fatalf("Expected %v lines, got %v", len(expected), len(lines))
		}
		for k, v := range expected {
			if lines[k] != v {
				t.Fatalf("Expected '%v', got '%v'", v, lines[k])
			}
		}
	})
}
