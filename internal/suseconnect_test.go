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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSUSEConnectData(t *testing.T) {
	data := &SUSEConnectData{}
	require.EqualValues(t, ':', data.separator())

	err := data.afterParseCheck()
	assert.Nil(t, err)
	assert.Equal(t, sccURLStr, data.SccURL)

	locs := data.locations()
	assert.Contains(t, locs, "/etc/SUSEConnect")
	assert.Contains(t, locs, "/run/secrets/SUSEConnect")

	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	data.setValues("unknown", "value")
	assert.Contains(t, buffer.String(), "Warning: Unknown key 'unknown'")
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
	require.Nil(t, err)
	assert.Equal(t, "https://smt.test.lan", mock.data.SccURL)
	assert.True(t, mock.data.Insecure)
}

func TestLocationsNotFound(t *testing.T) {
	var data SUSEConnectData
	locationShouldBeFound = false
	mock := SUSEConnectDataMock{data: &data}

	err := ReadConfiguration(&mock)
	require.Nil(t, err)

	assert.Equal(t, "https://scc.suse.com", mock.data.SccURL)
	assert.False(t, mock.data.Insecure)
}
