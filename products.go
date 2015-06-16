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

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// All the information we need from repositories as given by the registration
// server.
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Autorefresh bool   `json:"autorefresh"`
	Enabled     bool   `json:"enabled"`
}

// All the information we need from product as given by the registration
// server. It contains a slice of repositories in it.
type Product struct {
	ProductType  string       `json:"product_type"`
	Identifier   string       `json:"identifier"`
	Version      string       `json:"version"`
	Arch         string       `json:"arch"`
	Repositories []Repository `json:"repositories"`
}

// Parse the product as expected from the given reader. This function already
// checks whether the given reader is valid or not.
func parseProducts(reader io.Reader) ([]Product, error) {
	var products []Product

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return products,
			fmt.Errorf("Can't read product information: %v", err.Error())
	}

	err = json.Unmarshal(data, &products)
	if err != nil {
		return products,
			fmt.Errorf("Can't read product information: %v - %s", err.Error(), data)
	}
	return products, nil
}

// Request product information to the registration server. The `regCode`
// parameters is used to establish the connection with
// the registration server. The `installed` parameter contains the product to
// be requested.
// This function relies on [/connect/subscriptions/products](https://github.com/SUSE/connect/wiki/SCC-API-%28Implemented%29#product) API.
func requestProductsFromRegCode(data SUSEConnectData, regCode string,
	installed InstalledProduct) ([]Product, error) {
	var products []Product
	var err error

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: data.Insecure},
		Proxy:           http.ProxyFromEnvironment,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", data.SccURL, nil)
	if err != nil {
		return products,
			fmt.Errorf("Could not connect with registration server: %v\n", err)
	}

	values := req.URL.Query()
	values.Add("identifier", installed.Identifier)
	values.Add("version", installed.Version)
	values.Add("arch", installed.Arch)
	req.URL.RawQuery = values.Encode()
	req.URL.Path = "/connect/subscriptions/products"
	if data.SccURL == sccURLStr {
		req.Header.Add("Authorization", `Token token=`+regCode)
	}

	resp, err := client.Do(req)
	if err != nil {
		return products, err
	}
	if resp.StatusCode != 200 {
		return products,
			fmt.Errorf("Unexpected error while retrieving products with regCode %s: %s", regCode, resp.Status)
	}

	return parseProducts(resp.Body)
}

// Request product information to the registration server. The `data` and the
// `credentials` parameters are used in order to establish the connection with
// the registration server. The `installed` parameter contains the product to
// be requested.
func requestProducts(data SUSEConnectData, credentials Credentials,
	installed InstalledProduct) ([]Product, error) {
	var products []Product
	var regCodes []string
	var err error

	if data.SccURL == sccURLStr {
		regCodes, err = requestRegcodes(data, credentials)
		if err != nil {
			return products, err
		}
	} else {
		// SMT does not have this API and does not need a regcode
		regCodes = append(regCodes, "")
	}

	for _, regCode := range regCodes {
		p, err := requestProductsFromRegCode(data, regCode, installed)
		if err != nil {
			var emptyProducts []Product
			return emptyProducts, err
		}
		products = append(products, p...)
	}

	return products, nil
}
