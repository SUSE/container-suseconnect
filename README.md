# container-suseconnect [![Build Status](https://travis-ci.org/SUSE/container-suseconnect.svg?branch=master)](https://travis-ci.org/SUSE/container-suseconnect) [![GoDoc](https://godoc.org/github.com/SUSE/container-suseconnect?status.png)](https://godoc.org/github.com/SUSE/container-suseconnect)

container-suseconnect is a [ZYpp service](http://doc.opensuse.org/projects/libzypp/HEAD/zypp-plugins.html)
and command line interface.

It gives access to repositories during docker build and run using the host machine credentials.

## Command line interface

The application runs as ZYpp service per default if the name of the executable
is `container-suseconnect-zypp`. In every other case it assumes that a real user
executes the application. The help output of container-suseconnect shows all
available commands and indicates which one is the current default:

```
container-suseconnect -h
NAME:
   container-suseconnect

USAGE:
   This application can be used to retrieve basic metadata about SLES
   related products and module extensions.

   Please use the 'list-products' subcommand for listing all currently
   available products including their repositories and a short description.

   Use the 'list-modules' subcommand for listing available modules, where
   their 'Identifier' can be used to enable them via the ADDITIONAL_MODULES
   environment variable during container creation/run.

   The 'zypper' subcommand runs the application as zypper plugin and is only
   intended to use for debugging purposes.

VERSION:
   2.1.0

COMMANDS:
     list-products, lp  List available products (default)
     list-modules, lm   List available modules
     zypper, z, zypp    Run the zypper service plugin
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   © 2018 SUSE LCC
```

## Logging

By default, this program will log everything into the
`/var/log/suseconnect.log` file. Otherwise, you can optionally specify a custom
path with the `SUSECONNECT_LOG_FILE` environment variable. To retrieve the
contents of the log, you will have to mount the needed volume in your host.
You can do this with either the
[VOLUME](https://docs.docker.com/reference/builder/#volume) instruction in your
Dockerfile, or with the `-v` parameter when running `docker run`. Read more
about this [here](https://docs.docker.com/userguide/dockervolumes/).

Finally, if neither the default `/var/log/suseconnect.log` file nor the file
specified through the `SUSECONNECT_LOG_FILE` environment variable are writable,
then this program will default to the standard error.

## Example Dockerfiles:

Creating a Custom SLE 12 Image

The following Docker file creates a simple Docker image based on SLE 12:

```
FROM suse/sles12:latest

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

When the Docker host machine is registered against an internal SMT server, the Docker image requires the SSL certificate used by SMT:

```
FROM suse/sles12:latest

# Import the crt file of our private SMT server
ADD http://smt.test.lan/smt.crt /etc/pki/trust/anchors/smt.crt
RUN update-ca-certificates

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

Creating a Custom SLE 11 SP3 or SP4 Image

The following Docker file creates a simple Docker image based on SLE 11 SP3:

```
FROM suse/sles11sp3:latest

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

When the Docker host machine is registered against an internal SMT server, the Docker image requires the SSL certificate used by SMT:

```
FROM suse/sles11sp3:latest

# Import the crt file of our private SMT server
ADD http://smt.test.lan/smt.crt /etc/ssl/certs/smt.pem
RUN c_rehash /etc/ssl/certs

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

All recommended package modules are enabled by default. It is possible to enable
additionally non-recommended modules via the `identifier` by setting the
environment variable `ADDITIONAL_MODULES`:

```
FROM registry.suse.com/suse/sle15:latest

ENV ADDITIONAL_MODULES sle-module-desktop-applications,sle-module-development-tools

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in gvim
```

Examples taken from https://www.suse.com/documentation/sles-12/book_sles_docker/data/customizing_pre-build_images.html

### Building images on non SLE distributions

It is possible to build SLE based docker images on other distributions as well.
For that the following two files from a base SLE system are needed:

- `/etc/SUSEConnect` for the API access point
- `/etc/zypp/credentials.d/SCCcredentials` for providing the user credentials

These files can be copied from a SLE machine that has been successfully
registered at the SUSE Customer Center, RMT or SMT.

A Docker version of 18.09 or above is needed to provide a secure way to mount
the credentials into the image build process. This version can be installed for
example on LEAP15 (x86_64) via:

```
> sudo zypper addrepo https://download.opensuse.org/repositories/Virtualization:containers/openSUSE_Leap_15.0/Virtualization:containers.repo
> sudo zypper in docker-18.09.0_ce-lp150.304.1.x86_64
```

A `Dockerfile` for building a SLE15 image which contains `go` would then look like
this:

```
# syntax=docker/dockerfile:1.0.0-experimental
FROM registry.suse.com/suse/sle15:latest

ARG ADDITIONAL_MODULES
RUN --mount=type=secret,id=SUSEConnect,required \
    --mount=type=secret,id=SCCcredentials,required \
    zypper -n --gpg-auto-import-keys in go
```

This file mounts all necessary secrets into the image during the build process
and removes it afterwards. Please note that the first line of the `Dockerfile`
is mandatory at the time of writing. Both files `SUSEConnect` and
`SCCcredentials` needs to be available beside the `Dockerfile`.

After the file creation the image can be built by executing:

```bash
> export DOCKER_BUILDKIT=1
> docker build -t sle15-go \
    --build-arg ADDITIONAL_MODULES=PackageHub \
    --secret id=SUSEConnect,src=SUSEConnect \
    --secret id=SCCcredentials,src=SCCcredentials \
    .
```

At the time of writing (docker 18.09.0) it is necessary to enable the Docker
BuildKit by setting the environment variable `DOCKER_BUILDKIT`. The
`ADDITIONAL_MODULES` are used here to enable all needed repositories. After the
image has been built the package should be usable within the docker image:

```bash
> docker run sle15-go go version
go version go1.9.7 linux/amd64
```

Please keep in mind that it is not possible to use `container-suseconnect` or
`zypper` within the container after the build, because the secrets are not
available any more.

# License

Licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/SUSE/Portus/blob/master/LICENSE) for the full
license text.
