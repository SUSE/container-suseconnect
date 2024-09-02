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
)

func productHelper(t *testing.T, product Product, expectedVersion string) {
	if product.ProductType != "base" {
		t.Fatal("Wrong base for product")
	}
	if product.Identifier != "SLES" {
		t.Fatal("Wrong identifier for product")
	}
	if product.Version != expectedVersion {
		t.Fatal("Wrong version for product")
	}
	if product.Arch != "x86_64" {
		t.Fatal("Wrong arch for product")
	}
}

func productHelperSLE12(t *testing.T, product Product) {
	productHelper(t, product, "12")

	if len(product.Repositories) != 4 {
		t.Fatalf("Wrong number of repos %v", len(product.Repositories))
	}

	if product.Repositories[3].Name != "SLES12-Debuginfo-Pool" {
		t.Fatal("Unexpected value")
	}

	expectedURL := "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product_debug"
	if string(product.Repositories[3].URL) != expectedURL {
		t.Fatalf("Unexpected repository URL: %s", product.Repositories[3].URL)
	}
}

func productHelperSLE15RMT(t *testing.T, product Product) {
	productHelper(t, product, "15.1")

	if len(product.Repositories) != 6 {
		t.Fatal("Wrong number of repos")
	}

	if product.Extensions[0].Repositories[2].Name != "SLE-Module-Basesystem15-SP1-Pool" {
		t.Fatalf("Unexpected Extension Name: %v", product.Extensions[0].Repositories[2].Name)
	}

	expectedURL := "https://smt-ec2.susecloud.net/repo/SUSE/Products/SLE-Module-Basesystem/15-SP1/x86_64/product/?credentials=SCCcredentials"
	if string(product.Extensions[0].Repositories[2].URL) != expectedURL {
		t.Fatalf("Unexpected repository URL: %s", product.Extensions[0].Repositories[2].URL)
	}
}

// Tests for the parseProduct function.

func TestUnreadableProduct(t *testing.T) {
	file, err := os.Open("non-existant-file")
	if err == nil {
		file.Close()
		t.Fatal("This should've been an error...")
	}

	_, err = parseProducts(file)
	if err == nil || err.Error() != "Can't read product information: invalid argument" {
		t.Fatal("This is not the proper error we're expecting")
	}
}

func TestInvalidJsonForProduct(t *testing.T) {
	reader := strings.NewReader("invalid json is invalid")
	_, err := parseProducts(reader)

	if err == nil ||
		err.Error() != "Can't read product information: invalid character 'i' looking for beginning of value - invalid json is invalid" {

		t.Fatalf("This is not the proper error we're expecting: %v", err)
	}
}

func TestValidProduct(t *testing.T) {
	file, err := os.Open("testdata/products-sle12.json")
	if err != nil {
		t.Fatal("Something went wrong when reading the JSON file")
	}
	defer file.Close()

	products, err := parseProducts(file)
	if err != nil {
		t.Fatal("Unexpected error when reading a valid JSON file")
	}

	if len(products) != 1 {
		t.Fatalf("Unexpected number of products found. Got %d, expected %d", len(products), 1)
	}

	productHelperSLE12(t, products[0])
}

// Tests for the requestProduct function.

func TestInvalidRequestForProduct(t *testing.T) {
	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: ":", Insecure: true}

	_, err := RequestProducts(data, cr, ip)
	if err == nil || !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Fatalf("There should be a proper error: %v", err)
	}
}

func TestFaultyRequestForProduct(t *testing.T) {
	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: "http://", Insecure: true}

	_, err := RequestProducts(data, cr, ip)
	if err == nil || !strings.HasSuffix(err.Error(), "no Host in request URL") {
		t.Fatalf("There should be a proper error: %v", err)
	}
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
	if err == nil || err.Error() != "Unexpected error while retrieving products with regCode : 500 Internal Server Error" {
		t.Fatalf("It should have a proper error: %v", err)
	}
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
	if err != nil {
		t.Fatal("It should've run just fine...")
	}

	if len(products) != 1 {
		t.Fatalf("Unexpected number of products found. Got %d, expected %d", len(products), 1)
	}

	productHelperSLE12(t, products[0])
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
	if err != nil {
		t.Fatal("It should've run just fine...")
	}

	if len(products) != 1 {
		t.Fatalf("Unexpected number of products found. Got %d, expected %d", len(products), 1)
	}

	productHelperSLE15RMT(t, products[0])
}
