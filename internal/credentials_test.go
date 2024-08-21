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

func TestCredentials(t *testing.T) {
	cr := &Credentials{}

	if cr.separator() != '=' {
		t.Fatal("Wrong separator")
	}
	prepareLogger()
	err := cr.afterParseCheck()
	msg := "Can't find username"
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}

	cr.setValues("username", "suse")
	prepareLogger()
	msg = "Can't find password"
	err = cr.afterParseCheck()
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}

	cr.setValues("password", "1234")
	err = cr.afterParseCheck()
	if err != nil {
		t.Fatal("There should not be an error")
	}

	locs := cr.locations()
	if locs[0] != "/etc/zypp/credentials.d/SCCcredentials" {
		t.Fatal("Wrong location")
	}
	if locs[1] != "/run/secrets/SCCcredentials" {
		t.Fatal("Wrong location")
	}
	if locs[2] != "/run/secrets/credentials.d/SCCcredentials" {
		t.Fatal("Wrong location")
	}

	// It should log a proper warning.
	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	cr.setValues("unknown", "value")
	if !strings.Contains(buffer.String(), "Warning: Unknown key 'unknown'") {
		t.Fatal("Wrong warning!")
	}
}

// In the following test we will create a mock that just wraps up the
// `Credentials` struct and replaces its `location` function for something that
// can be tested. We test for a successful run, since all the possible errors
// have already been tested in the `configuration_test.go` file.

type CredentialsMock struct {
	cr *Credentials
}

func (mock *CredentialsMock) locations() []string {
	return []string{"testdata/credentials.txt"}
}

func (mock *CredentialsMock) onLocationsNotFound() bool {
	return mock.cr.onLocationsNotFound()
}

func (mock *CredentialsMock) separator() byte {
	return mock.cr.separator()
}

func (mock *CredentialsMock) setValues(key, value string) {
	mock.cr.setValues(key, value)
}

func (mock *CredentialsMock) afterParseCheck() error {
	return mock.cr.afterParseCheck()
}

func TestIntegrationCredentials(t *testing.T) {
	var credentials Credentials
	mock := CredentialsMock{cr: &credentials}

	err := ReadConfiguration(&mock)
	if err != nil {
		t.Fatal("This should've been a successful run")
	}
	if mock.cr.Username != "SCC_a6994b1d3ae14b35agc7cef46b4fff9a" {
		t.Fatal("Unexpected name value")
	}
	if mock.cr.Password != "10yb1x6bd159g741ad420fd5aa5083e4" {
		t.Fatal("Unexpected password value")
	}
	if mock.cr.SystemToken != "36531d07-a283-441b-a02a-1cd9a88b0d5d" {
		t.Fatal("Unexpected system_token value")
	}
	if mock.cr.onLocationsNotFound() {
		t.Fatalf("It should've been false")
	}
}
