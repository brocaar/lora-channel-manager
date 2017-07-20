package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/brocaar/loraserver/api/gw"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var lastUpdatedAt time.Time

var jsonCommentRegexp = regexp.MustCompile(`/\*.*\*/`)

type radioConfig struct {
	Enable bool
	Freq   int
}

type multiSFChannelConfig struct {
	Enable bool
	Radio  int
	IF     int
	Freq   int
}

type loRaSTDChannelConfig struct {
	Enable       bool
	Radio        int
	IF           int
	Bandwidth    int
	SpreadFactor int
	Freq         int
}

type fskChannelConfig struct {
	Enable    bool
	Radio     int
	IF        int
	Bandwidth int
	DataRate  int
	Freq      int
}

type configFile struct {
	SX1301Conf  map[string]interface{} `json:"SX1301_conf"`
	GatewayConf map[string]interface{} `json:"gateway_conf"`
}

type gatewayConfiguration struct {
	UpdatedAt            time.Time
	Radios               [radioCount]radioConfig
	MultiSFChannels      [channelCount]multiSFChannelConfig
	LoRaSTDChannelConfig loRaSTDChannelConfig
	FSKChannelConfig     fskChannelConfig
}

// channelByMinRadioCenterFreqency implements sort.Interface for []*gw.Channel.
// The sorting is based on the center frequency of the radio when placing the
// channel exactly on the left side of the available radio bandwidth.
type channelByMinRadioCenterFrequency []*gw.Channel

func (c channelByMinRadioCenterFrequency) Len() int      { return len(c) }
func (c channelByMinRadioCenterFrequency) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c channelByMinRadioCenterFrequency) Less(i, j int) bool {
	return c.minRadioCenterFreq(i) < c.minRadioCenterFreq(j)
}
func (c channelByMinRadioCenterFrequency) minRadioCenterFreq(i int) int {
	channelBandwidth := int(c[i].Bandwidth * 1000)
	radioBandwidth, ok := radioBandwidthPerChannelBandwidth[channelBandwidth]
	if !ok {
		radioBandwidth = defaultRadioBandwidth
	}
	return int(c[i].Frequency) - (channelBandwidth / 2) + (radioBandwidth / 2)
}

// UpdateConfigLoop checks for new configuration, writes new configuration
// to disk and invokes the packet-forwarder restart command.
func UpdateConfigLoop() {
	for {
		log.Info("checking for updated configuration")
		if err := updateConfig(); err != nil {
			log.Errorf("update config error: %s", err)
		}
		log.WithField("duration", ConfigPollInterval).Info("sleeping until next update check")
		time.Sleep(ConfigPollInterval)
	}
}

// updateConfig fetches the latest configuration from the gateway-server api,
// loads the base configuration file, injects the new configuration and writes
// this to disk.
func updateConfig() error {
	// get latest config
	conf, err := getGatewayConfig()
	if err != nil {
		return errors.Wrap(err, "get packet-forwarder config error")
	}

	if lastUpdatedAt.Equal(conf.UpdatedAt) {
		log.Info("no configuration update available")
		return nil
	}

	// load base config
	baseConf, err := loadConfigFile(BaseConfigFile)
	if err != nil {
		return errors.Wrap(err, "load config file error")
	}

	// merge the config
	if err = mergeConfig(baseConf, conf); err != nil {
		return errors.Wrap(err, "merge config error")
	}

	// generate config json
	b, err := json.Marshal(baseConf)
	if err != nil {
		return err
	}

	// write file to disk
	if err = ioutil.WriteFile(OutputConfigFile, b, 0644); err != nil {
		return errors.Wrap(err, "write file error")
	}
	log.WithField("path", OutputConfigFile).Info("configuration written to disk")

	// invoke restart command
	if err = invokePFRestart(); err != nil {
		return errors.Wrap(err, "invoke packet-forwarder restart error")
	}

	// set last updated timestamp
	lastUpdatedAt = conf.UpdatedAt

	return nil
}

func invokePFRestart() error {
	parts := strings.Fields(PFRestartCommand)
	if len(parts) == 0 {
		return errors.New("no packet-forwarder restart command configured")
	}

	var args []string
	if len(parts) > 1 {
		args = parts[1:len(parts)]
	}

	log.WithFields(log.Fields{
		"cmd":  parts[0],
		"args": args,
	}).Info("invoking packet-forwarder restart command")

	out, err := exec.Command(parts[0], args...).Output()
	if err != nil {
		return errors.Wrap(err, "execute command error")
	}
	log.WithField("output", string(out)).Info("packet-forwarder restart command invoked")

	return nil
}

func loadConfigFile(filePath string) (configFile, error) {
	var out configFile

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return out, errors.Wrap(err, "read file error")
	}

	// remove comments from json
	b = jsonCommentRegexp.ReplaceAll(b, []byte{})

	if err = json.Unmarshal(b, &out); err != nil {
		return out, errors.Wrap(err, "unmarshal config json error")
	}

	return out, nil
}

// mergeConfig merges the new configuration into the given configuration.
// Unfortunately we have to do this as the packet-forwarder sees these keys
// as complete overrides (it does not just update the leaves).
// We want to remain the other configuration (e.g. which radio chip is used,
// calibration values that are board specific).
// This is not pretty but it works.
func mergeConfig(config configFile, newConfig gatewayConfiguration) error {
	// update radios
	for i, r := range newConfig.Radios {
		radio, ok := config.SX1301Conf[fmt.Sprintf("radio_%d", i)].(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected radio_%d to be of type map[string]interface{}, got %T", i, config.SX1301Conf[fmt.Sprintf("radio_%d", i)])
		}
		radio["enable"] = r.Enable
		radio["freq"] = r.Freq
	}

	// update multi SF channels
	for i, c := range newConfig.MultiSFChannels {
		channel, ok := config.SX1301Conf[fmt.Sprintf("chan_multiSF_%d", i)].(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected chan_multiSF_%d to be of type map[string]interface{}, got %T", i, config.SX1301Conf[fmt.Sprintf("chan_multiSF_%d", i)])
		}
		channel["enable"] = c.Enable
		channel["radio"] = c.Radio
		channel["if"] = c.IF
	}

	// update LoRa std channel
	channel, ok := config.SX1301Conf["chan_Lora_std"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected chan_Lora_std to be of type map[string]interface{}, got %T", config.SX1301Conf["chan_Lora_std"])
	}
	channel["enable"] = newConfig.LoRaSTDChannelConfig.Enable
	channel["radio"] = newConfig.LoRaSTDChannelConfig.Radio
	channel["if"] = newConfig.LoRaSTDChannelConfig.IF
	channel["bandwidth"] = newConfig.LoRaSTDChannelConfig.Bandwidth
	channel["spread_factor"] = newConfig.LoRaSTDChannelConfig.SpreadFactor

	// update FSK channel
	channel, ok = config.SX1301Conf["chan_FSK"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected chan_FSK to be of type map[string]interface{}, got %T", config.SX1301Conf["chan_FSK"])
	}
	channel["enable"] = newConfig.FSKChannelConfig.Enable
	channel["radio"] = newConfig.FSKChannelConfig.Radio
	channel["if"] = newConfig.FSKChannelConfig.IF
	channel["bandwidth"] = newConfig.FSKChannelConfig.Bandwidth
	channel["datarate"] = newConfig.FSKChannelConfig.DataRate

	// update gateway mac / ID
	config.GatewayConf["gateway_ID"] = GatewayMAC.String()

	return nil
}

func getGatewayConfig() (gatewayConfiguration, error) {
	var conf gatewayConfiguration
	var multiSFCounter int

	configResp, err := GatewayClient.GetConfiguration(context.Background(), &gw.GetConfigurationRequest{
		Mac: GatewayMAC[:],
	})
	if err != nil {
		return conf, errors.Wrap(err, "get configuration error")
	}

	// set UpdatedAt
	ts, err := time.Parse(time.RFC3339Nano, configResp.UpdatedAt)
	if err != nil {
		return conf, errors.Wrap(err, "parse time error")
	}
	conf.UpdatedAt = ts

	// make sure the channels are sorted by the minimum radio center frequency
	channelsCopy := make([]*gw.Channel, len(configResp.Channels))
	copy(channelsCopy, configResp.Channels)
	sort.Sort(channelByMinRadioCenterFrequency(channelsCopy))

	// define the radios and their center frequency
	for _, c := range channelsCopy {
		channelBandwidth := int(c.Bandwidth * 1000)
		channelMax := int(c.Frequency) + (channelBandwidth / 2)
		radioBandwidth, ok := radioBandwidthPerChannelBandwidth[channelBandwidth]
		if !ok {
			radioBandwidth = defaultRadioBandwidth
		}
		minRadioCenterFreq := int(c.Frequency) - (channelBandwidth / 2) + (radioBandwidth / 2)

		for i, r := range conf.Radios {
			// the radio is not defined yet, use it
			if !r.Enable {
				conf.Radios[i].Enable = true
				conf.Radios[i].Freq = minRadioCenterFreq

				break
			}

			if channelMax <= r.Freq+(radioBandwidth/2) {
				break
			}
		}
	}

	// assign channels
	for _, c := range configResp.Channels {
		var radio int

		channelBandwidth := int(c.Bandwidth * 1000)
		channelMin := int(c.Frequency) - (channelBandwidth / 2)
		channelMax := int(c.Frequency) + (channelBandwidth / 2)
		radioBandwidth, ok := radioBandwidthPerChannelBandwidth[channelBandwidth]
		if !ok {
			radioBandwidth = defaultRadioBandwidth
		}

		// get the radio covering the channel frequency
		for i, r := range conf.Radios {
			if channelMin >= r.Freq-(radioBandwidth/2) && channelMax <= r.Freq+(radioBandwidth/2) {
				radio = i
				break
			}
		}

		if c.Modulation == gw.Modulation_FSK {
			// FSK channel
			if conf.FSKChannelConfig.Enable {
				return conf, errors.New("FSK channel already configured")
			}

			conf.FSKChannelConfig = fskChannelConfig{
				Enable:    true,
				Radio:     radio,
				IF:        int(c.Frequency) - conf.Radios[radio].Freq,
				Bandwidth: int(c.Bandwidth),
				DataRate:  int(c.BitRate),
				Freq:      int(c.Frequency),
			}

		} else if c.Modulation == gw.Modulation_LORA && len(c.SpreadFactors) == 1 {
			// LoRa STD (single SF) channel
			if conf.LoRaSTDChannelConfig.Enable {
				return conf, errors.New("LoRa std channel already configured")
			}

			conf.LoRaSTDChannelConfig = loRaSTDChannelConfig{
				Enable:       true,
				Radio:        radio,
				IF:           int(c.Frequency) - conf.Radios[radio].Freq,
				Bandwidth:    channelBandwidth,
				SpreadFactor: int(c.SpreadFactors[0]),
				Freq:         int(c.Frequency),
			}

		} else if c.Modulation == gw.Modulation_LORA {
			// LoRa multi-SF channels
			if multiSFCounter > channelCount {
				return conf, errors.New("exceeded maximum number of multi-SF channels")
			}

			conf.MultiSFChannels[multiSFCounter] = multiSFChannelConfig{
				Enable: true,
				Radio:  radio,
				IF:     int(c.Frequency) - conf.Radios[radio].Freq,
				Freq:   int(c.Frequency),
			}

			multiSFCounter++

		} else {
			return conf, fmt.Errorf("invalid modulation %s", c.Modulation)
		}
	}

	return conf, nil
}
