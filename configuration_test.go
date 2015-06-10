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

package main

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
		"data/product.json",
	}
	path = getLocationPath(strs)
	if path != "data/product.json" {
		t.Fatal("Wrong location path")
	}
}

type NotFoundConfiguration struct{}

func (cfg NotFoundConfiguration) separator() byte {
	return '.'
}

func (cfg NotFoundConfiguration) locations() []string {
	return []string{}
}

func (cfg NotFoundConfiguration) setValues(key, value string) {
}

func (cfg NotFoundConfiguration) afterParseCheck() error {
	return nil
}

func TestNotFound(t *testing.T) {
	var cfg NotFoundConfiguration

	err := read(&cfg)
	if err == nil || err.Error() != "No locations found!" {
		t.Fatal("Wrong error")
	}
}

type NotAllowedConfiguration struct{}

func (cfg NotAllowedConfiguration) separator() byte {
	return '.'
}

func (cfg NotAllowedConfiguration) locations() []string {
	return []string{"/etc/shadow"}
}

func (cfg NotAllowedConfiguration) setValues(key, value string) {
}

func (cfg NotAllowedConfiguration) afterParseCheck() error {
	return nil
}

func TestNotAllowed(t *testing.T) {
	var cfg NotAllowedConfiguration

	err := read(&cfg)
	if err == nil || err.Error() != "Can't open /etc/shadow file: open /etc/shadow: permission denied" {
		t.Fatal("Wrong error")
	}
}

func TestParseInvalid(t *testing.T) {
	var cfg NotAllowedConfiguration

	file, err := os.Open("/etc/shadow")
	if err == nil {
		file.Close()
		t.Fatal("There should be an error here")
	}
	err = parse(cfg, file)
	if err == nil || err.Error() != "Error when scanning configuration: invalid argument" {
		t.Fatal("Wrong error")
	}
}

type ErrorAfterParseConfiguration struct{}

func (cfg ErrorAfterParseConfiguration) separator() byte {
	return '.'
}

func (cfg ErrorAfterParseConfiguration) locations() []string {
	return []string{}
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
	err := parse(cfg, str)
	if err == nil || err.Error() != "Can't parse line: keywithoutvalue" {
		t.Fatal("Wrong error")
	}
}
