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
)

func main() {

	log.SetOutput(os.Stderr)

	var credentials Credentials
	if err := readConfiguration(&credentials); err != nil {
		log.Fatalf(err.Error())
	}

	installedProduct, err := getInstalledProduct()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Installed product: %v\n", installedProduct)

	var suseConnectData SUSEConnectData
	if err := readConfiguration(&suseConnectData); err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Registration server set to %v\n", suseConnectData.SccURL)

	products, err := requestProducts(suseConnectData, credentials, installedProduct)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, product := range products {
		dumpRepositories(os.Stdout, product)
	}
}
