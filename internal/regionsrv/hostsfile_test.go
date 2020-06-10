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

package regionsrv

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
)

// copyHostFileToTemp copies the hosts file from the fixtures directory into a
// file that is meant to be temporary and that has a randomized name. Use this
// temporary file instead of the original hosts file inside of the fixtures path
// on these tests.
func copyHostFileToTemp(mode os.FileMode) string {
	data, err := ioutil.ReadFile(fixturesPath("hosts"))
	if err != nil {
		fmt.Printf("Read file error: %v\n", err)
		return ""
	}

	path := fixturesPath(fmt.Sprintf("testfile%v", rand.Int()))
	err = ioutil.WriteFile(path, data, mode)
	if err != nil {
		fmt.Printf("Write file error: %v\n", err)
		return ""
	}

	return path
}

// Test suite below

func TestUpdateHostsFileCouldNotRead(t *testing.T) {
	hostsFile = "/bubblegloop-swamp"
	err := UpdateHostsFile("hostname", "1.1.1.1")

	if err == nil || !strings.Contains(err.Error(), "Can't read /bubblegloop-swamp file") {
		t.Fatalf("Expected 'Can't read /bubblegloop-swamp file', got '%v'\n", err)
	}
}

func TestUpdateHostsFileCouldNotWrite(t *testing.T) {
	hostsFile = copyHostFileToTemp(0400)
	if hostsFile == "" {
		t.Fatalf("Failed to initialize hosts file")
	}

	defer os.Remove(hostsFile)

	err := UpdateHostsFile("hostname", "1.1.1.1")
	if err == nil || !strings.Contains(err.Error(), "Can't write") {
		t.Fatalf("Expected a write error, got '%v'\n", err)
	}
}

func TestUpdateHostsFileSuccessful(t *testing.T) {
	hostsFile = copyHostFileToTemp(0644)
	if hostsFile == "" {
		t.Fatalf("Failed to initialize hosts file")
	}

	defer os.Remove(hostsFile)

	before, err := ioutil.ReadFile(hostsFile)

	err = UpdateHostsFile("test-hostname", "1.1.1.1")
	if err != nil {
		t.Fatalf("Expected a nil error, got: %v", err)
	}

	after, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		t.Fatalf("Expected a nil error, got: %v", err)
	}
	if !strings.Contains(string(after), string(before)) {
		t.Fatalf("%v\nshould contain\n%v", string(after), string(before))
	}

	expected := "1.1.1.1 test-hostname test-hostname"
	if !strings.Contains(string(after), expected) {
		t.Fatalf("%v\nshould contain\n%v", string(after), expected)
	}
}
