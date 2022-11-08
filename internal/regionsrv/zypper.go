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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

// ParseStdin parses the standard input as given by zypper and returns a map
// with the parsed parameters.
func ParseStdin() (map[string]string, error) {
	params := make(map[string]string)
	first := true

	// The zypper plugin protocol is based on STOMP. STOMP messages
	// are NUL-terminated. So read the entire message first
	reader := bufio.NewReader(os.Stdin)
	msg, err := reader.ReadBytes(0)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Now read the message line by line. URL resolver plugin messages
	// are just <key>:<value> header lines and don't contain a body.
	sr := bytes.NewReader(msg)
	scanner := bufio.NewScanner(sr)
	for scanner.Scan() {
		if first {
			first = false
		} else {
			vals := strings.SplitN(scanner.Text(), ":", 2)
			if len(vals) == 2 {
				params[vals[0]] = vals[1]
			}
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return params, nil
}

// PrintResponse prints to standard output with the format expected by zypper.
func PrintResponse(params map[string]string) error {
	cfg, err := ReadConfigFromServer()
	if err != nil {
		return err
	}

	// Error out if we have no information on the credentials.
	if cfg.Username == "" && cfg.Password == "" {
		return errors.New("no credentials given")
	}

	// Safe the contents of the CA file if it doesn't exist already.
	if err = SafeCAFile(cfg.Ca); err != nil {
		return err
	}

	printFromConfiguration(params["path"], cfg)
	return nil
}

func printFromConfiguration(path string, cfg *ContainerBuildConfig) {
	u := url.URL{
		Scheme: "https",
		Host:   cfg.ServerFqdn,
		Path:   path,
		User:   url.UserPassword(cfg.Username, "XXXX"),
	}

	log.Print("Received X-Instance-Data")
	log.Printf("Resulting URL: %s", u.String())

	// Add user info to URL to avoid password appearing in logs
	u.User = url.UserPassword(cfg.Username, cfg.Password)

	fmt.Printf("RESOLVEDURL\n")
	// Add an extra emptyline to separate Headers from payload
	fmt.Printf("X-Instance-Data:%s\n\n", cfg.InstanceData)
	// Message needs to be NUL-terminated
	fmt.Printf("%s\000", u.String())

}
