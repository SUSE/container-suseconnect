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

package containersuseconnect

import (
	"fmt"
	"runtime/debug"
	"strconv"
)

var (
	version  string
	revision string
)

// GetVersion returns the version and revision if available.
// If both version and revision are found, it returns "%version (%revision)".
// If only version is found, it returns "%version".
// If no version is found but a revision is, it returns "devel (%revision)".
// If no version and no revision are found, it returns "devel".
// The revision can taken from the binary if it was stamped with version
// control information.
func GetVersion() string {
	if version == "" {
		version = "devel"
	}

	if revision == "" {
		bi, ok := debug.ReadBuildInfo()

		if !ok {
			return version
		}

		var vcsRevision string
		var vcsModified bool

		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				vcsRevision = s.Value
			case "vcs.modified":
				vcsModified, _ = strconv.ParseBool(s.Value)
			}
		}

		if vcsRevision == "" {
			return version
		} else {
			revision = vcsRevision[:7]

			if vcsModified {
				revision = revision + "+modified"
			}
		}
	}

	return fmt.Sprintf("%s (%s)", version, revision)
}
