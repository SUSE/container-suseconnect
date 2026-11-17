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

func TestGetLogWriterFromRegularFile(t *testing.T) {
	tempFile, err := os.CreateTemp(t.TempDir(), "prefix")
	require.NoError(t, err)
	defer tempFile.Close()

	w, err := getLogWriter(tempFile.Name())
	require.NoError(t, err)
	assert.NotNil(t, w)
	w.Close()
}

func TestGetLogWriterFromRelativeFile(t *testing.T) {
	w, err := getLogWriter("test.log")
	assert.EqualError(t, err, "log path is not absolute: test.log")
	assert.Nil(t, w)
}

func TestGetLogWriterFromValidDir(t *testing.T) {
	dir := t.TempDir()

	w, err := getLogWriter(dir)
	assert.EqualError(t, err, fmt.Sprintf("open %s: is a directory", dir))
	assert.Nil(t, w)
}

func TestGetLogWriterFromMissingDir(t *testing.T) {
	dir := t.TempDir()

	// ensure the path is terminated by / othewise it will be considered a file
	path := filepath.Join(dir, "does", "not", "exits") + string(filepath.Separator)
	w, err := getLogWriter(path)
	assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", path))
	assert.Nil(t, w)
}

func TestGetLogWriterFromMissingDirFile(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "does", "not", "exits", "test.log")
	w, err := getLogWriter(path)
	assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", path))
	assert.Nil(t, w)
}

func TestGetLogWriterFromEmptyString(t *testing.T) {
	w, err := getLogWriter("")
	assert.EqualError(t, err, "path is empty")
	assert.Nil(t, w)
}

func TestGetLogWriterFromSpaces(t *testing.T) {
	w, err := getLogWriter("    ")
	assert.EqualError(t, err, "path is empty")
	assert.Nil(t, w)
}

func TestGetLogWriterFromSymbolInName(t *testing.T) {
	dir := t.TempDir()

	w, err := getLogWriter(filepath.Join(dir, "@"))
	require.NoError(t, err)
	assert.NotNil(t, w)
	w.Close()
}

func TestGetLogWriterFromDevice(t *testing.T) {
	w, err := getLogWriter("/dev/null")
	assert.EqualError(t, err, "path is not a regular file: /dev/null")
	assert.Nil(t, w)
}

func TestGetLogWriterFromNonWritableFile(t *testing.T) {
	// /proc/1/mem is a valid regular file, but
	// it is not writable by any user, even root
	w, err := getLogWriter("/proc/1/mem")
	// it can fail with various errors
	// not worth checking for an exact message
	assert.NotNil(t, err)
	assert.Nil(t, w)
}

func TestGetLogEnvIfNotSet(t *testing.T) {
	t.Setenv(LogEnv, "")
	path := getLogEnv()
	assert.Empty(t, path)
}

func TestGetLogEnvIfSet(t *testing.T) {
	t.Setenv(LogEnv, "/path/file.log")

	envPath := getLogEnv()
	assert.Equal(t, "/path/file.log", envPath)
}

func TestGetLogEnvTrimSpaces(t *testing.T) {
	t.Setenv(LogEnv, "    /path/file.log    ")

	envPath := getLogEnv()
	assert.Equal(t, "/path/file.log", envPath)
}

// Ensures that the log is always written to a file and Stderr.
func TestSetLoggerOutput(t *testing.T) {
	tempFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	defer tempFile.Close()

	t.Setenv(LogEnv, tempFile.Name())

	logLine := "This in a log entry in a file and Stderr"

	var stdData, fileData string

	stdData, err = captureStderr(t, func() {
		SetLoggerOutput()
		log.Println(logLine)
	})

	require.NoError(t, err, "Failed to capture Stderr")

	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(tempFile)
	require.NoError(t, err, "Failed to read temp file")

	fileData = buff.String()

	assert.Contains(t, fileData, logLine)
	assert.Contains(t, stdData, logLine)
}
