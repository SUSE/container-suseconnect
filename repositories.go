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

type Repository struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Url string `json:"url"`
	Autorefresh string `json:"autorefresh"`
	Enabled string `json:"enabled"`
}

type Product struct {
	ProductType string `json:"product_type"`
	Identifier string `json:"identifier"`
	Version string `json:"version"`
	Arch string `json:"arch"`
}
