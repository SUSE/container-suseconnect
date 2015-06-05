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

import (
	"strings"
	"testing"
)

var credentials = `
username=SCC_a6994b1d3ae14b35agc7cef46b4fff9a
password=10yb1x6bd159g741ad420fd5aa5083e4
`
var credentialsWithComments = `
username=SCC_a6994b1d3ae14b35agc7cef46b4fff9a
# Some comment
password=10yb1x6bd159g741ad420fd5aa5083e4

`
var credentialsWithoutUsername = `
password=10yb1x6bd159g741ad420fd5aa5083e4
`
var credentialsWithoutPassword = `
username=SCC_a6994b1d3ae14b35agc7cef46b4fff9a
`

func TestParseCredentials(t *testing.T) {
	reader := strings.NewReader(credentials)

	credentials, err := ParseCredentials(reader)
	if err != nil {
		t.Errorf(err.Error())
	}

	if credentials.Username != "SCC_a6994b1d3ae14b35agc7cef46b4fff9a" {
		t.Fail()
	}

	if credentials.Password != "10yb1x6bd159g741ad420fd5aa5083e4" {
		t.Fail()
	}
}

func TestParseCredentialsWithComments(t *testing.T) {
	reader := strings.NewReader(credentialsWithComments)

	credentials, err := ParseCredentials(reader)
	if err != nil {
		t.Errorf(err.Error())
	}

	if credentials.Username != "SCC_a6994b1d3ae14b35agc7cef46b4fff9a" {
		t.Errorf(credentials.Username)
	}

	if credentials.Password != "10yb1x6bd159g741ad420fd5aa5083e4" {
		t.Errorf(credentials.Password)
	}
}

func TestParseCredentialsWithoutUsername(t *testing.T) {
	reader := strings.NewReader(credentialsWithoutUsername)

	_, err := ParseCredentials(reader)
	if err == nil {
		t.Fail()
	}
}

func TestParseCredentialsWithoutPassword(t *testing.T) {
	reader := strings.NewReader(credentialsWithoutPassword)

	_, err := ParseCredentials(reader)
	if err == nil {
		t.Fail()
	}
}
