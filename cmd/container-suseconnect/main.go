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

// Gives access to repositories during docker build and run using the host
// machine credentials.
package main

import (
	"log"
	"os"

	cs "github.com/SUSE/container-suseconnect/pkg/container-suseconnect"
)

func main() {
	log.SetOutput(cs.GetLoggerFile())

	var credentials cs.Credentials
	if err := cs.ReadConfiguration(&credentials); err != nil {
		log.Fatalf(err.Error())
	}

	installedProduct, err := cs.GetInstalledProduct()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Installed product: %v\n", installedProduct)

	var suseConnectData cs.SUSEConnectData
	if err := cs.ReadConfiguration(&suseConnectData); err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Registration server set to %v\n", suseConnectData.SccURL)

	products, err := cs.RequestProducts(suseConnectData, credentials, installedProduct)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, product := range products {
		cs.DumpRepositories(os.Stdout, product)
	}
}
