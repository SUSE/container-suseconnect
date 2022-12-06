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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Subscription has all the information that we need for SLE subscriptions.
type Subscription struct {
	RegCode string `json:"regcode"`
	Status  string `json:"status"`
}

// Request registration codes to the registration server. The `data` and the
// `credentials` parameters are used in order to establish the connection with
// the registration server. The `installed` parameter contains the product to
// be requested.
// This function uses SCC's "/connect/systems/subscriptions" API
func requestRegcodes(data SUSEConnectData, credentials Credentials) ([]string, error) {
	var codes []string
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: data.Insecure},
		Proxy:           http.ProxyFromEnvironment,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", data.SccURL, nil)
	if err != nil {
		return codes,
			loggedError(NetworkError, "Could not connect with registration server: %v\n", err)
	}

	req.URL.Path = "/connect/systems/subscriptions"

	auth := url.UserPassword(credentials.Username, credentials.Password)
	req.URL.User = auth

	resp, err := client.Do(req)
	if err != nil {
		return codes, err
	}

	if resp.StatusCode == 404 {
		// we cannot requesting regcodes from s SMT server. It does not
		// has this API. Just return a empty string
		log.Println("Cannot fetch regcodes. Assuming it is SMT server")
		codes = append(codes, "")
		return codes, nil
	}

	if resp.StatusCode != 200 {
		return codes,
			loggedError(SubscriptionServerError, "Unexpected error while retrieving regcode: %s", resp.Status)
	}

	subscriptions, err := parseSubscriptions(resp.Body)
	if err != nil {
		return codes, err
	}

	for _, subscription := range subscriptions {
		if strings.ToUpper(subscription.Status) != "EXPIRED" {
			codes = append(codes, subscription.RegCode)
		} else {
			loggedError(SubscriptionServerError, "Skipping regCode: %s -- expired.", subscription.RegCode)
		}
	}
	return codes, err
}

// Parse the product as expected from the given reader. This function already
// checks whether the given reader is valid or not.
func parseSubscriptions(reader io.Reader) ([]Subscription, error) {
	var subscriptions []Subscription

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return subscriptions, loggedError(SubscriptionError, "Can't read subscriptions information: %v", err.Error())
	}

	err = json.Unmarshal(data, &subscriptions)
	if err != nil {
		return subscriptions, loggedError(SubscriptionError, "Can't read subscription: %v", err.Error())
	}
	if len(subscriptions) == 0 {
		return subscriptions, loggedError(SubscriptionError, "Got 0 subscriptions")
	}
	return subscriptions, nil
}
