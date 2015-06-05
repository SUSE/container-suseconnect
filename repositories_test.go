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
	"strings"
	"testing"
)

var sccReply = `
{
  "arch": "x86_64",
  "product_class": "7261",
  "extensions": [
  ],
  "product_type": "base",
  "identifier": "SLES",
  "free": true,
  "id": 100874,
  "friendly_name": "SUSE Linux Enterprise Server 12 x86_64",
  "available": true,
  "version": "12",
  "name": "SUSE Linux Enterprise Server 12 x86_64",
  "description": "SUSE Linux Enterprise offers a comprehensive suite of products built on a single code base. The platform addresses business needs from the smallest thin-client devices to the world's most powerful high-performance computing and mainframe servers. SUSE Linux Enterprise offers common management tools and technology certifications across the platform, and each product is enterprise-class.",
  "repositories": [
    {
      "distro_target": "sle-12-x86_64",
      "name": "SLES12-Updates",
      "description": "SLES12-Updates for sle-12-x86_64",
      "autorefresh": true,
      "url": "https://smt.test.lan/repo/SUSE/Updates/SLE-SERVER/12/x86_64/update",
      "id": "941",
      "enabled": true
    },
    {
      "distro_target": "sle-12-x86_64",
      "name": "SLES12-Debuginfo-Updates",
      "description": "SLES12-Debuginfo-Updates for sle-12-x86_64",
      "autorefresh": true,
      "url": "https://smt.test.lan/repo/SUSE/Updates/SLE-SERVER/12/x86_64/update_debug",
      "id": "942",
      "enabled": false
    },
    {
      "distro_target": "sle-12-x86_64",
      "name": "SLES12-Pool",
      "description": "SLES12-Pool for sle-12-x86_64",
      "autorefresh": false,
      "url": "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product",
      "id": "943",
      "enabled": true
    },
    {
    "distro_target": "sle-12-x86_64",
    "name": "SLES12-Debuginfo-Pool",
    "description": "SLES12-Debuginfo-Pool for sle-12-x86_64",
    "autorefresh": false,
    "url": "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product_debug",
    "id": "944",
    "enabled": false
    }
  ],
  "release_type": null,
  "former_identifier": "SUSE_SLES",
  "eula_url": "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product.license/",
  "cpe": "cpe:/o:suse:sles:12.0"
}
`

func TestParseSCCReply(t *testing.T) {
	reader := strings.NewReader(sccReply)

	product, err := ParseProduct(reader)
	if err != nil {
		t.Errorf(err.Error())
	}

	n := len(product.Repositories)
	if n != 4 {
		t.Errorf(fmt.Sprintf("Got: %v", n))
	}

	if product.Repositories[3].Name != "SLES12-Debuginfo-Pool" {
		t.Fail()
	}
	if product.Repositories[3].URL != "https://smt.test.lan/repo/SUSE/Products/SLE-SERVER/12/x86_64/product_debug" {
		t.Fail()
	}
}
