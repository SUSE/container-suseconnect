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
	"log"
	"os"
)

// Credentials holds the host credentials
// NOTE (@mssola): in SCC we have introduced the System-Token credential. For
// now this is not affecting the normal operation of this application, but it
// is something to keep in mind for further developments.
type Credentials struct {
	Username     string
	Password     string
	SystemToken  string
	InstanceData string
}

func (cr *Credentials) separator() byte {
	return '='
}

func (cr *Credentials) locations() []string {
	return []string{
		"/etc/zypp/credentials.d/SCCcredentials",
		"/run/secrets/SCCcredentials",
		"/run/secrets/credentials.d/SCCcredentials",
	}
}

func (cr *Credentials) onLocationsNotFound() bool {
	env_user := os.Getenv("SCC_CREDENTIAL_USERNAME")
	env_pass := os.Getenv("SCC_CREDENTIAL_PASSWORD")
	env_system_token := os.Getenv("SCC_CREDENTIAL_SYSTEM_TOKEN")

	if env_user != "" && env_pass != "" {
		cr.Username = env_user
		cr.Password = env_pass
		cr.SystemToken = env_system_token
		return true
	}

	return false
}

func (cr *Credentials) setValues(key, value string) {
	switch key {
	case "username":
		cr.Username = value
	case "password":
		cr.Password = value
	case "system_token":
		cr.SystemToken = value
	default:
		log.Printf("Warning: Unknown key '%v'", key)
	}
}

func (cr *Credentials) afterParseCheck() error {
	if cr.Username == "" {
		return loggedError(InvalidCredentialsError, "Can't find username")
	}

	if cr.Password == "" {
		return loggedError(InvalidCredentialsError, "Can't find password")
	}

	return nil
}
