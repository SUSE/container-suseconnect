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
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var credentialLocations = []string{
	"/etc/zypp/credentials.d/SCCcredentials",
	"/run/secrets/credentials.d/SCCcredentials",
}

// Credentials holds the host credentials
type Credentials struct {
	Username string
	Password string
}

// ParseCredentials parse the contents of the credentials file and returns a Credentials instance
func ParseCredentials(reader io.Reader) (Credentials, error) {
	credentials := Credentials{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// comments #
		if strings.Index(scanner.Text(), "#") == 0 {
			continue
		}

		// empty lines
		if scanner.Text() == "" {
			continue
		}

		parts := strings.SplitN(scanner.Text(), "=", 2)
		if len(parts) != 2 {
			return Credentials{}, fmt.Errorf("Can't parse line: %v", scanner.Text())
		}

		if parts[0] == "username" {
			credentials.Username = parts[1]
		}

		if parts[0] == "password" {
			credentials.Password = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return credentials, err
	}

	if credentials.Username == "" {
		return Credentials{}, fmt.Errorf("Can't find username")
	}

	if credentials.Password == "" {
		return Credentials{}, fmt.Errorf("Can't find password")
	}

	return credentials, nil
}

// ReadCredentials looks for a credential file (first inside of /etc/zypp/credentials.d/, then inside of /run/secrets/credentials.d)
// and returns a Credentials instance
func ReadCredentials() (Credentials, error) {
	var credentialsPath string
	for _, path := range credentialLocations {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else {
			credentialsPath = path
			break
		}
	}

	if credentialsPath == "" {
		return Credentials{}, fmt.Errorf("No credentials found")
	}

	credFile, err := os.Open(credentialsPath)
	if err != nil {
		return Credentials{}, fmt.Errorf("Can't open credentials file: %v", err.Error())
	}
	defer credFile.Close()

	return ParseCredentials(credFile)
}
