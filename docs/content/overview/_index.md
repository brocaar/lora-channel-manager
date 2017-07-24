---
title: LoRa Gateway Config
menu:
  main:
    parent: overview
    weight: 1
---

## LoRa Gateway Config

LoRa Gateway Config periodically reads channel-configuration from [LoRa Server](/loraserver/)
adn updates the [packet-forwarder](https://github.com/lora-net/packet_forwarder)
configuration in case of updates, and restarts the packet-forwarder process.

### How it works

At configured intervals it calls the [LoRa Server](/loraserver/) API
to fetch the channel-configuration for a given gateway MAC. In case of an
update, it reads a base configuration, updates the channel related keys
(`radio_`, `chan_`) and `gateway_ID` key and writes this as a new JSON file.
After writing the new JSON configuration file, it will issue the configured
packet-forwarder restart command.
