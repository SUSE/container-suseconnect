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
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const sle11CertDir = "/etc/ssl/certs"

var (
	suseRelease   = "/etc/SuSE-release"
	releaseRegexp = regexp.MustCompile("SUSE Linux Enterprise(.*)11")
)

// Check whether the current machine is a SLE11 or not.
func isSLE11() bool {
	f, err := os.Open(suseRelease)
	if err != nil {
		log.Printf("Could not open %v: %v", suseRelease, err)
		return false
	}

	// Just parsing the first line should be enough.
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := strings.TrimSpace(scanner.Text())
	_ = f.Close()
	return releaseRegexp.MatchString(line)
}

// Load a TLS configuration for SLE11.
func sle11TLSConfig() (*tls.Config, error) {
	fi, err := os.Stat(sle11CertDir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s not a dir", sle11CertDir)
	}
	pool := x509.NewCertPool()
	cfg := &tls.Config{RootCAs: pool}

	f, err := os.Open(sle11CertDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, _ := f.Readdirnames(-1)
	for _, name := range names {
		path := filepath.Join(sle11CertDir, name)
		fi, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink || fi.IsDir() {
			continue
		}

		pem, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		log.Printf("Adding PEM: %s", path)
		pool.AppendCertsFromPEM(pem)
	}
	return cfg, nil
}

// Returns the proper TLS config to be used.
func tlsConfig(insecure bool) *tls.Config {
	if !isSLE11() {
		return &tls.Config{InsecureSkipVerify: insecure}
	}
	if cfg, err := sle11TLSConfig(); err != nil {
		log.Printf("Error while adding SLE11 certificates: %v", err)
	} else {
		cfg.InsecureSkipVerify = insecure
		return cfg
	}
	return &tls.Config{InsecureSkipVerify: insecure}
}
