# LoRa Channel Manager

[![Build Status](https://travis-ci.org/brocaar/lora-channel-manager.svg?branch=master)](https://travis-ci.org/brocaar/lora-channel-manager)

**This component has been deprecated and has been merged into the
[LoRa Gateway Bridge](https://www.loraserver.io/lora-gateway-bridge/)!**

LoRa Channel Manager periodically reads channel-configuration from [LoRa Server](https://github.com/brocaar/loraserver/),
updates the [packet-forwarder](https://github.com/lora-net/packet_forwarder)
configuration in case of updates, and restarts the packet-forwarder process
in case of any changes.

## Documentation

Please refer to [https://docs.loraserver.io/lora-channel-manager/](https://docs.loraserver.io/lora-channel-manager/).

## Building from source

The easiest way to get started is by using the provided [docker-compose](https://docs.docker.com/compose/)
environment. To start a bash shell within the docker-compose environment,
execute the following command from the root of this project:

```bash
docker-compose run --rm channelmanager bash
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
repository to `$GOPATH/src/github.com/brocaar/lora-channel-manager`.

## Contributing

There are a couple of ways to get involved:

* Join the discussions at [https://forum.loraserver.io](https://forum.loraserver.io/)
* Report bugs or make feature-requests by opening an issue at [https://github.com/brocaar/lora-channel-manager/issues](https://github.com/brocaar/lora-channel-manager/issues)
* Fix issues or improve documentation by creating pull-requests

## License

LoRa Channel Manager is distributed under the MIT license. See also
[LICENSE](https://github.com/brocaar/lora-channel-manager/blob/master/LICENSE).
