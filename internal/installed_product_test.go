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
)

type NotFoundProvider struct{}

func (m NotFoundProvider) Location() string {
	return "testdata/not-found.xml"
}

func TestFailNonExistantProduct(t *testing.T) {
	var b NotFoundProvider

	_, err := readInstalledProduct(b)
	if err == nil {
		t.Fatal("This file should not exist...")
	}

	if err.Error() != "No base product detected" {
		t.Fatal("Wrong error message")
	}
}

type NotAllowedProvider struct{}

func (m NotAllowedProvider) Location() string {
	return "/etc/shadow"
}

func TestFailNotAllowedProduct(t *testing.T) {
	var b NotAllowedProvider

	_, err := readInstalledProduct(b)
	if err == nil {
		t.Fatal("This file should not be available...")
	}

	if err.Error() != "Can't open base product file: open /etc/shadow: permission denied" {
		t.Fatal("Wrong error message")
	}
}

type BadFormattedProvider struct{}

func (m BadFormattedProvider) Location() string {
	return "testdata/bad.xml"
}

func TestFailBadFormattedProduct(t *testing.T) {
	var b BadFormattedProvider

	_, err := readInstalledProduct(b)
	if err == nil {
		t.Fatal("This file should have a bad format")
	}

	if err.Error() != "Can't parse base product file: EOF" {
		t.Fatal("Wrong error message")
	}
}

type MockProvider struct{}

func (m MockProvider) Location() string {
	return "testdata/installed.xml"
}

func TestMockProvider(t *testing.T) {
	var b MockProvider

	p, err := readInstalledProduct(b)
	if err != nil {
		t.Fatal("It should've read it just fine")
	}

	if p.Identifier != "SLES" {
		t.Fatal("Wrong product name")
	}

	if p.Version != "12" {
		t.Fatal("Wrong product version")
	}

	if p.Arch != "x86_64" {
		t.Fatal("Wrong product arch")
	}

	if p.String() != "SLES-12-x86_64" {
		t.Fatal("Wrong product string")
	}
}

// This test is useless outside SUSE. Added so the go cover tool is happy.
func TestSUSE(t *testing.T) {
	var b SUSEProductProvider

	if _, err := os.Stat(b.Location()); os.IsNotExist(err) {
		_, err = GetInstalledProduct()
		if err == nil {
			t.Fatal("It should fail")
		}

		return
	}

	_, err := GetInstalledProduct()
	if err != nil {
		t.Fatal("We assume that is SUSE, so this should be fine")
	}
}
