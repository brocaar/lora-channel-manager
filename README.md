# LoRa Gateway Config

[![Build Status](https://travis-ci.org/brocaar/lora-gateway-config.svg?branch=master)](https://travis-ci.org/brocaar/lora-gateway-config)
[![Gitter chat](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/loraserver/lora-gateway-config)

LoRa Gateway Config reads periodically channel-configuration from [LoRa Server](https://github.com/brocaar/loraserver/)
updates the [packet-forwarder](https://github.com/lora-net/packet_forwarder)
configuration in case of updates, and restarts the packet-forwarder process.

## Documentation

Please refer to [https://docs.loraserver.io/lora-gateway-config/](https://docs.loraserver.io/lora-gateway-config/).

## Building from source

The easiest way to get started is by using the provided [docker-compose](https://docs.docker.com/compose/)
environment. To start a bash shell within the docker-compose environment,
execute the following command from the root of this project:

```bash
docker-compose run --rm gwconfig bash
```

A few example commands that you can run:

```bash
# cleanup workspace
make clean

# run the tests
make test

# compile (this will also compile the ui and generate the static files)
make build

# cross-compile for Linux ARM
GOOS=linux GOARCH=arm make build

# cross-compile for Windows AMD64
GOOS=windows BINEXT=.exe GOARCH=amd64 make build

# build the .tar.gz file
make package

# build the .tar.gz file for Linux ARM
GOOS=linux GOARCH=arm make package

# build the .tar.gz file for Windows AMD64
GOOS=windows BINEXT=.exe GOARCH=amd64 make package
```

Alternatively, you can run the same commands from any working [Go](https://golang.org)
environment. As all requirements are vendored, there is no need to go get
these. Make sure you have Go 1.7.x installed and that you clone this
repository to `$GOPATH/src/github.com/brocaar/lora-gateway-config`.

## Contributing

There are a couple of ways to get involved:

* Join the discussions and share your feedback at [https://gitter.io/loraserver/lora-gateway-config](https://gitter.io/loraserver/lora-gateway-config)
* Report bugs or make feature-requests by opening an issue at [https://github.com/brocaar/lora-gateway-config/issues](https://github.com/brocaar/lora-gateway-config/issues)
* Fix issues or improve documentation by creating pull-requests

## License

LoRa Gateway Config is distributed under the MIT license. See also
[LICENSE](https://github.com/brocaar/lora-gateway-config/blob/master/LICENSE).
