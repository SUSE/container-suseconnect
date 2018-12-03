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

package container_suseconnect

import (
	"os"
	"testing"
)

// The helper for testing the logger setup where:
//
//   - env: the value for the `logEnv` environment value.
//   - expected: the expected path of the logger.
//   - cleanup: whether the file has to be closed and removed.
func testLogger(t *testing.T, env, expected string, cleanup bool) {
	os.Setenv(logEnv, env)
	f := GetLoggerFile()
	if f.Name() != expected {
		t.Fatalf("Wrong file")
	}

	if cleanup {
		f.Close()
		err := os.Remove(f.Name())
		if err != nil {
			t.Fatalf("Problem when cleaning up")
		}
	}
}

func TestSetupLoggerDefault(t *testing.T) {
	defaultLogPath = "suseconnect.log"
	testLogger(t, "", defaultLogPath, true)
}

func TestSetupLoggerCustom(t *testing.T) {
	testLogger(t, "suse.log", "suse.log", true)
}

func TestLoggerStdErr(t *testing.T) {
	testLogger(t, "/var/log/suse.log", os.Stderr.Name(), false)
}
