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

// Gives access to repositories during docker build and run using the host
// machine credentials.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	cs "github.com/SUSE/container-suseconnect/pkg/container-suseconnect"
	"github.com/urfave/cli"
)

func main() {
	// Set the basic CLI metadata
	app := cli.NewApp()
	app.Copyright = fmt.Sprintf("© %d SUSE LCC", time.Now().Year())
	app.Name = "container-suseconnect"
	app.Version = "2.0.0"
	app.Usage = ""
	app.UsageText =
		`This application can be used to retrieve basic metadata about SLE
   related products and module extensions.

   Please use the 'list-product' subcommand for listing all currently
   available products including their repositories and a short description.

   Use the 'list-modules' subcommand for listing available modules, where
   their 'Identifier' can be used to enable them via the ADDITIONAL_MODULES
   environment variable during container creation/run.

   The 'zypper' subcommand runs the application as zypper plugin and is only
   intended to use for debugging purposes.`

	// Switch the application behavior regarding the basename
	switch filepath.Base(os.Args[0]) {
	case "container-suseconnect-zypp":
		app.Action = runZypperPlugin
	default:
		app.Action = runListProducts
	}

	// Set additional actions, which are always available
	app.Commands = []cli.Command{
		{
			Name:    "list-products",
			Aliases: []string{"lp"},
			Usage:   "List available products",
			Action:  runListProducts,
		},
		{
			Name:    "list-modules",
			Aliases: []string{"lm"},
			Usage:   "List available modules",
			Action:  runListModules,
		},
		{
			Name:    "zypper",
			Aliases: []string{"z", "zypp"},
			Usage:   "Run the zypper plugin",
			Action:  runZypperPlugin,
		},
	}

	// Run the application
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func requestProducts() ([]cs.Product, error) {
	credentials := cs.Credentials{}
	if err := cs.ReadConfiguration(&credentials); err != nil {
		return nil, err
	}

	installedProduct, err := cs.GetInstalledProduct()
	if err != nil {
		return nil, err
	}
	log.Printf("Installed product: %v\n", installedProduct)

	var suseConnectData cs.SUSEConnectData
	if err := cs.ReadConfiguration(&suseConnectData); err != nil {
		return nil, err
	}
	log.Printf("Registration server set to %v\n", suseConnectData.SccURL)

	products, err := cs.RequestProducts(suseConnectData, credentials, installedProduct)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// runZypperPlugin runs the application in zypper plugin mode, which dumps
// all available repositories for the installed product. Additional modules
// can be specified via the `ADDITIONAL_MODULES` environment variable, which
// reflect the module `identifier`.
func runZypperPlugin(_ *cli.Context) error {
	log.SetOutput(cs.GetLoggerFile())

	products, err := requestProducts()
	if err != nil {
		return err
	}

	for _, product := range products {
		cs.DumpRepositories(os.Stdout, product)
	}

	return nil
}

// runListModules lists all available modules and their metadata, which
// includes the `Name`, `Identifier` and the `Recommended` flag.
func runListModules(_ *cli.Context) error {
	products, err := requestProducts()
	if err != nil {
		return err
	}

	fmt.Printf("All available modules:\n\n")
	cs.ListModules(os.Stdout, products)
	return nil
}

// runListProducts lists all available products and their metadata
func runListProducts(_ *cli.Context) error {
	products, err := requestProducts()
	if err != nil {
		return err
	}

	fmt.Printf("All available products:\n\n")
	cs.ListProducts(os.Stdout, products, "none")
	return nil
}
