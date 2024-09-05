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

type NotFoundProvider struct{}

func (m NotFoundProvider) Location() string {
	return "testdata/not-found.xml"
}

func TestFailNonExistantProduct(t *testing.T) {
	var b NotFoundProvider

	_, err := readInstalledProduct(b)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "No base product detected")
}

type NotAllowedProvider struct{}

func (m NotAllowedProvider) Location() string {
	return "/etc/shadow"
}

func TestFailNotAllowedProduct(t *testing.T) {
	var b NotAllowedProvider

	_, err := readInstalledProduct(b)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can't open base product file: open /etc/shadow: permission denied")
}

type BadFormattedProvider struct{}

func (m BadFormattedProvider) Location() string {
	return "testdata/bad.xml"
}

func TestFailBadFormattedProduct(t *testing.T) {
	var b BadFormattedProvider

	_, err := readInstalledProduct(b)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can't parse base product file: EOF")
}

type MockProvider struct{}

func (m MockProvider) Location() string {
	return "testdata/installed.xml"
}

func TestMockProvider(t *testing.T) {
	var b MockProvider

	p, err := readInstalledProduct(b)
	require.Nil(t, err)

	assert.Equal(t, "SLES", p.Identifier)
	assert.Equal(t, "12", p.Version)
	assert.Equal(t, "x86_64", p.Arch)
	assert.Equal(t, "SLES-12-x86_64", p.String())
}

// This test is useless outside SUSE. Added so the go cover tool is happy.
func TestSUSE(t *testing.T) {
	var b SUSEProductProvider

	if _, err := os.Stat(b.Location()); os.IsNotExist(err) {
		_, err = GetInstalledProduct()
		assert.NotNil(t, err)

		return
	}

	_, err := GetInstalledProduct()
	assert.Nil(t, err)
}
