// Copyright (c) 2023 SUSE LLC. All rights reserved.
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
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Convert from bool to int.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// DumpRepositories dumps the repositories of the given product to the given
// writer.
func DumpRepositories(w io.Writer, product Product) {
	fmt.Fprintf(w, "# generated by container-suseconnect\n")
	fmt.Fprintf(w, "\n")

	// Always print the first product disregarding if it is recommended or not
	dumpRepositoriesRecursive(w, product, true)
}

// dumpRepositoriesRecursive prints available product repositories to the
// provided writer.
//
// The function takes all product extensions into account, which will be
// printed recursively too.
//
// `dumpAlways` specifies if the products repositories should be
// printed ignoring if it is recommended or not.
//
// The user has the option to enable certain modules via its `identifier` by
// setting them within the `ADDITIONAL_MODULES` environment variable. Multiple
// modules can be set comma separated.
func dumpRepositoriesRecursive(
	w io.Writer,
	product Product,
	dumpAlways bool,
) {
	for _, repo := range product.Repositories {
		if product.Recommended || dumpAlways ||
			moduleEnabledInEnv(product.Identifier) ||
			moduleEnabledInProductFiles(product.Identifier) {

			fmt.Fprintf(w, "[%s]\n", repo.Name)
			fmt.Fprintf(w, "name=%s\n", repo.Description)
			fmt.Fprintf(w, "baseurl=%s\n", repo.URL)
			fmt.Fprintf(w, "autorefresh=%d\n", boolToInt(repo.Autorefresh))
			fmt.Fprintf(w, "enabled=%d\n", boolToInt(repo.Enabled))
			fmt.Fprintf(w, "\n")
		}
	}

	// Continue the traversal for the extensions if needed
	if len(product.Extensions) > 0 {
		for _, extension := range product.Extensions {
			dumpRepositoriesRecursive(w, extension, false)
		}
	}
}

// moduleEnabledInEnv returns true if the provided `identifier` is included in
// the `ADDITIONAL_MODULES` environment variable, otherwise false.
func moduleEnabledInEnv(identifier string) bool {
	for _, i := range strings.Split(os.Getenv("ADDITIONAL_MODULES"), ",") {
		if identifier == i {
			return true
		}
	}

	return false
}

// moduleEnabledInProductFiles returns true if the provided `identifier` is
// a name of a file in the  /etc/product.d/*.prod, otherwise false.
func moduleEnabledInProductFiles(identifier string) bool {
	files, err := os.ReadDir("/etc/products.d")
	if err != nil {
		return false
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		info, err := file.Info()
		if err != nil {
			continue
		}

		if ext == ".prod" &&
			identifier == strings.TrimSuffix(file.Name(), ext) &&
			info.Mode()&os.ModeSymlink == 0 {
			return true
		}
	}

	return false
}

// ListModules prints the provided `products` slice the provided writer `w` in a
// human readable way
func ListModules(w io.Writer, products []Product) {
	for _, product := range products {
		if product.ProductType == "module" {
			fmt.Fprintf(w, "Name: %v\n", product.Name)
			fmt.Fprintf(w, "Identifier: %v\n", product.Identifier)
			fmt.Fprintf(w, "Recommended: %v\n", product.Recommended)
			fmt.Fprintf(w, "\n")
		}

		// Continue traversal with available product extensions
		ListModules(w, product.Extensions)
	}
}

// ListProducts prints the provided `products` slice the provided writer `w` in a
// human readable way
func ListProducts(w io.Writer, products []Product, baseProduct string) {
	for _, product := range products {
		fmt.Fprintf(w, "Name: %v\n", product.Name)
		fmt.Fprintf(w, "Type: %v\n", product.ProductType)
		fmt.Fprintf(w, "Identifier: %v\n", product.Identifier)
		fmt.Fprintf(w, "Based on: %v\n", baseProduct)
		fmt.Fprintf(w, "Recommended: %v\n", product.Recommended)

		// Strip HTML tags from the description
		r := regexp.MustCompile(`<[^>]*>\s*`)
		fmt.Fprintf(w, "Description: %v\n",
			strings.TrimSpace(r.ReplaceAllString(product.Description, "")))

		fmt.Fprintf(w, "Repositories:\n")
		for idx, repo := range product.Repositories {
			repoState := "disabled"
			if repo.Enabled {
				repoState = "enabled"
			}
			fmt.Fprintf(w, "%v. %v: %v (%v)\n", idx+1, repo.Name, repo.URL, repoState)
		}
		fmt.Fprintf(w, "\n")

		// Continue traversal with available product extensions
		ListProducts(w, product.Extensions, product.Identifier)
	}
}
