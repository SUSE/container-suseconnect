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
	"log"
)

const (
	sccURLStr = "https://scc.suse.com"
)

type SUSEConnectData struct {
	SccURL   string
	Insecure bool
}

func (data *SUSEConnectData) separator() byte {
	return ':'
}

func (data *SUSEConnectData) locations() []string {
	return []string{"/etc/SUSEConnect", "/run/secrets/SUSEConnect"}
}

func (data *SUSEConnectData) onLocationsNotFound() bool {
	data.SccURL = sccURLStr
	return true
}

func (data *SUSEConnectData) setValues(key, value string) {
	if key == "url" {
		data.SccURL = value
	} else if key == "insecure" {
		data.Insecure = value == "true"
	} else {
		log.Printf("Warning: Unknown key '%v'", key)
	}
}

func (data *SUSEConnectData) afterParseCheck() error {
	if data.SccURL == "" {
		data.SccURL = sccURLStr
	}
	return nil
}
