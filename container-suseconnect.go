//
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
//
package main

import (
	"log"
	"os"
	"net/url"
)

const (
	sccUrlStr = "https://scc.suse.com"
)

func main () {

	log.SetOutput(os.Stderr)

	credentials, err := ReadCredentials()
	if err != nil {
		log.Fatalf(err.Error())
	}

	installedProduct, err := ReadInstalledProduct()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Installed product: %v\n", installedProduct)

	regUrlStr := os.Getenv("SCC_URL")
	if regUrlStr == "" {
		regUrlStr = sccUrlStr
	}
	regUrl, err := url.Parse(regUrlStr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Registration server set to %v\n", regUrl.String())

	product, err := RequestProduct(*regUrl, credentials, installedProduct)
	if err != nil {
		log.Fatalf(err.Error())
	}

	DumpRepositories(os.Stdout, product)
}
