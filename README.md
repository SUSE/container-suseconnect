# container-suseconnect [![Golang Tests](https://github.com/SUSE/container-suseconnect/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/SUSE/container-suseconnect/actions/workflows/ci.yaml) [![GoDoc](https://godoc.org/github.com/SUSE/container-suseconnect?status.png)](https://godoc.org/github.com/SUSE/container-suseconnect)

container-suseconnect provides a [ZYpp service plugin](https://doc.opensuse.org/projects/libzypp/HEAD/zypp-plugins.html#plugin-services),
a [ZYpp url resolver plugin](https://doc.opensuse.org/projects/libzypp/HEAD/zypp-plugins.html#plugin-url-resolver)
and command line interface.

It gives access to repositories during docker build and run using the host machine credentials.

## Command line interface

The application runs as ZYpp service per default if the name of the executable
is `container-suseconnect-zypp`. It will run as a ZYpp URL resolver if the name
of the executable is `susecloud`. In every other case it assumes that a real
user executes the application. The help output of container-suseconnect shows
all available commands and indicates the current defaults:

```bash
container-suseconnect -h
NAME:
   container-suseconnect - Access zypper repositories from within containers

USAGE:
   This application can be used to retrieve basic metadata about SLES
   related products and module extensions.

   Please use the 'list-products' subcommand for listing all currently
   available products including their repositories and a short description.

   Use the 'list-modules' subcommand for listing available modules, where
   their 'Identifier' can be used to enable them via the ADDITIONAL_MODULES
   environment variable during container creation/run. When enabling multiple
   modules the identifiers are expected to be comma-separated.

   The 'zypper' subcommand runs the application as zypper plugin and is only
   intended to use for debugging purposes.

VERSION:
   2.3.0

COMMANDS:
   list-products, lp  List available products (default)
   list-modules, lm   List available modules
   zypper, z, zypp    Run the zypper service plugin
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

COPYRIGHT:
   © 2020 SUSE LCC
```

## Logging

By default, this program will log everything into the
`/var/log/suseconnect.log` file. Otherwise, you can optionally specify a custom
path with the `SUSECONNECT_LOG_FILE` environment variable. To retrieve the
contents of the log, you will have to mount the needed volume in your host.

You can do this with either the
[VOLUME](https://docs.docker.com/engine/reference/builder/#volume) instruction
in your Dockerfile, or with the `-v` parameter when running `docker run`. Read
more about this [here](https://docs.docker.com/storage/volumes/).

Finally, if neither the default `/var/log/suseconnect.log` file nor the file
specified through the `SUSECONNECT_LOG_FILE` environment variable are writable,
then this program will log to the standard error output by default.

## Example Dockerfile

Creating a SLE 15 Image

The following Docker file creates a simple Docker image based on SLE 15:

```Dockerfile
FROM registry.suse.com/suse/sle15:latest

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

When the Docker host machine is registered against an internal SMT server, the
Docker image requires the SSL certificate used by SMT:

```Dockerfile
FROM registry.suse.com/suse/sle15:latest

# Import the crt file of our private SMT server
ADD http://smt.test.lan/smt.crt /etc/pki/trust/anchors/smt.crt
RUN update-ca-certificates

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

Creating a Custom SLE 15 SP5 Image

The following Docker file creates a simple Docker image based on SLE 15 SP5:

```Dockerfile
FROM registry.suse.com/suse/sle15:15.5

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

After registering the host machine against an internal SMT server, the Docker
image requires the SSL certificate used by SMT:

```Dockerfile
FROM registry.suse.com/suse/sle15:15.5

# Import the crt file of our private SMT server
ADD http://smt.test.lan/smt.crt /etc/ssl/certs/smt.pem
RUN c_rehash /etc/ssl/certs

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in vim
```

All recommended package modules are enabled by default. It is possible to
enable additionally non-recommended modules via the `identifier` by setting the
environment variable `ADDITIONAL_MODULES`. When enabling multiple modules the
identifiers are expected to be comma-separated:

```Dockerfile
FROM registry.suse.com/suse/sle15:latest

ENV ADDITIONAL_MODULES sle-module-desktop-applications,sle-module-development-tools

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in gvim
```

Examples taken from
<https://documentation.suse.com/sles/12-SP4/html/SLES-all/docker-building-images.html#Customizing-Pre-build-Images>

### Building images on SLE systems registered with RMT or SMT

When the host system used for building the docker images is registered against
RMT or SMT it is only possible to build containers for the same SLE code base
as the host system is running on. I.e. if you docker host is a SLE15 system you
can only build SLE15 based images out of the box.

If you want to build for a different SLE version than what is running on the
host machine you will need to inject matching credentials for that target
release into the build. For details on how to achieve that please follow the
steps outlined in the [Building images on non SLE
distributions](#building-images-on-non-sle-distributions)
section.

### Building images on-demand SLE instances in the public cloud

When building container images on SLE instances that were launched as so-called
"on-demand" or "pay as you go" instances on a Public Cloud (as AWS, GCE or
Azure) some additional steps have to be performed.

For installing packages and updates the "on-demand" public cloud instance are
connected to a public cloud specific update infrastructure which is based
around RMT servers operated by SUSE on the various Public Cloud Providers.

To be able access this update infrastructure instances need to perform
additional steps to locate the required services and authenticate with them.
More details on this are outlined in a
[Blog Series on suse.com starting here](https://suse.com/c/a-new-update-infrastructure-for-the-public-cloud/).

In order to build containers on this type of instances a new service was
introduced, that service is called `containerbuild-regionsrv` and will be
available in the public cloud images provided through the Marketplaces of the
various Public Cloud Providers. So before building an image this service has
to be started on the public cloud instance

```bash
systemctl start containerbuild-regionsrv
```

In order to have it started automatically (e.g. after reboot) please use:

```bash
systemctl enable containerbuild-regionsrv
```

The zypper plugins provided by `container-suseconnect` will then connect to
this service for getting authentication details and information about which
update server to talk to. In order for that to work the container has to be
built with host networking enabled. I.e. you need to call `docker build` with
`--network host`:

```bash
docker build --network host <builddir>
```

Since update infrastructure in the Public Clouds is based upon RMT, the same
restrictions with regard to building SLE images for SLE versions differing from
the SLE version of the host apply here as well. (See above)

### Building images on non SLE distributions

It is possible to build SLE based docker images on other distributions as well.

For that the following two files from a base SLE system are needed:

- `/etc/SUSEConnect` for the API access point
- `/etc/zypp/credentials.d/SCCcredentials` for providing the user credentials

These files can be copied from a SLE machine that has been successfully
registered at the SUSE Customer Center, RMT or SMT.

A Docker version of 18.09 or above is needed to provide a secure way to mount
the credentials into the image build process.

A `Dockerfile` for building a SLE15 image which contains `go` would then look like the following:

```Dockerfile
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

After creating the file you can build the image by executing:

```bash
docker build -t sle15-go \
    --build-arg ADDITIONAL_MODULES=sle-module-development-tools \
    --secret id=SUSEConnect,src=SUSEConnect \
    --secret id=SCCcredentials,src=SCCcredentials \
    .
```

At the time of writing (docker 18.09.0) it is necessary to enable the Docker
BuildKit by setting the environment variable `DOCKER_BUILDKIT`. The
`ADDITIONAL_MODULES` are used here to enable all needed repositories. After the
image has been built the package should be usable within the docker image:

```bash
docker run sle15-go go version
go version go1.9.7 linux/amd64
```

Please keep in mind that it is not possible to use `container-suseconnect` or
`zypper` within the container after the build, because the secrets are not
available any more.


### Alternative approach when using podman as the container runtime

When using `podman` as the container runtime, add the following to
`/etc/containers/mounts.conf` (when building containers as `root`) or
`~/.config/containers/mounts.conf` (when building containers as a regular
user):

```text
<path_on_host>/SUSEConnect:/etc/SUSEConnect
<path_on_host>/SCCcredentials:/etc/zypp/credentials.d/SCCcredentials
```

No further change are needed in `Dockerfile`.

### Obtaining the SUSEConnect and SCCcredentials secrets

Ideally the host system on which your container builds are running is
registered. In all other cases, you need to manually start the container you
would like to register and execute `SUSEConnect` inside:

```bash
SUSEConnect -e <youremailaddress> -r <yourregistrationcode>
cat /etc/SUSEconnect
cat /etc/zypp/credentials.d/SCCcredentials
```

As a last resort, for example for interactive use of the container, you can
pass the obtained username and password via environment variables

```bash
env SCC_CREDENTIAL_USERNAME=<credential_username>
env SCC_CREDENTIAL_PASSWORD=<credential_password>
```

These can be passed in the container also at start time via the `docker run -e ENV=key` options.

# License

Licensed under the Apache License, Version 2.0. See
[LICENSE](./LICENSE) for the full license text.
