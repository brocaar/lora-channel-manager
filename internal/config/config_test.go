package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/brocaar/loraserver/api/gw"
	"github.com/brocaar/lorawan"
)

func init() {
	log.SetLevel(log.ErrorLevel)
}

type testGatewayClient struct {
	GetConfigurationRequestChan chan gw.GetConfigurationRequest
	GetConfigurationResponse    gw.GetConfigurationResponse
}

func (c *testGatewayClient) GetConfiguration(ctx context.Context, req *gw.GetConfigurationRequest, opts ...grpc.CallOption) (*gw.GetConfigurationResponse, error) {
	c.GetConfigurationRequestChan <- *req
	return &c.GetConfigurationResponse, nil
}

func TestGetPacketForwarderConfig(t *testing.T) {
	Convey("Given a mocked GatewayClient", t, func() {
		client := testGatewayClient{
			GetConfigurationRequestChan: make(chan gw.GetConfigurationRequest, 100),
		}
		GatewayClient = &client
		GatewayMAC = lorawan.EUI64{1, 2, 3, 4, 5, 6, 7, 8}
		now := time.Now().UTC()

		testTable := []struct {
			Name                            string
			GetConfigurationResponse        gw.GetConfigurationResponse
			ExpectedGetConfigurationRequest gw.GetConfigurationRequest
			ExpectedGatewayConfig           gatewayConfiguration
			ExpectedError                   error
		}{
			{
				Name: "EU 868 band config (minimal configuration)",
				GetConfigurationResponse: gw.GetConfigurationResponse{
					UpdatedAt: now.Format(time.RFC3339Nano),
					Channels: []*gw.Channel{
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868100000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868300000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868500000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
					},
				},
				ExpectedGetConfigurationRequest: gw.GetConfigurationRequest{
					Mac: GatewayMAC[:],
				},
				ExpectedGatewayConfig: gatewayConfiguration{
					UpdatedAt: now,
					Radios: [radioCount]radioConfig{
						{
							Enable: true,
							Freq:   868500000,
						},
					},
					MultiSFChannels: [channelCount]multiSFChannelConfig{
						{
							Enable: true,
							Radio:  0,
							IF:     -400000,
							Freq:   868100000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     -200000,
							Freq:   868300000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     0,
							Freq:   868500000,
						},
					},
				},
			},
			{
				Name: "EU 868 band config + CFList + LoRa single-SF + FSK",
				GetConfigurationResponse: gw.GetConfigurationResponse{
					UpdatedAt: now.Format(time.RFC3339Nano),
					Channels: []*gw.Channel{
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868100000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868300000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868500000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     867100000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     867300000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     867500000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     867700000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     867900000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     868300000,
							Bandwidth:     250,
							SpreadFactors: []int32{7},
						},
						{
							Modulation: gw.Modulation_FSK,
							Frequency:  868800000,
							Bandwidth:  125,
							BitRate:    50000,
						},
					},
				},
				ExpectedGetConfigurationRequest: gw.GetConfigurationRequest{
					Mac: GatewayMAC[:],
				},
				ExpectedGatewayConfig: gatewayConfiguration{
					UpdatedAt: now,
					Radios: [radioCount]radioConfig{
						{
							Enable: true,
							Freq:   867500000,
						},
						{
							Enable: true,
							Freq:   868500000,
						},
					},
					MultiSFChannels: [channelCount]multiSFChannelConfig{
						{
							Enable: true,
							Radio:  1,
							IF:     -400000,
							Freq:   868100000,
						},
						{
							Enable: true,
							Radio:  1,
							IF:     -200000,
							Freq:   868300000,
						},
						{
							Enable: true,
							Radio:  1,
							IF:     0,
							Freq:   868500000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     -400000,
							Freq:   867100000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     -200000,
							Freq:   867300000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     0,
							Freq:   867500000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     200000,
							Freq:   867700000,
						},
						{
							Enable: true,
							Radio:  0,
							IF:     400000,
							Freq:   867900000,
						},
					},
					LoRaSTDChannelConfig: loRaSTDChannelConfig{
						Enable:       true,
						Radio:        1,
						IF:           -200000,
						Bandwidth:    250000,
						SpreadFactor: 7,
						Freq:         868300000,
					},
					FSKChannelConfig: fskChannelConfig{
						Enable:    true,
						Radio:     1,
						IF:        300000,
						Bandwidth: 125,
						DataRate:  50000,
						Freq:      868800000,
					},
				},
			},
			{
				Name: "US band (0-7 + 64)",
				GetConfigurationResponse: gw.GetConfigurationResponse{
					UpdatedAt: now.Format(time.RFC3339Nano),
					Channels: []*gw.Channel{
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     902300000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     902500000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     902700000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     902900000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     903100000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     903300000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     903500000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     903700000,
							Bandwidth:     125,
							SpreadFactors: []int32{7, 8, 9, 10},
						},
						{
							Modulation:    gw.Modulation_LORA,
							Frequency:     903000000,
							Bandwidth:     500,
							SpreadFactors: []int32{8},
						},
					},
				},
				ExpectedGetConfigurationRequest: gw.GetConfigurationRequest{
					Mac: GatewayMAC[:],
				},
				ExpectedGatewayConfig: gatewayConfiguration{
					UpdatedAt: now,
					Radios: [radioCount]radioConfig{
						{
							Enable: true,
							Freq:   902700000,
						},
						{
							Enable: true,
							Freq:   903700000,
						},
					},
					MultiSFChannels: [channelCount]multiSFChannelConfig{
						{
							Enable: true,
							Freq:   902300000,
							Radio:  0,
							IF:     -400000,
						},
						{

							Enable: true,
							Freq:   902500000,
							Radio:  0,
							IF:     -200000,
						},
						{

							Enable: true,
							Freq:   902700000,
							Radio:  0,
							IF:     0,
						},
						{

							Enable: true,
							Freq:   902900000,
							Radio:  0,
							IF:     200000,
						},
						{

							Enable: true,
							Freq:   903100000,
							Radio:  0,
							IF:     400000,
						},
						{

							Enable: true,
							Freq:   903300000,
							Radio:  1,
							IF:     -400000,
						},
						{

							Enable: true,
							Freq:   903500000,
							Radio:  1,
							IF:     -200000,
						},
						{

							Enable: true,
							Freq:   903700000,
							Radio:  1,
							IF:     0,
						},
					},
					LoRaSTDChannelConfig: loRaSTDChannelConfig{
						Enable:       true,
						Freq:         903000000,
						Radio:        0,
						IF:           300000,
						Bandwidth:    500000,
						SpreadFactor: 8,
					},
				},
			},
		}

		for i, test := range testTable {
			Convey(fmt.Sprintf("Testing: %s [%d]", test.Name, i), func() {
				client.GetConfigurationResponse = test.GetConfigurationResponse

				So(client.GetConfigurationRequestChan, ShouldHaveLength, 0)

				pfConfig, err := getGatewayConfig()
				So(err, ShouldResemble, test.ExpectedError)

				So(client.GetConfigurationRequestChan, ShouldHaveLength, 1)
				So(<-client.GetConfigurationRequestChan, ShouldResemble, test.ExpectedGetConfigurationRequest)

				if test.ExpectedError == nil {
					So(pfConfig.Radios, ShouldResemble, test.ExpectedGatewayConfig.Radios)
					So(pfConfig, ShouldResemble, test.ExpectedGatewayConfig)
				}
			})
		}
	})
}

func TestUpdateConfig(t *testing.T) {
	Convey("Given a mocked GatewayClient", t, func() {
		now := time.Now()

		// create temp dir
		tempDir, err := ioutil.TempDir("", "test")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempDir)

		client := testGatewayClient{
			GetConfigurationRequestChan: make(chan gw.GetConfigurationRequest, 100),
		}
		client.GetConfigurationResponse = gw.GetConfigurationResponse{
			UpdatedAt: now.Format(time.RFC3339Nano),
			Channels: []*gw.Channel{
				{
					Modulation:    gw.Modulation_LORA,
					Frequency:     868100000,
					Bandwidth:     125,
					SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
				},
				{
					Modulation:    gw.Modulation_LORA,
					Frequency:     868300000,
					Bandwidth:     125,
					SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
				},
				{
					Modulation:    gw.Modulation_LORA,
					Frequency:     868500000,
					Bandwidth:     125,
					SpreadFactors: []int32{7, 8, 9, 10, 11, 12},
				},
			},
		}

		GatewayClient = &client
		GatewayMAC = lorawan.EUI64{1, 2, 3, 4, 5, 6, 7, 8}
		PFRestartCommand = fmt.Sprintf("touch %s", filepath.Join(tempDir, "restart"))
		BaseConfigFile = filepath.Join("test/test.json")
		OutputConfigFile = filepath.Join(tempDir, "out.json")

		Convey("When calling updateConfig", func() {
			err := updateConfig()
			So(err, ShouldBeNil)

			Convey("Then the new configuration has been written", func() {
				Convey("Then the new configuration contains the expected values", func() {
					conf, err := loadConfigFile(OutputConfigFile)
					So(err, ShouldBeNil)

					gwConfig, err := getGatewayConfig()
					So(err, ShouldBeNil)

					// test radios
					for i, r := range gwConfig.Radios {
						radio := conf.SX1301Conf[fmt.Sprintf("radio_%d", i)].(map[string]interface{})
						expected := map[string]interface{}{
							"enable": r.Enable,
							"freq":   r.Freq,
						}
						for k, v := range expected {
							So(radio[k], ShouldEqual, v)
						}
					}

					// test multi SF channels
					for i, c := range gwConfig.MultiSFChannels {
						channel := conf.SX1301Conf[fmt.Sprintf("chan_multiSF_%d", i)].(map[string]interface{})
						expected := map[string]interface{}{
							"enable": c.Enable,
							"radio":  c.Radio,
							"if":     c.IF,
						}
						for k, v := range expected {
							So(channel[k], ShouldEqual, v)
						}
					}

					// test LoRa std channel
					channel := conf.SX1301Conf["chan_Lora_std"].(map[string]interface{})
					expected := map[string]interface{}{
						"enable":        gwConfig.LoRaSTDChannelConfig.Enable,
						"radio":         gwConfig.LoRaSTDChannelConfig.Radio,
						"if":            gwConfig.LoRaSTDChannelConfig.IF,
						"bandwidth":     gwConfig.LoRaSTDChannelConfig.Bandwidth,
						"spread_factor": gwConfig.LoRaSTDChannelConfig.SpreadFactor,
					}
					for k, v := range expected {
						So(channel[k], ShouldEqual, v)
					}

					// test FSK channel
					channel = conf.SX1301Conf["chan_FSK"].(map[string]interface{})
					expected = map[string]interface{}{
						"enable":    gwConfig.FSKChannelConfig.Enable,
						"radio":     gwConfig.FSKChannelConfig.Radio,
						"if":        gwConfig.FSKChannelConfig.IF,
						"bandwidth": gwConfig.FSKChannelConfig.Bandwidth,
						"datarate":  gwConfig.FSKChannelConfig.DataRate,
					}
					for k, v := range expected {
						So(channel[k], ShouldEqual, v)
					}

					// test gateway mac / gateway_ID
					So(conf.GatewayConf["gateway_ID"], ShouldEqual, GatewayMAC.String())
				})
			})

			Convey("Then the restart packet-forwarder command has been invoked", func() {
				_, err := os.Stat(filepath.Join(tempDir, "restart"))
				So(err, ShouldBeNil)
			})
		})
	})
}
