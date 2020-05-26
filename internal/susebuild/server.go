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

// Utilities for talking to a suse-build server.
package susebuild

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

// SuseBuildConfig contains all the data that is available through the
// suse-build server running on the host.
type SuseBuildConfig struct {
	InstanceData string `json:"instance-data"`
	ServerFqdn   string `json:"server-fqdn"`
	ServerIp     string `json:"server-ip"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Ca           string `json:"ca"`
}

// suseBuildAddress returns a string containing the full address of the TCP
// server. You may tweak this by providing `SUSE_BUILD_IP` and/or
// `SUSE_BUILD_PORT`, otherwise `0.0.0.0` and `7956` are taken as defaults for
// the IP and the port respectively.
func suseBuildAddress() string {
	ip := os.Getenv("SUSE_BUILD_IP")
	if ip == "" {
		ip = "0.0.0.0"
	}

	port := os.Getenv("SUSE_BUILD_PORT")
	if port == "" {
		port = "7956"
	}

	return fmt.Sprintf("%s:%s", ip, port)
}

// ServerAvailable returns true if the suse-build server is reachable, false
// otherwise.
func ServerReachable() error {
	addr, err := net.ResolveTCPAddr("tcp", suseBuildAddress())
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	_ = conn.Close()
	return nil
}

// ReadConfigFromServer performs a request agains the suse-build server running
// in the host, and it parses the given response so it can be used as a
// SuseBuildConfig.
func ReadConfigFromServer() (*SuseBuildConfig, error) {
	addr, err := net.ResolveTCPAddr("tcp", suseBuildAddress())
	log.Printf("Trying to reach suse build server at '%v'", addr.String())
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	log.Printf("Reading from suse build server...")

	// After testing it turns out that we need something a bit over 2048, but
	// let's leave some extra room just in case...
	reply := make([]byte, 8192)

	n, err := conn.Read(reply)
	if err != nil {
		return nil, err
	}

	data := &SuseBuildConfig{}
	if err := json.Unmarshal(reply[:n], &data); err != nil {
		return nil, err
	}
	return data, nil
}
