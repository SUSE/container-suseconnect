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
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// copyHostFileToTemp copies the hosts file from the fixtures directory into a
// file that is meant to be temporary and that has a randomized name. Use this
// temporary file instead of the original hosts file inside of the fixtures path
// on these tests.
func copyHostFileToTemp(mode os.FileMode) string {
	data, err := os.ReadFile(fixturesPath("hosts"))
	if err != nil {
		fmt.Printf("Read file error: %v\n", err)
		return ""
	}

	path := fixturesPath(fmt.Sprintf("testfile%v", rand.Int()))
	err = os.WriteFile(path, data, mode)
	if err != nil {
		fmt.Printf("Write file error: %v\n", err)
		return ""
	}

	return path
}

func TestUpdateHostsFileCouldNotRead(t *testing.T) {
	hostsFile = "/bubblegloop-swamp"
	err := UpdateHostsFile("hostname", "1.1.1.1")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "can't read /bubblegloop-swamp file")
}

func TestUpdateHostsFileCouldNotWrite(t *testing.T) {
	hostsFile = copyHostFileToTemp(0o400)
	require.NotEmpty(t, hostsFile)
	defer os.Remove(hostsFile)

	err := UpdateHostsFile("hostname", "1.1.1.1")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "can't write")
}

func TestUpdateHostsFileSuccessful(t *testing.T) {
	hostsFile = copyHostFileToTemp(0o644)
	require.NotEmpty(t, hostsFile)
	defer os.Remove(hostsFile)

	before, err := os.ReadFile(hostsFile)
	require.Nil(t, err)

	err = UpdateHostsFile("test-hostname", "1.1.1.1")
	require.Nil(t, err)

	after, err := os.ReadFile(hostsFile)
	require.Nil(t, err)

	assert.Contains(t, string(after), string(before))
	assert.Contains(t, string(after), "1.1.1.1 test-hostname test-hostname")
}

func TestUpdateHostsFileUpdateExistingEntry(t *testing.T) {
	hostsFile = copyHostFileToTemp(0o644)
	require.NotEmpty(t, hostsFile)
	defer os.Remove(hostsFile)

	err := UpdateHostsFile("ip6-localnet", "1.1.1.1")
	require.Nil(t, err)

	after, err := os.ReadFile(hostsFile)
	require.Nil(t, err)

	assert.Contains(t, string(after), "1.1.1.1 ip6-localnet ip6-localnet")
}
