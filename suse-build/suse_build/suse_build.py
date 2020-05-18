#!/usr/bin/env python
# -*- encoding: utf-8 -*-

# Copyright (c) 2019 SUSE LLC.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from cloudregister.registerutils import get_instance_data, get_smt, get_config
import base64
import os
import socketserver
import json

class SuseBuildTCPServer(socketserver.BaseRequestHandler):
    """
    A TCP server that submits configuration details that are relevant to
    SUSEConnect.
    """

    def instance_data_header(self):
        """
        Returns the instance data as retrieved from the SMT server.
        """

        instance_data = bytes(get_instance_data(get_config()), 'utf-8')
        return base64.b64encode(instance_data).decode()

    def smt_fqdn(self):
        """
        Get the FQDN from the SMT being used.
        """

        return get_smt().get_FQDN()

    def handle(self):
        """
        This is the method being called for each request. It returns a JSON response
        with all the relevant information.
        """

        resp = {
            'instance-data': self.instance_data_header(),
            "server": self.smt_fqdn()
        }
        self.request.sendall(bytes(json.dumps(resp), 'utf-8'))

def main():
    """
    main entry point of the program.
    """

    ip = os.getenv("SUSE_BUILD_IP", '0.0.0.0')
    port = int(os.getenv("SUSE_BUILD_PORT", 7956))

    with socketserver.TCPServer((ip, port), SuseBuildTCPServer) as server:
        server.serve_forever()

if __name__ == "__main__":
    main()
