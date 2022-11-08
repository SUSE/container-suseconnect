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

// Package regionsrv implements all the utilities to interact with on-deman
// Public clouds.
package regionsrv

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

// ContainerBuildConfig contains all the data that is available through the
// containerbuild-regionsrv server running on the host.
type ContainerBuildConfig struct {
	InstanceData string `json:"instance-data"`
	ServerFqdn   string `json:"server-fqdn"`
	ServerIP     string `json:"server-ip"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Ca           string `json:"ca"`
}

// containerBuildSrvAddress returns a string containing the full address of the TCP
// server. You may tweak this by providing `CONTAINER_BUILD_IP` and/or
// `CONTAINER_BUILD_PORT`, otherwise `0.0.0.0` and `7956` are taken as defaults for
// the IP and the port respectively.
func containerBuildSrvAddress() string {
	ip := os.Getenv("CONTAINER_BUILD_IP")
	if ip == "" {
		ip = "0.0.0.0"
	}

	port := os.Getenv("CONTAINER_BUILD_PORT")
	if port == "" {
		port = "7956"
	}

	return fmt.Sprintf("%s:%s", ip, port)
}

// ServerReachable returns true if the containerbuild-regionsrv server is
// reachable, false otherwise.
func ServerReachable() error {
	addr, err := net.ResolveTCPAddr("tcp", containerBuildSrvAddress())
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

// ReadConfigFromServer performs a request agains the containerbuild-regionsrv
// server running in the host, and it parses the given response so it can be
// used as a ContainerBuildConfig.
func ReadConfigFromServer() (*ContainerBuildConfig, error) {
	addr, err := net.ResolveTCPAddr("tcp", containerBuildSrvAddress())
	log.Printf("Trying to reach suse build server at '%v'", addr.String())
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	log.Printf("Reading from containerbuild-regionsrv ...")

	d := json.NewDecoder(conn)

	data := &ContainerBuildConfig{}
	if err := d.Decode(&data); err != nil {
		return nil, err
	}

	// If something is really bad on the server side, it may return an empty
	// response. Catch this error here.
	if data.InstanceData == "" && data.ServerFqdn == "" &&
		data.ServerIP == "" && data.Ca == "" {
		return nil, errors.New("empty response from the server")
	}
	return data, nil
}
