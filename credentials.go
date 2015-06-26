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
	"fmt"
	"log"
)

// Credentials holds the host credentials
type Credentials struct {
	Username string
	Password string
}

func (cr *Credentials) separator() byte {
	return '='
}

func (cr *Credentials) locations() []string {
	return []string{
		"/etc/zypp/credentials.d/SCCcredentials",
		"/run/secrets/credentials.d/SCCcredentials",
	}
}

func (cr *Credentials) onLocationsNotFound() bool {
	return false
}

func (cr *Credentials) setValues(key, value string) {
	if key == "username" {
		cr.Username = value
	} else if key == "password" {
		cr.Password = value
	} else {
		log.Printf("Warning: Unknown key '%v'", key)
	}
}

func (cr *Credentials) afterParseCheck() error {
	if cr.Username == "" {
		return fmt.Errorf("Can't find username")
	}
	if cr.Password == "" {
		return fmt.Errorf("Can't find password")
	}
	return nil
}
