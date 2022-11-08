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
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var hostsFile = "/etc/hosts"

// UpdateHostsFile updates the hosts file with the given hostname and IP.
func UpdateHostsFile(hostname string, ip string) error {
	content, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		return fmt.Errorf("can't read %s file: %v", hostsFile, err.Error())
	}

	lines := strings.Split(string(content), "\n")
	newcontent := ""
	hostChecked := false
	shorthost := strings.Split(hostname, ".")[0]
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == hostname {
			if fields[0] != ip {
				log.Printf("updating hosts entry for %s", hostname)
				line = fmt.Sprintf("%s %s %s\n", ip, hostname, shorthost)
			}
			hostChecked = true
		}
		newcontent += line + "\n"
	}

	if !hostChecked {
		newcontent += fmt.Sprintf("%s %s %s\n", ip, hostname, shorthost)
	}

	err = ioutil.WriteFile(hostsFile, []byte(newcontent), 0644)
	if err != nil {
		return fmt.Errorf("can't write %s file: %v", hostsFile, err.Error())
	}

	return nil
}
