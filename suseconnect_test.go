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
	"strings"
	"testing"
)

var suseConnect = `
---
## SUSEConnect configuration file example

## URL of the registration server. (default: https://scc.suse.com)
# url: https://scc.suse.com
 url: 	https://smt.test.lan

## Registration code to use for the base product on the system
# regcode:

## Language code to use for error messages (default: $LANG)
# language:

## Do not verify SSL certificates when using https (default: false)
insecure: false
`

var suseConnectWithoutUrl = `
---
## SUSEConnect configuration file example

## URL of the registration server. (default: https://scc.suse.com)
# url: https://scc.suse.com

## Registration code to use for the base product on the system
# regcode:

## Language code to use for error messages (default: $LANG)
# language:

## Do not verify SSL certificates when using https (default: false)
insecure: false
`

func TestParseSUSEConnect(t *testing.T) {
	reader := strings.NewReader(suseConnect)

	url, err := ParseSUSEConnect(reader)
	if err != nil {
		t.Errorf(err.Error())
	}

	if url != "https://smt.test.lan" {
		t.Fail()
	}
}

func TestParseSUSEConnectWithoutUrl(t *testing.T) {
	reader := strings.NewReader(suseConnectWithoutUrl)

	url, err := ParseSUSEConnect(reader)
	if err != nil {
		t.Errorf(err.Error())
	}

	if url != "" {
		t.Fail()
	}
}
