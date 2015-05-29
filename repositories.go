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
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

type Repository struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Url string `json:"url"`
	Autorefresh bool `json:"autorefresh"`
	Enabled bool `json:"enabled"`
}

type Product struct {
	ProductType string `json:"product_type"`
	Identifier string `json:"identifier"`
	Version string `json:"version"`
	Arch string `json:"arch"`
	Repositories []Repository `json:"repositories'`
}

func ParseProduct(reader io.Reader) (Product, error) {
	product := Product{}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return product, err
	}

	err = json.Unmarshal(data, &product)
	if err != nil {
		return product, fmt.Errorf("Can't read product information: %v", err.Error())
	}
	return product, nil
}

// request product information to the registration server
// url is the registration server url
// installedProduct is the product you are requesting
func RequestProduct(url string, credentials Credentials, installed InstalledProduct) (Product, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	values := req.URL.Query()

	values.Add("identifier", installed.Identifier)
	values.Add("version", installed.Version)
	values.Add("arch", installed.Arch)
	req.URL.RawQuery = values.Encode()
	req.URL.Path = "/connect/systems/products"

	resp, err := client.Do(req)
	if err != nil {
		return Product{}, err
	}

	return ParseProduct(resp.Body)
}




