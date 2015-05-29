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
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	baseProductLoc string = "/etc/products.d/baseproduct"
)

type InstalledProduct struct {
	Identifier string `xml:"name"`
	Version string `xml:"version`
	Arch string `xml:"arch"`
}

// parses installed product data
func ParseInstalledProduct(reader io.Reader) (InstalledProduct, error) {
	xmlData, err := ioutil.ReadAll(reader)
	if err != nil {
		return InstalledProduct{}, fmt.Errorf("Can't read base product file: %v", err.Error())
	}

	var p InstalledProduct
	xml.Unmarshal(xmlData, &p)
	if err != nil {
		return InstalledProduct{}, fmt.Errorf("Can't parse base product file: %v", err.Error())
	}

	return p, nil
}

// read the product file from the standard location
func ReadInstalledProduct() (InstalledProduct, error) {
	if _, err := os.Stat(baseProductLoc); os.IsNotExist(err) {
		return InstalledProduct{}, fmt.Errorf("No base product detected")
	}

	xmlFile, err := os.Open(baseProductLoc)
	if err != nil {
		return InstalledProduct{}, fmt.Errorf("Can't open base product file: %v", err.Error())
	}
	defer xmlFile.Close()

	return ParseInstalledProduct(xmlFile)
}
