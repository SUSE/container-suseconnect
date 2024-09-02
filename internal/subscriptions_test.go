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

func subscriptionHelper(t *testing.T, subscription Subscription) {
	if subscription.RegCode != "35098ff7" {
		t.Fatal("Wrong regcode for subscription")
	}
}

// Tests for the parseSubscriptions function.

func TestUnreadableSubscription(t *testing.T) {
	file, err := os.Open("non-existant-file")
	if err == nil {
		file.Close()
		t.Fatal("This should've been an error...")
	}

	_, err = parseSubscriptions(file)
	if err == nil || err.Error() != "Can't read subscriptions information: invalid argument" {
		t.Fatalf("This is not the proper error we're expecting: %v", err)
	}
}

func TestInvalidJsonForSubscriptions(t *testing.T) {
	reader := strings.NewReader("invalid json is invalid")
	_, err := parseSubscriptions(reader)

	if err == nil ||
		err.Error() != "Can't read subscription: invalid character 'i' looking for beginning of value" {

		t.Fatalf("This is not the proper error we're expecting: %v", err)
	}
}

func TestEmptySubscriptions(t *testing.T) {
	file, err := os.Open("testdata/empty-subscriptions.json")
	if err != nil {
		t.Fatal("Something went wrong when reading the JSON file")
	}
	defer file.Close()

	subscriptions, err := parseSubscriptions(file)
	if err == nil || err.Error() != "Got 0 subscriptions" {
		t.Fatal("Unexpected error when reading a valid JSON file")
	}
	if len(subscriptions) != 0 {
		t.Fatalf("It should be empty")
	}
}

func TestValidSubscriptions(t *testing.T) {
	file, err := os.Open("testdata/subscriptions.json")
	if err != nil {
		t.Fatal("Something went wrong when reading the JSON file")
	}
	defer file.Close()

	subscriptions, err := parseSubscriptions(file)
	if err != nil {
		t.Fatal("Unexpected error when reading a valid JSON file")
	}

	if len(subscriptions) != 2 {
		t.Fatalf("Unexpected number of subscriptions found. Got %d, expected %d", len(subscriptions), 1)
	}

	subscriptionHelper(t, subscriptions[0])
}

// Tests for the requestRegcodes function.

func TestInvalidRequestForRegcodes(t *testing.T) {
	var cr Credentials
	data := SUSEConnectData{SccURL: ":", Insecure: true}

	_, err := requestRegcodes(data, cr)
	if err == nil || !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Fatalf("There should be a proper error: %v", err)
	}
}

func TestFaultyRequestForRegcodes(t *testing.T) {
	var cr Credentials
	data := SUSEConnectData{SccURL: "http://", Insecure: true}

	_, err := requestRegcodes(data, cr)
	if err == nil || !strings.Contains(err.Error(), "no Host in request URL") {
		t.Fatalf("There should be a proper error: %v", err)
	}
}

func TestRemoteErrorWhileRequestingRegcodes(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something bad happened", 500)
	}))
	defer ts.Close()

	var cr Credentials
	data := SUSEConnectData{SccURL: ts.URL, Insecure: true}

	_, err := requestRegcodes(data, cr)
	if err == nil || err.Error() != "Unexpected error while retrieving regcode: 500 Internal Server Error" {
		t.Fatalf("There should be a proper error:  %v", err)
	}
}

func TestValidRequestForRegcodes(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/subscriptions.json")
		if err != nil {
			fmt.Fprintln(w, "FAIL!")
			return
		}
		io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	var cr Credentials
	data := SUSEConnectData{SccURL: ts.URL, Insecure: true}

	codes, err := requestRegcodes(data, cr)
	if err != nil {
		t.Fatal("It should've run just fine...")
	}

	// This also tests that we're not including expired regcodes
	if len(codes) != 1 {
		t.Fatalf("Unexpected number of products found. Got %d, expected %d", len(codes), 1)
	}

	if codes[0] != "35098ff7" {
		t.Fatalf("Got the wrong registration code: %v", codes[0])
	}
}

func TestRequestEmptyRegcodes(t *testing.T) {
	// We setup a fake http server that mocks a registration server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/empty-subscriptions.json")
		if err != nil {
			fmt.Fprintln(w, "FAIL!")
			return
		}
		io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	var cr Credentials
	data := SUSEConnectData{SccURL: ts.URL, Insecure: true}

	codes, err := requestRegcodes(data, cr)
	if err == nil || err.Error() != "Got 0 subscriptions" {
		t.Fatal("Unexpected error when reading a valid JSON file")
	}

	if len(codes) != 0 {
		t.Fatalf("It should be 0")
	}
}
