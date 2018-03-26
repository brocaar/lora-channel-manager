---
title: Kerlink gateways
menu:
  main:
    parent: install
    weight: 3
---

# Kerlink gateways

## Kerlink IOT station

The Kerlink IOT station has a mechanism to start "custom" application on boot.
These steps will install the LoRa Channel Manager ARM build on the Kerlink.

1. Create the directories needed:

    `mkdir -p /mnt/fsuser-1/lora-channel-manager`

2. Download and extract the LoRa Channel Manager ARM binary into the above
   directory. See [downloads]({{< ref "overview/downloads.md" >}}).
   Make sure the binary is marked as executable.

3. Save the following content as `/mnt/fsuser-1/lora-channel-manager/lora-channel-manager.sh`:

    ```bash
    #!/bin/sh

    LOGGER="logger -p local1.notice"

    # firewall rule for loraserver API
    iptables -A INPUT -p tcp --sport 8002 -j ACCEPT
    
    cd /mnt/fsuser-1/lora-channel-manager/

    ./lora-channel-manager --gw-mac YOUR_GW_MAC --gw-server YOUR_GW_SERVER:8002 --pf-restart-command "killall -15 execute_spf.sh" --base-config-file /mnt/fsuser-1/lora-channel-manager/global_conf_EU868.json --output-config-file /mnt/fsuser-1/spf/etc/global_conf.json --gw-client-jwt-token YOUR_GW_TOKEN 2>&1 | $LOGGER -t lora-channel-manager
    ```

    Make sure to replace:
    - `YOUR_GW_MAC` with your gateway MAC address in HEX format (e.g. 0A0B0C0D...)
    - `YOUR_GW_SERVER` with the hostname/IP of your loraserver
    - `YOUR_GW_TOKEN` with token, generated to this gateway on your App Server Web interface (Organization - Gateways - Select this gateway - Gateway token, click Generate and copy generated token)
    
    Also make sure the file is marked as executable.
    
    Place `global_conf_EU868.json` to `/mnt/fsuser-1/lora-channel-manager/` directory

4. Save the following content as `/mnt/fsuser-1/lora-gateway-bridge/manifest.xml`:

    ```xml
    <?xml version="1.0"?>
    <manifest>
    <app name="lora-channel-manager" binary="lora-channel-manager.sh" >
    <start param="" autostart="y"/>
    <stop kill="9"/>
    </app>
    </manifest>
    ```
 5. Because `lora-channel-manager` make new JSON configuration file as one line, we need patch packet forwarder starting script `/mnt/fsuser-1/spf/bin/execute_spf.sh`:
    - find function jsonval()
    - insert line `awk -v k="text" '{n=split($0,a,","); for (i=1; i<=n; i++) print a[i]}' |\` between two line started from `sed...`

    Your result must be:
    ```
    jsonval() {
        sed -e 's!\([^\\]\)"!\1!g' -e 's!,$!!' $1 | \
        awk -v k="text" '{n=split($0,a,","); for (i=1; i<=n; i++) print a[i]}' |\
        sed -n "s!.*${2}: *\(.*\)!\1!p"
    }
    ```
 6. Reboot gateway or run `/etc/init.d/knet restart`
 7. See logfile for details `tail -f /mnt/fsuser-1/spf/var/log/spf.log`
