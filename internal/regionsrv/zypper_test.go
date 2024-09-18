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

package regionsrv

import (
	"os"
	"strings"
	"testing"

	"github.com/mssola/capture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStdinSuccessful(t *testing.T) {
	content := []byte("RESOLVEURL\nkey:value1\nanother:value2")
	tmp, err := os.CreateTemp("", "file")
	require.Nil(t, err)

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	_, err = tmp.Write(content)
	require.Nil(t, err)

	_, err = tmp.Seek(0, 0)
	require.Nil(t, err)

	// stdin switcheroo
	old := os.Stdin
	defer func() {
		os.Stdin = old
	}()
	os.Stdin = tmp

	params, err := ParseStdin()
	require.Nil(t, err)
	assert.Len(t, params, 2)
	assert.Equal(t, "value1", params["key"])
	assert.Equal(t, "value2", params["another"])
}

func TestParseStdinBadInput(t *testing.T) {
	old := os.Stdin
	defer func() {
		os.Stdin = old
	}()
	os.Stdin = nil

	_, err := ParseStdin()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid argument")
}

func TestPrintResponseBadResponseFromServer(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{
			bootstrapped: make(chan bool, 1),
			response:     ContainerBuildConfig{},
		}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		params := map[string]string{}

		err := PrintResponse(params)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "empty response from the server")
	})
}

func TestPrintResponseNoCredentials(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{
			bootstrapped: make(chan bool, 1),
			response: ContainerBuildConfig{
				InstanceData: "instance data",
			},
		}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		params := map[string]string{}

		err := PrintResponse(params)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "no credentials given")
	})
}

func TestPrintFromConfiguration(t *testing.T) {
	withSuppressedLog(func() {
		res := capture.All(func() {
			printFromConfiguration("/path", &ContainerBuildConfig{
				InstanceData: "instance data",
				ServerFqdn:   "test.fqdn.com",
				ServerIP:     "1.1.1.1",
				Username:     "banjo",
				Password:     "kazooie",
				Ca:           "ca",
			})
		})

		lines := strings.Split(string(res.Stdout), "\n")
		expected := []string{
			"RESOLVEDURL",
			"X-Instance-Data:instance data",
			"",
			"https://banjo:kazooie@test.fqdn.com/path\x00",
		}

		require.Equal(t, len(expected), len(lines))

		for k, v := range expected {
			assert.Equal(t, v, lines[k])
		}
	})
}
