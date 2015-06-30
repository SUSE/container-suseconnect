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

# License

Licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/SUSE/Portus/blob/master/LICENSE) for the full
license text.
