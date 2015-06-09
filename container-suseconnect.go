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
	"net/url"
	"os"
)

const (
	sccURLStr = "https://scc.suse.com"
)

func main() {

	log.SetOutput(os.Stderr)

	credentials, err := ReadCredentials()
	if err != nil {
		log.Fatalf(err.Error())
	}

	installedProduct, err := getInstalledProduct()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Installed product: %v\n", installedProduct)

	suseConnectData, err := ReadSUSEConnect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if suseConnectData.SccURL == "" {
		suseConnectData.SccURL = sccURLStr
	}
	regURL, err := url.Parse(suseConnectData.SccURL)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Registration server set to %v\n", regURL.String())

	product, err := RequestProduct(*regURL, credentials, installedProduct, suseConnectData.Insecure)
	if err != nil {
		log.Fatalf(err.Error())
	}

	DumpRepositories(os.Stdout, product)
}
