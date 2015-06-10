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

package main

import "testing"

func TestCredentials(t *testing.T) {
	cr := &Credentials{}

	if cr.separator() != '=' {
		t.Fatal("Wrong separator")
	}
	err := cr.afterParseCheck()
	if err == nil || err.Error() != "Can't find username" {
		t.Fatal("Wrong error")
	}

	cr.setValues("username", "suse")
	err = cr.afterParseCheck()
	if err == nil || err.Error() != "Can't find password" {
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
	if locs[1] != "/run/secrets/credentials.d/SCCcredentials" {
		t.Fatal("Wrong location")
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
	return []string{"data/credentials.txt"}
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

	err := readConfiguration(&mock)
	if err != nil {
		t.Fatal("This should've been a successful run")
	}
	if mock.cr.Username != "SCC_a6994b1d3ae14b35agc7cef46b4fff9a" {
		t.Fatal("Unexpected name value")
	}
	if mock.cr.Password != "10yb1x6bd159g741ad420fd5aa5083e4" {
		t.Fatal("Unexpected password value")
	}
}
