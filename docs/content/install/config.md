---
title: Configuration
menu:
  main:
    parent: install
    weight: 5
---

# Configuration

To list all configuration options, start `lora-channel-manager` with the
`--help` flag. This will display:

```text
GLOBAL OPTIONS:
   --gw-mac value                mac address of the gateway [$GW_MAC]
   --gw-server value             hostname:ip of the gateway api server (default: "127.0.0.1:8002") [$GW_SERVER]
   --gw-client-ca-cert value     ca certificate used by the gateway-server client (optional) [$GW_CLIENT_CA_CERT]
   --gw-client-tls-cert value    tls certificate used by the gateway-server client (optional) [$GW_CLIENT_TLS_CERT]
   --gw-client-tls-key value     tls key used by the gateway-server client (optional) [$GW_CLIENT_TLS_KEY]
   --gw-client-jwt-token value   jwt token used by the gateway-server client for authentication (issued by LoRa Server) [$GW_CLIENT_JWT_TOKEN]
   --base-config-file value      path to the base configuration file [$BASE_CONFIG_FILE]
   --output-config-file value    path to the output configuration file [$OUTPUT_CONFIG_FILE]
   --pf-restart-command value    command which must be executed on configuration changes to restart the packet-forwarder [$PF_RESTART_COMMAND]
   --config-poll-interval value  interval between polling new configuration (default: 5m0s) [$CONFIG_POLL_INTERVAL]
   --help, -h                    show help
   --version, -v                 print the version
```

Both cli arguments and environment-variables can be used to pass configuration
options.

## Gateway MAC

The gateway MAC address must be given in HEX format, e.g. `0102030405060708`.

## Gateway API server

This is the `IP:PORT` pointing to the gateway API server. This API server is
exposed by the [LoRa Server](/loraserver/) service.

## Configuration files

LoRa Channel Manager reads a base configuration file (`--base-config-file`),
updates the `radio_`, `chan_` and `gateway_ID` keys and writes the end-result
into `--output-config-file`. The base configuration file must already contain
all other configuration values.

**Note:** the file to which `--output-config-file` point will be overwritten!

### Option one

When your current setup uses a `global_conf.json` and `local_conf.json` file,
combine these two files into a third file containing all configuration and use
this file as your `--base-config-file`. The `--output-config-file` could then
be the path to your `local_conf.json` file.

### Option two

An alternative way is to set both `--base-config-file` and
`--output-config-file` to the path of the `global_conf.json`. That way you are
still able to make overwrites by using a `local_conf.json` file. As this means
that your original `global_conf.json` file will be overwritten, **make sure you
keep a backup of the original `global_conf.json`!**.


## JWT token

The JWT token (`--gw-client-jwt-token`) must be set to authenticate the gateway
to the API. This token can be retrieved through the [LoRa App Server](/lora-app-server/)
web-interface or through the [LoRa Server](/loraserver/) API.

## Packet-forwarder restart command

This command configured by the `--pf-restart-command` will be executed by
LoRa Channel Manager each time there is a configuration update. The command thati
needs to be configured is gateway dependent. Please refer the manual of your
gateway.
