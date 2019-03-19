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

package containersuseconnect

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ProductProvider is used to retrieve the location of the file containing the
// information about the installed product.
type ProductProvider interface {
	// Returns the path to the XML file containing the info about the installed
	// product.
	Location() string
}

// Implements the ProductProvider interface so we can fetch the location of the
// SUSE baseproduct file.
type SUSEProductProvider struct{}

func (b SUSEProductProvider) Location() string {
	return "/etc/products.d/baseproduct"
}

// Contains all the info that we need from the installed product.
type InstalledProduct struct {
	Identifier string `xml:"name"`
	Version    string `xml:"version"`
	Arch       string `xml:"arch"`
}

func (p InstalledProduct) String() string {
	return fmt.Sprintf("%s-%s-%s", p.Identifier, p.Version, p.Arch)
}

// Parses installed product data. The passed reader is guaranteed to be
// readable.
func parseInstalledProduct(reader io.Reader) (InstalledProduct, error) {
	// We can ignore this error because of the pre-condition of the `reader`
	// being actually readable.
	xmlData, _ := ioutil.ReadAll(reader)

	var p InstalledProduct
	err := xml.Unmarshal(xmlData, &p)
	if err != nil {
		return InstalledProduct{},
			loggedError("Can't parse base product file: %v", err.Error())
	}
	return p, nil
}

// Read the product file from the standard location
func readInstalledProduct(provider ProductProvider) (InstalledProduct, error) {
	if _, err := os.Stat(provider.Location()); os.IsNotExist(err) {
		return InstalledProduct{}, loggedError("No base product detected")
	}

	xmlFile, err := os.Open(provider.Location())
	if err != nil {
		return InstalledProduct{},
			loggedError("Can't open base product file: %v", err.Error())
	}
	defer xmlFile.Close()

	return parseInstalledProduct(xmlFile)
}

// GetInstalledProduct gets the installed product on a SUSE machine.
func GetInstalledProduct() (InstalledProduct, error) {
	var b SUSEProductProvider
	return readInstalledProduct(b)
}
