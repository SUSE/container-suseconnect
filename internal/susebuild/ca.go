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

package susebuild

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	hashFilePath = "/etc/pki/susebuild.md5"
	caFilePath   = "/etc/pki/trust/anchors/susebuild.pem"
)

// commander is a very simple interface that just implements the `Run` function,
// which returns an error. This interface has merely been introduced to ease up
// testing.
type commander interface {
	Run() error
}

// Returns true if the CA file has to be updated, false otherwise.
func updateNeeded(contents string) bool {
	if _, err := os.Stat(hashFilePath); os.IsNotExist(err) {
		return true
	}

	data, err := ioutil.ReadFile(hashFilePath)
	if err != nil {
		return true
	}
	sum := strings.TrimSpace(string(data))

	hash := md5.New()
	io.WriteString(hash, contents)

	return sum != string(hash.Sum(nil))
}

// safeCAFile implements `SafeCAFile` by assuming a `commander` type will be
// given.
func safeCAFile(cmd commander, contents string) error {
	if !updateNeeded(contents) {
		return nil
	}

	// Nuke everything before populating things back again.

	_ = os.Remove(hashFilePath)
	_ = os.Remove(caFilePath)

	// Safe the file
	err := ioutil.WriteFile(caFilePath, []byte(contents), 0644)
	if err != nil {
		return err
	}

	// Execute `update-ca-certificates` now.
	if err = cmd.Run(); err != nil {
		return err
	}

	// Safe the new checksum
	hash := md5.New()
	io.WriteString(hash, contents)
	_ = ioutil.WriteFile(hashFilePath, hash.Sum(nil), 0644)

	return nil
}

// SafeCAFile creates a certificate file into the right location if it isn't
// already there. This function will call `update-ca-certificates` whenever the
// CA file has been updated.
func SafeCAFile(contents string) error {
	cmd := exec.Command("update-ca-certificates")
	return safeCAFile(cmd, contents)
}
