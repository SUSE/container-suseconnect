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
	"bufio"
	"io"
	"os"
	"strings"
)

// The Configuration interface allows us to fetch information from
// configuration files located somewhere in the system and that follow a
// grammar like: 'key' 'separator' 'value'.
type Configuration interface {
	// Returns the character separating the key from the value.
	separator() byte

	// Returns a slice of possible locations for the configuration file. Note
	// that the order matters: the first elements will be evaluated first
	// during parse time.
	locations() []string

	// This function will be called if this app could not found any file on the
	// specified locations. Returns a bool: true if it has been handled by the
	// implementer of this interface, false if it should fail.
	onLocationsNotFound() bool

	// Called after a line with relevant information has been successfully
	// parsed. The key and the value are guaranteed to be already trimmed
	// by the caller.
	setValues(key, value string)

	// Checks that should be done after the parsing has been performed. It's
	// assumed that the error returned has already been logged.
	afterParseCheck() error
}

// From the given slice of locations, return the first location that actually
// exists on the system. It returns an empty string on error.
func getLocationPath(locations []string) string {
	for _, path := range locations {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// ReadConfiguration reads the configuration and updates the given object.
func ReadConfiguration(config Configuration) error {
	path := getLocationPath(config.locations())
	if path == "" {
		// Leave early if locations could not be found but it can be handled by
		// the implementer.
		if config.onLocationsNotFound() {
			return nil
		}
		return loggedError("Warning: SUSE credentials not found: %v - automatic handling of repositories not done.", config.locations())
	}

	file, err := os.Open(path)
	if err != nil {
		return loggedError("Can't open %s file: %v", path, err.Error())
	}
	defer file.Close()

	return parse(config, file)
}

// Parses the contents given by the reader and updated the given configuration.
func parse(config Configuration, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// Comments & empty lines.
		if strings.IndexAny(scanner.Text(), "#-") == 0 {
			continue
		}
		if scanner.Text() == "" {
			continue
		}

		// Each line should be constructed as 'key' 'separator' 'value'.
		line := scanner.Text()
		parts := strings.SplitN(line, string(config.separator()), 2)
		if len(parts) != 2 {
			return loggedError("Can't parse line: %v", line)
		}

		// And finally trim the key and the value and pass it to the config.
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		config.setValues(key, value)
	}

	// Final checks.
	if err := scanner.Err(); err != nil {
		return loggedError("Error when scanning configuration: %v", err)
	}
	if err := config.afterParseCheck(); err != nil {
		return err
	}
	return nil
}
