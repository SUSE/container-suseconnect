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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Default path for the log file.
const DefaultLogPath = "/var/log/suseconnect.log"

// Environment variable used to specify a custom path for the log file.
const LogEnv = "SUSECONNECT_LOG_FILE"

// getLogEnv returns the value set to the [LogEnv] environment variable.
func getLogEnv() string {
	return strings.TrimSpace(os.Getenv(LogEnv))
}

// getLogWritter checks if the path can be open and written to
// and returns an [io.WriteCloser] if there are no errors.
func getLogWritter(path string) (io.WriteCloser, error) {
	path = strings.TrimSpace(path)

	if len(path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}

	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("log path is not absulute: %s", path)
	}

	lf, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)

	if err != nil {
		return nil, err
	}

	if fi, err := lf.Stat(); err == nil {
		if !fi.Mode().IsRegular() {
			lf.Close()
			return nil, fmt.Errorf("path is not a regular file: %s", path)
		}
	}

	_, err = lf.WriteString(fmt.Sprintf("container-suseconnect %s\n", Version))

	if err != nil {
		lf.Close()
		return nil, err
	}

	return lf, nil
}

// SetLoggerOutput configures the logger to write to /dev/stderr
// and to a file.
//
// If [LogEnv] is set and writable it writes to the file defined,
// otherwise it writes to [DefaultLogPath].
func SetLoggerOutput() {
	// ensure we are logging to stderr and nowhere else
	log.SetOutput(os.Stderr)

	path := getLogEnv()

	if len(path) == 0 {
		path = DefaultLogPath
	}

	w, err := getLogWritter(path)

	if err == nil {
		writter := io.MultiWriter(os.Stderr, w)
		log.SetOutput(writter)
		log.Printf("Log file location: %s\n", path)
	} else {
		log.Printf("Failed to set up log file '%s'\n", path)
		log.Println(err)
	}
}

// Log the given formatted string with its parameters, and return it
// as a new error.
func loggedError(errorCode int, format string, params ...interface{}) *SuseConnectError {
	msg := fmt.Sprintf(format, params...)
	log.Print(msg)
	return &SuseConnectError{
		ErrorCode: errorCode,
		message:   msg,
	}
}
