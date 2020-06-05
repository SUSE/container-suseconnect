#!/usr/bin/env python
# -*- encoding: utf-8 -*-

# Copyright (c) 2020 SUSE LLC. All rights reserved.
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

import os

from setuptools import setup


def version():
    '''Return the version. Pass when version file could not be parsed'''

    pwd = os.path.dirname(os.path.abspath(__file__))
    try:
        with open(os.path.join(pwd, '..', 'internal/' 'version.go')) as version_file:
            data = version_file.read()
            matches = re.findall(r'Version = \"(.+?)\"', data)
            if matches == None or len(matches) != 1:
                return ''
            return matches[0]
    except FileNotFoundError:
        return ''


setup(
    name="containerbuild-regionsrv",
    version=version(),
    author="SUSE Containers Team",
    author_email="containers@suse.com",
    description="Services that provides the needed data from cloud-regionsrv for container-suseconnect",
    long_description="TCP server that listens on a given port and replies back with the needed information for authenticating into SMT servers running in the Public Clouds.",
    license="Apache License 2.0",
    keywords="SUSEConnect",
    url="https://github.com/SUSE/container-suseconnect",
    packages=['containerbuild-regionsrv'],
    classifiers=[
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache License 2.0',
        'Operating System :: POSIX :: Linux',
    ], data_files=[],
    entry_points={
        'console_scripts': [
            'containerbuild-regionsrv = containerbuild_regionsrv.containerbuild_regionsrv:main'
        ]
    }
)
