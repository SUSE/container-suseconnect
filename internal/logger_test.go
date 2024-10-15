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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLogWritterFromRegularFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "prefix")
	require.Nil(t, err)

	defer os.Remove(tempFile.Name())

	w, err := getLogWritter(tempFile.Name())
	assert.Nil(t, err)
	assert.NotNil(t, w)
	w.Close()
}

func TestGetLogWritterFromRelativeFile(t *testing.T) {
	w, err := getLogWritter("test.log")
	assert.EqualError(t, err, "log path is not absulute: test.log")
	assert.Nil(t, w)
}

func TestGetLogWritterFromValidDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	w, err := getLogWritter(dir)
	assert.EqualError(t, err, fmt.Sprintf("open %s: is a directory", dir))
	assert.Nil(t, w)
}

func TestGetLogWritterFromMissingDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	// ensure the path is terminated by / othewise it will be considered a file
	path := filepath.Join(dir, "does", "not", "exits") + string(filepath.Separator)
	w, err := getLogWritter(path)
	assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", path))
	assert.Nil(t, w)
}

func TestGetLogWritterFromMissingDirFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "does", "not", "exits", "test.log")
	w, err := getLogWritter(path)
	assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", path))
	assert.Nil(t, w)
}

func TestGetLogWritterFromEmptyString(t *testing.T) {
	w, err := getLogWritter("")
	assert.EqualError(t, err, "path is empty")
	assert.Nil(t, w)
}

func TestGetLogWritterFromSpaces(t *testing.T) {
	w, err := getLogWritter("    ")
	assert.EqualError(t, err, "path is empty")
	assert.Nil(t, w)
}

func TestGetLogWritterFromSymbolInName(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	w, err := getLogWritter(filepath.Join(dir, "@"))
	assert.Nil(t, err)
	assert.NotNil(t, w)
	w.Close()
}

func TestGetLogWritterFromDevice(t *testing.T) {
	w, err := getLogWritter("/dev/null")
	assert.EqualError(t, err, "path is not a regular file: /dev/null")
	assert.Nil(t, w)
}

func TestGetLogWritterFromNonWritableFile(t *testing.T) {
	// /proc/1/mem is a valid regular file, but
	// it is not writable by any user, even root
	w, err := getLogWritter("/proc/1/mem")
	// it can fail with various errors
	// not worth checking for an exact message
	assert.NotNil(t, err)
	assert.Nil(t, w)
}

func TestGetLogEnvIfNotSet(t *testing.T) {
	// ensure no variable is set
	os.Unsetenv(LogEnv)
	path := getLogEnv()
	assert.Empty(t, path)
}

func TestGetLogEnvIfSet(t *testing.T) {
	// ensure no variable is set
	os.Unsetenv(LogEnv)
	err := os.Setenv(LogEnv, "/path/file.log")
	require.Nil(t, err)
	defer os.Unsetenv(LogEnv)

	envPath := getLogEnv()
	assert.Equal(t, "/path/file.log", envPath)
}

func TestGetLogEnvTrimSpaces(t *testing.T) {
	// ensure no variable is set
	os.Unsetenv(LogEnv)
	err := os.Setenv(LogEnv, "    /path/file.log    ")
	require.Nil(t, err)
	defer os.Unsetenv(LogEnv)

	envPath := getLogEnv()
	assert.Equal(t, "/path/file.log", envPath)
}

// Ensures that the log is always written to a file and Stderr.
func TestSetLoggerOutput(t *testing.T) {
	// ensure no variable is set
	os.Unsetenv(LogEnv)

	tempFile, err := os.CreateTemp("", "")
	require.Nil(t, err)
	defer os.Remove(tempFile.Name())

	err = os.Setenv(LogEnv, tempFile.Name())
	require.Nil(t, err)
	defer os.Unsetenv(LogEnv)

	logLine := "This in a log entry in a file and Stderr"

	var stdData, fileData string

	stdData, err = captureStderr(t, func() {
		SetLoggerOutput()
		log.Println(logLine)
	})

	assert.Nil(t, err, "Failed to capture Stderr")

	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(tempFile)
	assert.Nil(t, err, "Failed to read temp file")

	fileData = buff.String()

	assert.Contains(t, fileData, logLine)
	assert.Contains(t, stdData, logLine)
}
