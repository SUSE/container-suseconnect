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

func TestUnreadableSubscription(t *testing.T) {
	invalidFile := (*os.File)(nil)
	_, err := parseSubscriptions(invalidFile)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can't read subscriptions information: invalid argument")
}

func TestInvalidJsonForSubscriptions(t *testing.T) {
	reader := strings.NewReader("invalid json is invalid")
	_, err := parseSubscriptions(reader)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can't read subscription: invalid character 'i' looking for beginning of value")
}

func TestEmptySubscriptions(t *testing.T) {
	file, err := os.Open("testdata/empty-subscriptions.json")
	require.NotNil(t, file)
	defer file.Close()

	subscriptions, err := parseSubscriptions(file)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Got 0 subscriptions")
	assert.Empty(t, subscriptions)
}

func TestValidSubscriptions(t *testing.T) {
	file, err := os.Open("testdata/subscriptions.json")
	require.NotNil(t, file)
	defer file.Close()

	subscriptions, err := parseSubscriptions(file)
	require.Nil(t, err)

	if assert.Len(t, subscriptions, 2) {
		assert.Equal(t, "35098ff7", subscriptions[0].RegCode)
	}
}

func TestInvalidRequestForRegcodes(t *testing.T) {
	var cr Credentials
	data := SUSEConnectData{SccURL: ":", Insecure: true}

	_, err := requestRegcodes(data, cr)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing protocol scheme")
}

func TestFaultyRequestForRegcodes(t *testing.T) {
	var cr Credentials
	data := SUSEConnectData{SccURL: "http://", Insecure: true}

	_, err := requestRegcodes(data, cr)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "no Host in request URL")
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
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Unexpected error while retrieving regcode: 500 Internal Server Error")
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
	require.Nil(t, err)

	if assert.Len(t, codes, 1) {
		assert.Equal(t, "35098ff7", codes[0])
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
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Got 0 subscriptions")
	assert.Empty(t, codes)
}
