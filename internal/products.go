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
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Repository has all the information we need from repositories as given by the
// registration server.
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Autorefresh bool   `json:"autorefresh"`
	Enabled     bool   `json:"enabled"`
}

// Product has all the information we need from product as given by the registration
// server. It contains a slice of repositories in it.
type Product struct {
	ProductType  string       `json:"product_type"`
	Identifier   string       `json:"identifier"`
	Version      string       `json:"version"`
	Arch         string       `json:"arch"`
	Repositories []Repository `json:"repositories"`
	Extensions   []Product    `json:"extensions"`
	Recommended  bool         `json:"recommended"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
}

// Take the "Product" as returned from an RMT and adjust the repository URLs
// to include a "credentials" parameter. This is somewhat specific to the RMT
// instances available for Public Cloud on-demand instances, as they need to
// authenticate to be able to access repositories.
func fixRepoUrlsForRMT(p *Product) error {
	for i := range p.Repositories {
		repourl, err := url.Parse(p.Repositories[i].URL)
		if err != nil {
			loggedError(RepositoryError, "Unable to parse repository URL: %s - %v", p.Repositories[i].URL, err)
			return err
		}

		params := repourl.Query()
		if params.Get("credentials") == "" {
			params["credentials"] = []string{"SCCcredentials"}
		}

		repourl.RawQuery = params.Encode()
		p.Repositories[i].URL = repourl.String()
	}

	// Products can have nested Products (Extensions) handle those recursively
	for i := range p.Extensions {
		if err := fixRepoUrlsForRMT(&p.Extensions[i]); err != nil {
			return err
		}
	}

	return nil
}

// Parse the product as expected from the given reader. This function already
// checks whether the given reader is valid or not.
func parseProducts(reader io.Reader) ([]Product, error) {
	var products []Product

	data, err := io.ReadAll(reader)
	if err != nil {
		return products,
			loggedError(RepositoryError, "Can't read product information: %v", err.Error())
	}

	// Depending on which API was used the JSON we get passed contains
	// either a list of products ("/connect/subscriptions/products") with
	// a single element in it or a single product
	// ("/connect/systems/products").  So we need to Unmarshall() slightly
	// different for both cases.
	err = json.Unmarshal(data, &products)
	if err != nil {
		products = nil
		// When connected to RMT we only got a single "Product", so let's try
		// to unmarshall that. (And add a "credential" parameter to it if that
		// is not present)
		var product Product
		err = json.Unmarshal(data, &product)
		if err == nil {
			fixRepoUrlsForRMT(&product)
			products = append(products, product)
		}
	}

	if err != nil {
		return products,
			loggedError(RepositoryError, "Can't read product information: %v - %s", err.Error(), data)
	}

	return products, nil
}

// Request product information to the registration server. The `regCode`
// parameters is used to establish the connection with
// the registration server. The `installed` parameter contains the product to
// be requested.
// This function relies on [/connect/subscriptions/products](https://github.com/SUSE/connect/wiki/SCC-API-%28Implemented%29#product) API.
func requestProductsFromRegCodeOrSystem(data SUSEConnectData, regCode string,
	credentials Credentials, installed InstalledProduct,
) ([]Product, error) {
	var products []Product
	var err error

	req, err := http.NewRequest("GET", data.SccURL, nil)
	if err != nil {
		return products, loggedError(NetworkError, "Could not connect with registration server: %v\n", err)
	}

	values := req.URL.Query()
	values.Add("identifier", installed.Identifier)
	values.Add("version", installed.Version)
	values.Add("arch", installed.Arch)
	req.URL.RawQuery = values.Encode()

	if len(regCode) > 0 {
		req.Header.Add("Authorization", `Token token=`+regCode)
		req.URL.Path = "/connect/subscriptions/products"
	} else {
		// we're connected to a RMT, which does not provide the /connect/subscriptions/products
		// endpoint. Fallback to "/connect/systems/products" here.
		req.URL.Path = "/connect/systems/products"
		auth := url.UserPassword(credentials.Username, credentials.Password)
		req.URL.User = auth

		if credentials.SystemToken != "" {
			req.Header.Add("System-Token", credentials.SystemToken)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: data.Insecure,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return products, err
	}

	if resp.StatusCode != 200 {
		var payload map[string]interface{}
		dec := json.NewDecoder(resp.Body)

		if err := dec.Decode(&payload); err == nil {
			if err, ok := payload["error"]; ok {
				log.Println(err)
			}
		}

		return products, loggedError(SubscriptionServerError, "Unexpected error while retrieving products with regCode %s: %s", regCode, resp.Status)
	}

	return parseProducts(resp.Body)
}

// RequestProducts fetches product information to the registration server. The
// `data` and the `credentials` parameters are used in order to establish the
// connection with the registration server. The `installed` parameter contains
// the product to be requested.
func RequestProducts(data SUSEConnectData, credentials Credentials,
	installed InstalledProduct,
) ([]Product, error) {
	var products []Product
	var regCodes []string
	var err error

	regCodes, err = requestRegcodes(data, credentials)
	if err != nil {
		return products, err
	}

	for _, regCode := range regCodes {
		p, _err := requestProductsFromRegCodeOrSystem(data, regCode, credentials, installed)
		if _err != nil {
			err = _err
			continue
		}

		products = append(products, p...)
	}

	if len(products) > 0 {
		err = nil
	}

	return products, err
}
