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
	"testing"
	"strings"
)

var installedProduct = `
<?xml version="1.0" encoding="UTF-8"?>
<product schemeversion="0">
  <vendor>SUSE</vendor>
  <name>SLES</name>
  <version>12</version>
  <baseversion>12</baseversion>
  <patchlevel>0</patchlevel>
  <predecessor>SUSE_SLES</predecessor>
  <release>0</release>
  <endoflife>2024-10-31</endoflife>
  <arch>x86_64</arch>
  <cpeid>cpe:/o:suse:sles:12</cpeid>
  <productline>sles</productline>
  <register>
      <target>sle-12-x86_64</target>
    <updates>
      <repository repoid="obsrepository://build.suse.de/SUSE:Updates:SLE-SERVER:12:x86_64/update" />
      <repository repoid="obsrepository://build.suse.de/SUSE:Updates:SLE-SERVER:12:x86_64/update_debug" />
    </updates>
  </register>
  <upgrades />
  <updaterepokey>A43242DKD</updaterepokey>
  <summary>SUSE Linux Enterprise Server 12</summary>
  <shortsummary>SLES12</shortsummary>
  <description>SUSE Linux Enterprise offers a comprehensive
        suite of products built on a single code base.
        The platform addresses business needs from
        the smallest thin-client devices to the world's
        most powerful high-performance computing
        and mainframe servers. SUSE Linux Enterprise
        offers common management tools and technology
        certifications across the platform, and
        each product is enterprise-class.</description>
  <linguas>
    <language>cs</language>
    <language>da</language>
    <language>de</language>
    <language>en</language>
    <language>en_GB</language>
    <language>en_US</language>
    <language>es</language>
    <language>fi</language>
    <language>fr</language>
    <language>hu</language>
    <language>it</language>
    <language>ja</language>
    <language>nb</language>
    <language>nl</language>
    <language>pl</language>
    <language>pt</language>
    <language>pt_BR</language>
    <language>ru</language>
    <language>sv</language>
    <language>zh</language>
    <language>zh_CN</language>
    <language>zh_TW</language>
  </linguas>
  <urls>
    <url name="releasenotes">https://www.suse.com/releasenotes/x86_64/SUSE-SLES/12/release-notes-sles.rpm</url>
  </urls>
  <buildconfig>
    <producttheme>SLES</producttheme>
  </buildconfig>
  <installconfig>
    <defaultlang>en_US</defaultlang>
    <datadir>suse</datadir>
    <descriptiondir>suse/setup/descr</descriptiondir>
    <releasepackage name="sles-release" flag="EQ" version="12" release="1.377" />
    <distribution>SUSE_SLE</distribution>
  </installconfig>
  <runtimeconfig />
  <productdependency relationship="provides" name="SUSE_SLE" baseversion="12" patchlevel="0" flag="EQ" />
  <productdependency relationship="provides" name="SUSE_SLE-SP0" baseversion="12" patchlevel="0" flag="EQ" />
</product>
`

func TestInstalledProductParsing(t *testing.T) {
	reader := strings.NewReader(installedProduct)

	product, err := ParseInstalledProduct(reader)
	if (err != nil) {
		t.FailNow()
	}

	if product.Identifier != "SLES" {
		t.Errorf(product.Identifier)
	}
}
