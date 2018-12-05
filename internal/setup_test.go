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
	"regexp"
	"strings"
	"testing"
)

// Handy functions to be used by the test suite.

// Private global value for the tests. It stores all the contents that have
// been logged after a `prepareLogger` call.
var logged *bytes.Buffer

// It initializes the logger infrastructure for tests.
func prepareLogger() {
	logged = bytes.NewBuffer([]byte{})
	log.SetOutput(logged)
}

// Make sure that the logged string matches the given expected string.
func shouldHaveLogged(t *testing.T, str string) {
	original := logged.String()
	if strings.TrimSpace(original) == "" {
		t.Fatal("Nothing has been logged.\n")
	}

	// The logged string includes the timestamp, get rid of it.
	re := regexp.MustCompile("^\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}\\s")
	logStr := strings.TrimSpace(re.ReplaceAllString(original, ""))

	if strings.TrimSpace(str) != logStr {
		t.Fatalf("Should have logged: %v, not %v\n", str, logStr)
	}
}
