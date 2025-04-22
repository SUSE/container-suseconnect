// Copyright (c) 2020 SUSE LLC. All rights reserved.
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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

// testServer implements a net.Listener with some mocking attributes.
type testServer struct {
	server       net.Listener
	bootstrapped chan bool
	badResponse  bool
	response     ContainerBuildConfig
}

// run this testServer by taking the mocking attributes into account.
func (ts *testServer) run() (err error) {
	ts.server, err = net.Listen("tcp", "0.0.0.0:7956")
	if err != nil {
		err = fmt.Errorf("could not start test server: %v", err)
		ts.bootstrapped <- true

		return
	}

	for {
		ts.bootstrapped <- true

		conn, err := ts.server.Accept()
		if err != nil {
			break
		}

		if conn == nil {
			err = errors.New("could not create connection")
			break
		}

		if ts.badResponse {
			io.WriteString(conn, "{")
		}

		b, err := json.Marshal(ts.response)
		if err != nil {
			break
		}

		io.WriteString(conn, string(b))
		conn.Close()
	}

	return
}

// close the server if needed.
func (ts *testServer) close() error {
	if ts.server == nil {
		return nil
	}

	return ts.server.Close()
}

// Execute fn by suppressing any logger output.
func withSuppressedLog(fn func()) {
	log.SetOutput(io.Discard)
	fn()
	log.SetOutput(os.Stdout)
}

// Test suite below

func TestReadConfigFromServerFailConnection(t *testing.T) {
	withSuppressedLog(func() {
		_, err := ReadConfigFromServer()
		if !strings.Contains(err.Error(), "connection refused") {
			t.Fatalf("should be a connection refused error, got '%v'", err)
		}
	})
}

func TestReadConfigFromServerBadResponse(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{
			bootstrapped: make(chan bool, 1),
			badResponse:  true,
		}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		_, err := ReadConfigFromServer()
		if !strings.Contains(err.Error(), "invalid character '{' looking for beginning of object key string") {
			t.Fatalf("should be a 'invalid character '{' looking for beginning of object key string', got '%v'", err)
		}
	})
}

func TestEmptyResponseFromServer(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{
			bootstrapped: make(chan bool, 1),
			response:     ContainerBuildConfig{},
		}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		_, err := ReadConfigFromServer()
		if !strings.Contains(err.Error(), "empty response from the server") {
			t.Fatalf("should be a 'empty response from the server', got '%v'", err)
		}
	})
}

func TestValidResponse(t *testing.T) {
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

		cfg, err := ReadConfigFromServer()
		if err != nil {
			t.Fatalf("should be nil but got: %v", err)
		}

		if cfg.InstanceData != ts.response.InstanceData {
			t.Fatalf("Expected '%v', got '%v'", cfg.InstanceData, ts.response.InstanceData)
		}
	})
}

func TestServerReachableNope(t *testing.T) {
	withSuppressedLog(func() {
		err := ServerReachable()
		if !strings.Contains(err.Error(), "connection refused") {
			t.Fatalf("should be a connection refused error, got '%v'", err)
		}
	})
}

func TestServerReachableSuccessful(t *testing.T) {
	withSuppressedLog(func() {
		ts := &testServer{bootstrapped: make(chan bool, 1)}
		defer ts.close()

		go ts.run()
		<-ts.bootstrapped

		err := ServerReachable()
		if err != nil {
			t.Fatalf("should be nil but got: %v", err)
		}
	})
}
