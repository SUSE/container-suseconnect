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

// Zypper url resolver plugin for /susecloud schemes.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	cs "github.com/SUSE/container-suseconnect/internal"
	"github.com/SUSE/container-suseconnect/internal/susebuild"
	"github.com/urfave/cli/v2"
)

const regionServerConfig = "/etc/regionserverclnt.cfg"

func main() {
	app := cli.NewApp()
	app.Copyright = fmt.Sprintf("(c) %d SUSE LCC", time.Now().Year())
	app.Name = "susecloud"
	app.Version = cs.Version
	app.Usage = ""
	app.UsageText =
		`Zypper URL resolver plugin that transforms plugin:/susecloud urls into
	their proper counterparts inside of a SUSE container.`
	app.Action = runZypperUrlResolver

	log.SetOutput(cs.GetLoggerFile())
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// The main action: reads the arguments as given by zypper on the stdin, and
// prints into stdout the response to be used.
func runZypperUrlResolver(_ *cli.Context) error {
	if err := susebuild.ServerReachable(); err != nil {
		return fmt.Errorf("Could not reach build server from the host: %v", err)
	}

	input, err := parseStdin()
	if err != nil {
		return fmt.Errorf("Could not parse input: %s", err)
	}

	return printResponse(input)
}

// Parses the standard input as given by zypper and returns a map with the
// parsed parameters.
func parseStdin() (map[string]string, error) {
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

// Prints to standard output with the format expected by zypper.
func printResponse(params map[string]string) error {
	cfg, err := susebuild.ReadConfigFromServer()
	if err != nil {
		return err
	}

	u, err := url.Parse(cfg.Server)
	if err != nil {
		return err
	}

	q := u.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	log.Printf("Resulting X-Instance-Data: %s", cfg.InstanceData)
	log.Printf("Resulting URL: %s", u.String())

	fmt.Printf("RESOLVEDURL\n")
	// Add an extra emptyline to separate Headers from payload
	fmt.Printf("X-Instance-Data:%s\n\n", cfg.InstanceData)
	// Message needs to be NUL-terminated
	fmt.Printf("%s\000", u.String())

	return nil
}
