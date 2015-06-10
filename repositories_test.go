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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func productHelper(t *testing.T, product Product) {
	if product.ProductType != "base" {
		t.Fatal("Wrong base for product")
	}
	if product.Identifier != "SLES" {
		t.Fatal("Wrong identifier for product")
	}
	if product.Version != "12" {
		t.Fatal("Wrong version for product")
	}
	if product.Arch != "x86_64" {
		t.Fatal("Wrong arch for product")
	}
	if len(product.Repositories) != 4 {
		t.Fatal("Wrong number of repos")
	}
	if product.Repositories[3].Name != "SLES12-Debuginfo-Pool" {
		t.Fatal("Unexpected value")
	}
	if product.Repositories[3].URL != "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product_debug" {
		t.Fatal("Unexpected value")
	}
}

// Tests for the parseProduct function.

func TestUnreadableProduct(t *testing.T) {
	file, err := os.Open("non-existant-file")
	if err == nil {
		file.Close()
		t.Fatal("This should've been an error...")
	}

	_, err = parseProduct(file)
	if err == nil || err.Error() != "Can't read product information: invalid argument" {
		t.Fatal("This is not the proper error we're expecting")
	}
}

func TestInvalidJson(t *testing.T) {
	reader := strings.NewReader("invalid json is invalid")
	_, err := parseProduct(reader)

	if err == nil ||
		err.Error() != "Can't read product information: invalid character 'i' looking for beginning of value" {

		t.Fatal("This is not the proper error we're expecting")
	}
}

func TestValidProduct(t *testing.T) {
	file, err := os.Open("data/product.json")
	if err != nil {
		t.Fatal("Something went wrong when reading the JSON file")
	}
	defer file.Close()

	product, err := parseProduct(file)
	if err != nil {
		t.Fatal("Unexpected error when reading a valid JSON file")
	}
	productHelper(t, product)
}

// Tests for the requestProduct function.

func TestInvalidRequest(t *testing.T) {
	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: ":", Insecure: true}

	_, err := requestProduct(data, cr, ip)
	if err == nil || err.Error() != "Could not connect with registration server: parse :: missing protocol scheme\n" {
		t.Fatal("There should be a proper error")
	}
}

func TestFaultyRequest(t *testing.T) {
	var cr Credentials
	var ip InstalledProduct
	data := SUSEConnectData{SccURL: "http://", Insecure: true}

	_, err := requestProduct(data, cr, ip)
	str := "Get http://:@/connect/systems/products?arch=&identifier=&version=: http: no Host in request URL"
	if err == nil || err.Error() != str {
		t.Fatal("There should be a proper error")
	}
}

func TestValidRequest(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("data/product.json")
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

	product, err := requestProduct(data, cr, ip)
	if err != nil {
		t.Fatal("It should've run just fine...")
	}
	productHelper(t, product)
}
