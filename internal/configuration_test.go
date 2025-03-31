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
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfiguration struct {
	loc  []string
	fail bool
}

func (cfg testConfiguration) separator() byte {
	return '.'
}

func (cfg testConfiguration) locations() []string {
	return cfg.loc
}

func (cfg testConfiguration) onLocationsNotFound() bool {
	return false
}

func (cfg testConfiguration) setValues(key, value string) {
}

func (cfg testConfiguration) afterParseCheck() error {
	if cfg.fail {
		return errors.New("I'm grumpy, and I want to error")
	}

	return nil
}

func TestGetLocationPath(t *testing.T) {
	path := getLocationPath([]string{})
	assert.Empty(t, path)

	strs := []string{
		"does/not/exist",
		"testdata/products-sle12.json",
	}

	path = getLocationPath(strs)
	assert.Equal(t, "testdata/products-sle12.json", path)
}

func TestGetLocationPathEmpty(t *testing.T) {
	path := getLocationPath([]string{})
	assert.Empty(t, path)
}

func TestReadConfigurationNotFound(t *testing.T) {
	prepareLogger()

	cfg := testConfiguration{
		loc: []string{"/path/not/found/cfg"},
	}

	err := ReadConfiguration(&cfg)
	require.NotNil(t, err)

	msg := "Warning: SUSE credentials not found: [/path/not/found/cfg] - automatic handling of repositories not done."
	assert.EqualError(t, err, msg)
	shouldHaveLogged(t, msg)
}

func TestParseInvalid(t *testing.T) {
	prepareLogger()

	var cfg testConfiguration
	str := strings.NewReader("@")
	err := parse(cfg, str)
	msg := "Can't parse line: @"
	assert.EqualError(t, err, msg)

	shouldHaveLogged(t, msg)
}

func TestParseFailAfterCheck(t *testing.T) {
	cfg := testConfiguration{
		fail: true,
	}

	str := strings.NewReader("")
	err := parse(cfg, str)
	msg := "I'm grumpy, and I want to error"
	assert.EqualError(t, err, msg)
}

func TestParseFailNoSeparator(t *testing.T) {
	prepareLogger()

	var cfg testConfiguration
	str := strings.NewReader("keywithoutvalue")
	err := parse(cfg, str)
	assert.NotNil(t, err)

	msg := "Can't parse line: keywithoutvalue"
	assert.EqualError(t, err, msg)

	shouldHaveLogged(t, msg)
}
