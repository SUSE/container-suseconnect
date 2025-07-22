// Copyright (c) 2025 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package containersuseconnect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionWithVersionRevisionUnset(t *testing.T) {
	version = ""
	revision = ""

	assert.Equal(t, "devel", GetVersion())
}

func TestGetVersionWithVersionSet(t *testing.T) {
	version = "1.2.3"
	revision = ""

	assert.Equal(t, "1.2.3", GetVersion())
}

func TestGetVersionWithRevisionSet(t *testing.T) {
	version = ""
	revision = "gh12345"

	assert.Equal(t, "devel (gh12345)", GetVersion())
}

func TestGetVersionWithVersionRevisionSet(t *testing.T) {
	version = "1.2.3"
	revision = "gh12345"

	assert.Equal(t, "1.2.3 (gh12345)", GetVersion())
}
