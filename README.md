# container-suseconnect [![Build Status](https://travis-ci.org/SUSE/container-suseconnect.svg?branch=master)](https://travis-ci.org/SUSE/container-suseconnect) [![GoDoc](https://godoc.org/github.com/SUSE/container-suseconnect?status.png)](https://godoc.org/github.com/SUSE/container-suseconnect)

container-suseconnect is a [ZYpp service](http://doc.opensuse.org/projects/libzypp/HEAD/zypp-plugins.html).

Gives access to repositories during docker build and run using the host machine credentials.

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
additionally non-recommended modules via the environment variable
`ADDITIONAL_MODULES`:
```
FROM registry.suse.com/suse/sle15:latest

ENV ADDITIONAL_MODULES sle-module-desktop-applications,sle-module-development-tools

RUN zypper --gpg-auto-import-keys ref -s
RUN zypper -n in gvim
```
Examples taken from https://www.suse.com/documentation/sles-12/book_sles_docker/data/customizing_pre-build_images.html


# License

Licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/SUSE/Portus/blob/master/LICENSE) for the full
license text.
