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
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

// testCommand implements the commander interface with an attribute for mocking
// purposes.
type testCommand struct {
	shouldFail bool
}

// Returns an error if `shouldFail` is set to true and nil otherwise.
func (t testCommand) Run() error {
	if t.shouldFail {
		return errors.New("I AM ERROR")
	}
	return nil
}

// Run this before each test to get the fixtures path right.
func beforeTest() {
	hashFilePath = fixturesPath("valid.sha256")
	caFilePath = fixturesPath("valid.pem")
}

// Returns the full path to the given fixture file.
func fixturesPath(file string) string {
	path, _ := os.Getwd()
	return filepath.Join(path, "fixtures", file)
}

// Tests start here

func TestNoUpdateIsNeeded(t *testing.T) {
	beforeTest()

	if updateNeeded("valid") {
		t.Fatal("Should not be needed")
	}

	if !updateNeeded("nope") {
		t.Fatal("Should be needed")
	}
}

func TestUpdateIsNeededNoFile(t *testing.T) {
	beforeTest()

	hashFilePath = "/tmp/wubalubadubdub"
	b := updateNeeded("thing")
	if !b {
		t.Fatal("Expected update to be needed")
	}
}

func TestUpdateIsNeededCouldNotReadFile(t *testing.T) {
	beforeTest()

	hashFilePath = "/proc/1/mem"
	b := updateNeeded("thing")
	if !b {
		t.Fatal("Expected update to be needed")
	}
}

func TestSaveCAFileBadWrite(t *testing.T) {
	beforeTest()

	hashFilePath = fixturesPath(fmt.Sprintf("file%v.sha256", rand.Int()))
	caFilePath = "/wubalubadubdub"
	cmd := testCommand{shouldFail: false}

	err := saveCAFile(cmd, "valid")
	os.Remove(hashFilePath)
	os.Remove(caFilePath)

	if err == nil {
		t.Fatal("Should've failed")
	}
}

func TestSaveCAFileBadCommand(t *testing.T) {
	beforeTest()

	hashFilePath = fixturesPath(fmt.Sprintf("file%v.sha256", rand.Int()))
	caFilePath = fixturesPath(fmt.Sprintf("file%v.pem", rand.Int()))
	cmd := testCommand{shouldFail: true}

	err := saveCAFile(cmd, "valid")
	os.Remove(hashFilePath)
	os.Remove(caFilePath)

	if err == nil {
		t.Fatal("Expected error to be non-nil\n")
	}

	if err.Error() != "I AM ERROR" {
		t.Fatalf("Expected another error, got: %v\n", err)
	}
}

func TestSaveCAFileSuccess(t *testing.T) {
	beforeTest()

	hashFilePath = fixturesPath("tmp.sha256")
	cmd := testCommand{shouldFail: false}

	err := saveCAFile(cmd, "valid")
	if err != nil {
		os.Remove(hashFilePath)
		t.Fatalf("Expected error to be nil: %v\n", err)
	}

	b, _ := os.ReadFile(hashFilePath)
	os.Remove(hashFilePath)

	hash := sha256.New()
	io.WriteString(hash, "valid")
	if string(b) != string(hash.Sum(nil)) {
		t.Fatal("Bad checksum")
	}

	b, _ = os.ReadFile(caFilePath)
	if string(b) != "valid" {
		t.Fatalf("Wrong certificate. Expected 'valid', got '%v'\n", string(b))
	}
}
