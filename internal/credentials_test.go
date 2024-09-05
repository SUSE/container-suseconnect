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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	cr := &Credentials{}
	require.EqualValues(t, '=', cr.separator())

	prepareLogger()
	err := cr.afterParseCheck()
	assert.NotNil(t, err)
	msg := "Can't find username"
	assert.EqualError(t, err, msg)
	shouldHaveLogged(t, msg)

	cr.setValues("username", "suse")
	prepareLogger()
	err = cr.afterParseCheck()
	assert.NotNil(t, err)
	msg = "Can't find password"
	assert.EqualError(t, err, msg)
	shouldHaveLogged(t, msg)

	cr.setValues("password", "1234")
	err = cr.afterParseCheck()
	assert.Nil(t, err)

	locs := cr.locations()
	assert.Contains(t, locs, "/etc/zypp/credentials.d/SCCcredentials")
	assert.Contains(t, locs, "/run/secrets/SCCcredentials")
	assert.Contains(t, locs, "/run/secrets/credentials.d/SCCcredentials")

	// It should log a proper warning.
	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	cr.setValues("unknown", "value")
	assert.Contains(t, buffer.String(), "Warning: Unknown key 'unknown'")
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
	require.Nil(t, err)

	assert.Equal(t, "SCC_a6994b1d3ae14b35agc7cef46b4fff9a", mock.cr.Username)
	assert.Equal(t, "10yb1x6bd159g741ad420fd5aa5083e4", mock.cr.Password)
	assert.Equal(t, "36531d07-a283-441b-a02a-1cd9a88b0d5d", mock.cr.SystemToken)
	assert.False(t, mock.cr.onLocationsNotFound())
}
