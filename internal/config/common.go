package config

import (
	"time"

	"github.com/brocaar/loraserver/api/gw"
	"github.com/brocaar/lorawan"
)

// GatewayMAC contains the MAC of the gateway.
var GatewayMAC lorawan.EUI64

// GatewayClient is the client interfacing with the gateway-server.
var GatewayClient gw.GatewayClient

// ConfigPollInterval contains the interval between polling new configuration.
var ConfigPollInterval time.Duration

// PFRestartCommand contains the command to restart the packet-forwarder.
var PFRestartCommand string

// BaseConfigFile contains the path to the base config file.
var BaseConfigFile string

// OutputConfigFile contains the path to the output config file.
var OutputConfigFile string

// radioBandwidthPerChannelBandwidth defines the bandwidth that a single radio
// can cover per channel bandwidth
var radioBandwidthPerChannelBandwidth = map[int]int{
	500000: 1100000, // 500kHz channel
	250000: 1000000, // 250kHz channel
	125000: 925000,  // 125kHz channel
}

// defaultRadioBandwidth defines the radio bandwidth in case the channel
// bandwidth does not match any of the above values.
const defaultRadioBandwidth = 925000

// radioCount defines the number of radios available
const radioCount = 2

// channelCount defines the number of available channels
const channelCount = 8
