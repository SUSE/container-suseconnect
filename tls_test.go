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

import "testing"

func TestIsSLE11(t *testing.T) {
	sr := suseRelease
	defer func() {
		suseRelease = sr
	}()

	suseRelease = "data/sles12-release"
	if isSLE11() {
		t.Fatalf("It should not have been a SLE11: %s", sr)
	}

	suseRelease = "data/sle11sp3-release"
	if !isSLE11() {
		t.Fatalf("It should have been a SLE11: %s", sr)
	}
	suseRelease = "data/unknown-release"
	if isSLE11() {
		t.Fatalf("It should not have been a SLE11: %s", sr)
	}
}
