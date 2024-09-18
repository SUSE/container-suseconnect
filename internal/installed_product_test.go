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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testProductProvider struct {
	location string
}

func (p testProductProvider) Location() string {
	return p.location
}

func TestValidProductFile(t *testing.T) {
	p := testProductProvider{
		location: "testdata/installed.xml",
	}

	ip, err := readInstalledProduct(p)
	require.Nil(t, err)

	assert.Equal(t, ip.Vendor, "SUSE")
	assert.Equal(t, ip.Identifier, "SLES")
	assert.Equal(t, ip.Version, "12")
	assert.Equal(t, ip.Arch, "x86_64")
	assert.Equal(t, ip.String(), "SLES-12-x86_64")
}

func TestInvalidProductFile(t *testing.T) {
	p := testProductProvider{
		location: "testdata/bad.xml",
	}

	_, err := readInstalledProduct(p)

	assert.EqualError(t, err, "Can't parse base product file: EOF")
}

func TestMissingProductFile(t *testing.T) {
	p := testProductProvider{
		location: "/path/that/does/not/exists/file.xml",
	}

	_, err := readInstalledProduct(p)

	assert.EqualError(t, err, "No base product detected")
}

func TestGetInstalledProduct(t *testing.T) {
	var b SUSEProductProvider

	// if the file does not exist is not SUSE
	if _, err := os.Stat(b.Location()); os.IsNotExist(err) {
		t.Skip("Not a SUSE-based system")
	}

	ip, err := GetInstalledProduct()
	assert.Nil(t, err)

	// if the file exist, it should be SUSE/openSUSE
	assert.Contains(t, []string{"SUSE", "openSUSE"}, ip.Vendor)
}
