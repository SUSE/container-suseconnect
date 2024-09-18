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
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCommand implements the commander interface with an attribute for mocking
// purposes.
type testCommand struct {
	fail bool
}

func (t testCommand) Run() error {
	if t.fail {
		return errors.New("I am an error")
	}

	return nil
}

func NewSuccessCmd() commander {
	return testCommand{fail: false}
}

func NewErrorCmd() commander {
	return testCommand{fail: true}
}

// Run this before each test to get the fixtures path right.
func beforeTest() {
	hashFilePath = fixturesPath("valid.md5")
	caFilePath = fixturesPath("valid.pem")
}

// Returns the full path to the given fixture file.
func fixturesPath(file string) string {
	path, _ := os.Getwd()
	return filepath.Join(path, "fixtures", file)
}

func TestUpdateIsNeededHashMismatch(t *testing.T) {
	beforeTest()

	assert.True(t, updateNeeded("nope"))
}

func TestNoUpdateIsNeededHashMatch(t *testing.T) {
	beforeTest()

	assert.False(t, updateNeeded("valid"))
}

func TestUpdateIsNeededNoFile(t *testing.T) {
	beforeTest()

	hashFilePath = "/tmp/wubalubadubdub"
	assert.True(t, updateNeeded("thing"))
}

func TestUpdateIsNeededCouldNotReadFile(t *testing.T) {
	beforeTest()

	hashFilePath = "/proc/1/mem"
	assert.True(t, updateNeeded("thing"))
}

func TestSafeCAFileBadWrite(t *testing.T) {
	beforeTest()

	hashFilePath = fixturesPath(fmt.Sprintf("file%v.md5", rand.Int()))
	defer os.Remove(hashFilePath)

	caFilePath = "/path/that/does/not/exist/file"
	defer os.Remove(caFilePath)

	err := safeCAFile(NewSuccessCmd(), "valid")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "open /path/that/does/not/exist/file: no such file or directory")
}

func TestSafeCAFileBadCommand(t *testing.T) {
	beforeTest()

	hashFilePath = fixturesPath(fmt.Sprintf("file%v.md5", rand.Int()))
	defer os.Remove(hashFilePath)

	caFilePath = fixturesPath(fmt.Sprintf("file%v.pem", rand.Int()))
	defer os.Remove(caFilePath)

	err := safeCAFile(NewErrorCmd(), "valid")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "I am an error")
}

func TestSafeCAFileSuccess(t *testing.T) {
	beforeTest()

	hashFilePath = fixturesPath("tmp.md5")
	defer os.Remove(hashFilePath)

	err := safeCAFile(NewSuccessCmd(), "valid")
	require.Nil(t, err)

	b, _ := os.ReadFile(hashFilePath)

	hash := md5.New()
	io.WriteString(hash, "valid")
	assert.Equal(t, string(b), string(hash.Sum(nil)))

	b, _ = os.ReadFile(caFilePath)
	assert.Equal(t, "valid", string(b))
}
