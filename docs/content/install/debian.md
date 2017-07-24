---
title: Debian / Rasberry Pi
menu:
  main:
    parent: install
    weight: 2
---

## Debian / Raspberry Pi

These steps describe how to setup the LoRa Gateway Config utility on a
Debian / Raspbian based gateway. This process has been tested using:

* Debian / Raspbian Jessie

### LoRa Server Debian repository

The LoRa Server project provides pre-compiled binaries packaged as Debian (.deb)
packages. In order to activate this repository, execute the following
commands:

```bash
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 1CE2AFD36DBCCA00

export DISTRIB_ID=`lsb_release -si`
export DISTRIB_CODENAME=`lsb_release -sc`
sudo echo "deb https://repos.loraserver.io/${DISTRIB_ID,,} ${DISTRIB_CODENAME} testing" | sudo tee /etc/apt/sources.list.d/loraserver.list
sudo apt-get update
```

### Install LoRa Gateway Config

In order to instal LoRa Gateway Config, execute the following command:

```bash
sudo apt-get install lora-gateway-config
```

After installation, modify the configuration file which is located at
`/etc/default/lora-gateway-config`.

Settings you probably want to set / change:

* `GW_MAC`
* `GW_SERVER`
* `GW_CLIENT_JWT_TOKEN`
* `BASE_CONFIG_FILE`
* `OUTPUT_CONFIG_FILE`

Please refer to [configuration]({{< relref "config.md" >}}) for more information
about these configuration variables.

### Starting LoRa Gateway Config

To (re)start LoRa Gateway Config, use the following commands:

```bash
sudo systemctl [start|stop|restart|status] lora-gateway-config
```

### LoRa Gateway Config log output

Now you've setup LoRa Gateway Config, it is a good time to verify that
LoRa Gateway Config is actually up-and-running. This can be done by
looking at the LoRa Gateway Config log output.

```bash
sudo journalctl -u lora-gateway-config -f -n 50
```

Example output:

```
level=info msg="starting LoRa Gateway Config" base_config_file="/opt/semtech/packet_forwarder/lora_pkt_fwd/global_conf.json" docs="https://docs.loraserver.io/" output_config_file="/opt/semtech/packet_forwarder/lora_pkt_fwd/local_conf.json" version=0.1.0
level=info msg="connecting to gateway-server" ca-cert= server="localhost:8002" tls-cert= tls-key=
level=info msg="checking for updated configuration"
level=info msg="configuration written to disk" path="/opt/semtech/packet_forwarder/lora_pkt_fwd/local_conf.json"
level=info msg="invoking packet-forwarder restart command" args=[restart packet-forwarder] cmd=systemctl
level=info msg="packet-forwarder restart command invoked" output=
level=info msg="sleeping until next update check" duration=10s
level=info msg="checking for updated configuration"
level=info msg="no configuration update available"
level=info msg="sleeping until next update check" duration=10s
```
