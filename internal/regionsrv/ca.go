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
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	oldHashFilePath = "/etc/pki/containerbuild-regionsrv.md5"
	hashFilePath    = "/etc/pki/containerbuild-regionsrv.sha256"
	caFilePath      = "/etc/pki/trust/anchors/containerbuild-regionsrv.pem"
)

// commander is a very simple interface that just implements the `Run` function,
// which returns an error. This interface has merely been introduced to ease up
// testing.
type commander interface {
	Run() error
}

// Returns true if the CA file needs an update, false otherwise.
func updateNeeded(contents string) bool {
	if _, err := os.Stat(hashFilePath); os.IsNotExist(err) {
		return true
	}

	data, err := os.ReadFile(hashFilePath)
	if err != nil {
		return true
	}

	hash := sha256.New()
	io.WriteString(hash, contents)

	return strings.TrimSpace(string(data)) != string(hash.Sum(nil))
}

// saveCAFile implements `SaveCAFile` by assuming a `commander` type will be
// given.
func saveCAFile(cmd commander, contents string) error {
	if !updateNeeded(contents) {
		return nil
	}

	// Nuke everything before populating things back again.
	os.Remove(oldHashFilePath)
	os.Remove(hashFilePath)
	os.Remove(caFilePath)

	// Save the file
	err := os.WriteFile(caFilePath, []byte(contents), 0o644)
	if err != nil {
		return err
	}

	// Execute `update-ca-certificates` now.
	if err = cmd.Run(); err != nil {
		return err
	}

	// Save the new checksum
	hash := sha256.New()
	io.WriteString(hash, contents)
	os.WriteFile(hashFilePath, hash.Sum(nil), 0o644)

	return nil
}

// SaveCAFile creates a certificate file into the right location if it isn't
// already there. This function will call `update-ca-certificates` whenever the
// CA file has been updated.
func SaveCAFile(contents string) error {
	cmd := exec.Command("update-ca-certificates")
	return saveCAFile(cmd, contents)
}
