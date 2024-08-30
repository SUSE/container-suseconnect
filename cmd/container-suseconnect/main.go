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

	cs "github.com/SUSE/container-suseconnect/internal"
	"github.com/SUSE/container-suseconnect/internal/regionsrv"
	"github.com/urfave/cli/v2"
)

func main() {
	cs.SetLoggerOutput()

	// Set the basic CLI metadata
	app := cli.NewApp()
	app.Copyright = fmt.Sprintf("© %d SUSE LCC", time.Now().Year())
	app.Name = "container-suseconnect"
	app.Version = cs.Version
	app.Usage = "Access zypper repositories from within containers"
	app.UsageText =
		`This application can be used to retrieve basic metadata about SLES
   related products and module extensions.

   Please use the 'list-products' subcommand for listing all currently
   available products including their repositories and a short description.

   Use the 'list-modules' subcommand for listing available modules, where
   their 'Identifier' can be used to enable them via the ADDITIONAL_MODULES
   environment variable during container creation/run. When enabling multiple
   modules the identifiers are expected to be comma-separated.

   The 'zypper' subcommand runs the application as zypper plugin and is only
   intended to use for debugging purposes.`

	// Switch the default application behavior in relation to the basename
	defaultUsageAdditionZypp := ""
	defaultUsageAdditionListProducts := ""
	switch filepath.Base(os.Args[0]) {
	case "container-suseconnect-zypp":
		app.Action = runZypperPlugin
		defaultUsageAdditionZypp = " (default)"
	case "susecloud":
		app.Action = runZypperURLResolver
		defaultUsageAdditionZypp = " (default)"
	default:
		app.Action = runListProducts
		defaultUsageAdditionListProducts = " (default)"
	}

	// Set additional actions, which are always available
	app.Commands = cli.Commands{
		{
			Name:    "list-products",
			Aliases: []string{"lp"},
			Usage: fmt.Sprintf("List available products%v",
				defaultUsageAdditionListProducts),
			Action: runListProducts,
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
			Usage: fmt.Sprintf("Run the zypper service plugin%v",
				defaultUsageAdditionZypp),
			Action: runZypperPlugin,
		},
	}

	// Run the application
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// requestProducts collects a slice of products for the currently available
// environment
func requestProducts() ([]cs.Product, error) {
	credentials := cs.Credentials{}
	suseConnectData := cs.SUSEConnectData{}

	// read config from "containerbuild-regionsrv" service, if that service is
	// running, we're running inside a public cloud instance in that case read
	// config from "mounted" files if the service is not available
	if err := regionsrv.ServerReachable(); err == nil {
		log.Printf("containerbuild-regionsrv reachable, reading config\n")
		cloudCfg, err := regionsrv.ReadConfigFromServer()
		if err != nil {
			return nil, err
		}
		credentials.Username = cloudCfg.Username
		credentials.Password = cloudCfg.Password
		credentials.InstanceData = cloudCfg.InstanceData
		suseConnectData.SccURL = "https://" + cloudCfg.ServerFqdn
		suseConnectData.Insecure = false

		if cloudCfg.Ca != "" {
			regionsrv.SafeCAFile(cloudCfg.Ca)
		}
		regionsrv.UpdateHostsFile(cloudCfg.ServerFqdn, cloudCfg.ServerIP)
	} else {
		if err := cs.ReadConfiguration(&credentials); err != nil {
			return nil, err
		}
		if err := cs.ReadConfiguration(&suseConnectData); err != nil {
			return nil, err
		}
	}

	installedProduct, err := cs.GetInstalledProduct()
	if err != nil {
		return nil, err
	}
	log.Printf("Installed product: %v\n", installedProduct)
	log.Printf("Registration server set to %v\n", suseConnectData.SccURL)

	products, err := cs.RequestProducts(suseConnectData, credentials, installedProduct)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// Read the arguments as given by zypper on the stdin and print into stdout the
// response to be used.
func runZypperURLResolver(_ *cli.Context) error {
	if err := regionsrv.ServerReachable(); err != nil {
		return fmt.Errorf("could not reach build server from the host: %v", err)
	}

	input, err := regionsrv.ParseStdin()
	if err != nil {
		return fmt.Errorf("could not parse input: %s", err)
	}

	return regionsrv.PrintResponse(input)
}

// runZypperPlugin runs the application in zypper plugin mode, which dumps
// all available repositories for the installed product. Additional modules
// can be specified via the `ADDITIONAL_MODULES` environment variable, which
// reflect the module `identifier`.
func runZypperPlugin(_ *cli.Context) error {
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
