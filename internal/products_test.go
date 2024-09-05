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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnreadableProduct(t *testing.T) {
	invalidFile := (*os.File)(nil)
	_, err := parseProducts(invalidFile)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can't read product information: invalid argument")
}

func TestInvalidJsonForProduct(t *testing.T) {
	reader := strings.NewReader("invalid json is invalid")
	_, err := parseProducts(reader)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can't read product information: invalid character 'i' looking for beginning of value - invalid json is invalid")
}

func TestValidProduct(t *testing.T) {
	file, err := os.Open("testdata/products-sle12.json")
	require.Nil(t, err)
	defer file.Close()

	products, err := parseProducts(file)
	require.Nil(t, err)
	require.Len(t, products, 1)

	product := products[0]
	assert.Equal(t, "12", product.Version)
	assert.Equal(t, "base", product.ProductType)
	assert.Equal(t, "SLES", product.Identifier)
	assert.Equal(t, "x86_64", product.Arch)

	if assert.Len(t, product.Repositories, 4) {
		assert.Equal(t, "SLES12-Debuginfo-Pool", product.Repositories[3].Name)
		expectedURL := "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product_debug"
		assert.Equal(t, expectedURL, product.Repositories[3].URL)
	}
}

func TestInvalidRequestForProduct(t *testing.T) {
	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: ":", Insecure: true}

	_, err := RequestProducts(data, cr, ip)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing protocol scheme")
}

func TestFaultyRequestForProduct(t *testing.T) {
	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: "http://", Insecure: true}

	_, err := RequestProducts(data, cr, ip)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "no Host in request URL")
}

func TestRemoteErrorWhileRequestingProducts(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	firstRequest := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First request should return 404 to make the function request
		// products and return 500 in the second request
		if firstRequest {
			http.Error(w, "something bad happened", 404)
		} else {
			http.Error(w, "something bad happened", 500)
		}
		firstRequest = false
	}))
	defer ts.Close()

	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: ts.URL, Insecure: true}

	_, err := RequestProducts(data, cr, ip)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Unexpected error while retrieving products with regCode : 500 Internal Server Error")
}

func TestValidRequestForProduct(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/products-sle12.json")
		if err != nil {
			fmt.Fprintln(w, "FAIL!")
			return
		}
		io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: ts.URL, Insecure: true}

	products, err := RequestProducts(data, cr, ip)
	require.Nil(t, err)
	require.Len(t, products, 1)

	product := products[0]
	assert.Equal(t, "12", product.Version)
	assert.Equal(t, "base", product.ProductType)
	assert.Equal(t, "SLES", product.Identifier)
	assert.Equal(t, "x86_64", product.Arch)

	if assert.Len(t, product.Repositories, 4) {
		assert.Equal(t, "SLES12-Debuginfo-Pool", product.Repositories[3].Name)
		expectedURL := "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product_debug"
		assert.Equal(t, expectedURL, product.Repositories[3].URL)
	}
}

func TestValidRequestForProductUsingRMT(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SMT servers return 404 on this URL
		if r.URL.Path == "/connect/systems/subscriptions" {
			http.Error(w, "", http.StatusNotFound)
		}

		// The result also looks slightly different
		resFile := "testdata/products-sle15-rmt.json"
		file, err := os.Open(resFile)
		if err != nil {
			fmt.Fprintln(w, "FAIL!")
			return
		}
		io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: ts.URL, Insecure: true}

	products, err := RequestProducts(data, cr, ip)
	require.Nil(t, err)
	require.Len(t, products, 1)

	product := products[0]
	assert.Equal(t, "15.1", product.Version)
	assert.Equal(t, "base", product.ProductType)
	assert.Equal(t, "SLES", product.Identifier)
	assert.Equal(t, "x86_64", product.Arch)

	if assert.Len(t, product.Repositories, 6) && assert.Len(t, product.Extensions, 1) {
		assert.Equal(t, "SLE-Module-Basesystem15-SP1-Pool", product.Extensions[0].Repositories[2].Name)
		expectedURL := "https://smt-ec2.susecloud.net/repo/SUSE/Products/SLE-Module-Basesystem/15-SP1/x86_64/product/?credentials=SCCcredentials"
		assert.Equal(t, expectedURL, product.Extensions[0].Repositories[2].URL)
	}
}
