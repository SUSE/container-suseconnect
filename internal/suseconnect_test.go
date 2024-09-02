//
// Copyright (c) 2015 SUSE LLC. All rights reserved.
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

package containersuseconnect

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestSUSEConnectData(t *testing.T) {
	data := &SUSEConnectData{}

	if data.separator() != ':' {
		t.Fatal("Wrong separator")
	}

	err := data.afterParseCheck()
	if err != nil {
		t.Fatal("There should not be an error")
	}

	if data.SccURL != sccURLStr {
		t.Fatal("The URL should be the one from sccURLstr")
	}

	locs := data.locations()

	if locs[0] != "/etc/SUSEConnect" {
		t.Fatal("Wrong location")
	}

	if locs[1] != "/run/secrets/SUSEConnect" {
		t.Fatal("Wrong location")
	}

	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	data.setValues("unknown", "value")

	// It should log a proper warning.
	if !strings.Contains(buffer.String(), "Warning: Unknown key 'unknown'") {
		t.Fatal("Wrong warning!")
	}
}

// In the following test we will create a mock that just wraps up the
// `SUSEConnectData` struct and replaces its `location` function for something
// that can be tested. We test for a successful run, since all the possible
// errors have already been tested in the `configuration_test.go` file.

type SUSEConnectDataMock struct {
	data *SUSEConnectData
}

var locationShouldBeFound = true

func (mock *SUSEConnectDataMock) locations() []string {
	if locationShouldBeFound {
		return []string{"testdata/suseconnect.txt"}
	}

	return []string{"testdata/notfound.txt"}
}

func (mock *SUSEConnectDataMock) onLocationsNotFound() bool {
	return mock.data.onLocationsNotFound()
}

func (mock *SUSEConnectDataMock) separator() byte {
	return mock.data.separator()
}

func (mock *SUSEConnectDataMock) setValues(key, value string) {
	mock.data.setValues(key, value)
}

func (mock *SUSEConnectDataMock) afterParseCheck() error {
	return mock.data.afterParseCheck()
}

func TestIntegrationSUSEConnectData(t *testing.T) {
	var data SUSEConnectData
	locationShouldBeFound = true
	mock := SUSEConnectDataMock{data: &data}

	err := ReadConfiguration(&mock)
	if err != nil {
		t.Fatal("This should've been a successful run")
	}

	if mock.data.SccURL != "https://smt.test.lan" {
		t.Fatal("Unexpected URL value")
	}

	if !mock.data.Insecure {
		t.Fatal("Unexpected Insecure value")
	}
}

func TestLocationsNotFound(t *testing.T) {
	var data SUSEConnectData
	locationShouldBeFound = false
	mock := SUSEConnectDataMock{data: &data}

	err := ReadConfiguration(&mock)
	if err != nil {
		t.Fatal("This should've been a successful run")
	}

	if mock.data.SccURL != "https://scc.suse.com" {
		t.Fatal("It should've been scc.suse.com")
	}
}
