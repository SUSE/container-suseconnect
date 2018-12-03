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
	"errors"
	"os"
	"strings"
	"testing"
)

func TestGetLocationPath(t *testing.T) {
	path := getLocationPath([]string{})
	if path != "" {
		t.Fatal("It should be empty")
	}

	strs := []string{
		"does/not/exist",
		"../../test/products-sle12.json",
	}
	path = getLocationPath(strs)
	if path != "../../test/products-sle12.json" {
		t.Fatalf("Wrong location path: %v", path)
	}
}

type NotFoundConfiguration struct{}

func (cfg NotFoundConfiguration) separator() byte {
	return '.'
}

func (cfg NotFoundConfiguration) locations() []string {
	return []string{}
}

func (cfg NotFoundConfiguration) onLocationsNotFound() bool {
	return false
}

func (cfg NotFoundConfiguration) setValues(key, value string) {
}

func (cfg NotFoundConfiguration) afterParseCheck() error {
	return nil
}

func TestNotFound(t *testing.T) {
	var cfg NotFoundConfiguration

	prepareLogger()
	err := ReadConfiguration(&cfg)
	if err == nil || err.Error() != "No locations found: []" {
		t.Fatalf("Wrong error: %v", err)
	}
	shouldHaveLogged(t, "No locations found: []")
}

type NotAllowedConfiguration struct{}

func (cfg NotAllowedConfiguration) separator() byte {
	return '.'
}

func (cfg NotAllowedConfiguration) locations() []string {
	return []string{"/etc/shadow"}
}

func (cfg NotAllowedConfiguration) onLocationsNotFound() bool {
	return false
}

func (cfg NotAllowedConfiguration) setValues(key, value string) {
}

func (cfg NotAllowedConfiguration) afterParseCheck() error {
	return nil
}

func TestNotAllowed(t *testing.T) {
	var cfg NotAllowedConfiguration

	prepareLogger()
	err := ReadConfiguration(&cfg)
	msg := "Can't open /etc/shadow file: open /etc/shadow: permission denied"
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}
	shouldHaveLogged(t, msg)
}

func TestParseInvalid(t *testing.T) {
	var cfg NotAllowedConfiguration

	file, err := os.Open("/etc/shadow")
	if err == nil {
		file.Close()
		t.Fatal("There should be an error here")
	}
	prepareLogger()
	err = parse(cfg, file)
	msg := "Error when scanning configuration: invalid argument"
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}
	shouldHaveLogged(t, msg)
}

type ErrorAfterParseConfiguration struct{}

func (cfg ErrorAfterParseConfiguration) separator() byte {
	return '.'
}

func (cfg ErrorAfterParseConfiguration) locations() []string {
	return []string{}
}

func (cfg ErrorAfterParseConfiguration) onLocationsNotFound() bool {
	return false
}

func (cfg ErrorAfterParseConfiguration) setValues(key, value string) {
}

func (cfg ErrorAfterParseConfiguration) afterParseCheck() error {
	return errors.New("I'm grumpy, and I want to error!")
}

func TestParseFailAfterCheck(t *testing.T) {
	var cfg ErrorAfterParseConfiguration

	str := strings.NewReader("")
	err := parse(cfg, str)
	if err == nil || err.Error() != "I'm grumpy, and I want to error!" {
		t.Fatal("Wrong error")
	}
}

func TestParseFailNoSeparator(t *testing.T) {
	var cfg ErrorAfterParseConfiguration

	str := strings.NewReader("keywithoutvalue")
	prepareLogger()
	err := parse(cfg, str)
	msg := "Can't parse line: keywithoutvalue"
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}
	shouldHaveLogged(t, msg)
}
