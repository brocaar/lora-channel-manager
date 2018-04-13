---
title: LoRa Channel Manager
menu:
  main:
    parent: overview
    weight: 1
---

# LoRa Channel Manager

**This component has been deprecated and has been merged into the
[LoRa Gateway Bridge](https://www.loraserver.io/lora-gateway-bridge/)!**

LoRa Channel Manager periodically reads channel-configuration from [LoRa Server](/loraserver/),
updates the [packet-forwarder](https://github.com/lora-net/packet_forwarder)
configuration in case of updates, and restarts the packet-forwarder process
in case of any changes.

## How it works

At configured intervals it calls the [LoRa Server](/loraserver/) API
to fetch the channel-configuration for a given gateway MAC. In case of an
update, it reads a base configuration, updates the channel related keys
(`radio_`, `chan_`) and `gateway_ID` key and writes this as a new JSON file.
After writing the new JSON configuration file, it will issue the configured
packet-forwarder restart command.
