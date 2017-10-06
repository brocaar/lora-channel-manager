package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/brocaar/lora-channel-manager/internal/config"
	"github.com/brocaar/loraserver/api/gw"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var version string // set by the compiler

type jwt struct {
	token string
}

func (j jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j jwt) RequireTransportSecurity() bool {
	return false
}

func run(c *cli.Context) error {
	// set config variables
	if err := config.GatewayMAC.UnmarshalText([]byte(c.String("gw-mac"))); err != nil {
		log.Fatalf("invalid gw-mac: %s", err)
	}
	config.BaseConfigFile = c.String("base-config-file")
	config.OutputConfigFile = c.String("output-config-file")
	config.PFRestartCommand = c.String("pf-restart-command")
	config.ConfigPollInterval = c.Duration("config-poll-interval")

	log.WithFields(log.Fields{
		"version":            version,
		"docs":               "https://docs.loraserver.io/",
		"base_config_file":   config.BaseConfigFile,
		"output_config_file": config.OutputConfigFile,
	}).Info("starting LoRa Channel Manager")

	// connect to gateway api server
	log.WithFields(log.Fields{
		"server":   c.String("gw-server"),
		"ca-cert":  c.String("gw-client-ca-cert"),
		"tls-cert": c.String("gw-client-tls-cert"),
		"tls-key":  c.String("gw-client-tls-key"),
	}).Info("connecting to gateway-server")
	gwDialOptions := []grpc.DialOption{
		grpc.WithPerRPCCredentials(jwt{token: c.String("gw-client-jwt-token")}),
	}
	if c.String("gw-client-tls-cert") != "" && c.String("gw-client-tls-key") != "" {
		gwDialOptions = append(gwDialOptions, grpc.WithTransportCredentials(
			mustGetTransportCredentials(c.String("gw-client-tls-cert"), c.String("gw-client-tls-key"), c.String("gw-client-ca-cert"), false),
		))
	} else {
		gwDialOptions = append(gwDialOptions, grpc.WithInsecure())
	}
	gwConn, err := grpc.Dial(c.String("gw-server"), gwDialOptions...)
	if err != nil {
		log.Fatalf("gateway-server dial error: %s", err)
	}
	config.GatewayClient = gw.NewGatewayClient(gwConn)

	// run update config loop
	go config.UpdateConfigLoop()

	// wait for stop signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")

	return nil
}

func mustGetTransportCredentials(tlsCert, tlsKey, caCert string, verifyClientCert bool) credentials.TransportCredentials {
	var caCertPool *x509.CertPool
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		log.WithFields(log.Fields{
			"cert": tlsCert,
			"key":  tlsKey,
		}).Fatalf("load key-pair error: %s", err)
	}

	if caCert != "" {
		rawCaCert, err := ioutil.ReadFile(caCert)
		if err != nil {
			log.WithField("ca", caCert).Fatalf("load ca cert error: %s", err)
		}

		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(rawCaCert)
	}

	if verifyClientCert {
		return credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		})
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})
}

func main() {
	app := cli.NewApp()
	app.Name = "lora-channel-manager"
	app.Usage = "channel-configuration daemon for LoRa gateways"
	app.Version = version
	app.Copyright = "see http://github.com/brocaar/lora-channel-manager for copyright information"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "gw-mac",
			Usage:  "mac address of the gateway",
			EnvVar: "GW_MAC",
		},
		cli.StringFlag{
			Name:   "gw-server",
			Usage:  "hostname:ip of the gateway api server",
			Value:  "127.0.0.1:8002",
			EnvVar: "GW_SERVER",
		},
		cli.StringFlag{
			Name:   "gw-client-ca-cert",
			Usage:  "ca certificate used by the gateway-server client (optional)",
			EnvVar: "GW_CLIENT_CA_CERT",
		},
		cli.StringFlag{
			Name:   "gw-client-tls-cert",
			Usage:  "tls certificate used by the gateway-server client (optional)",
			EnvVar: "GW_CLIENT_TLS_CERT",
		},
		cli.StringFlag{
			Name:   "gw-client-tls-key",
			Usage:  "tls key used by the gateway-server client (optional)",
			EnvVar: "GW_CLIENT_TLS_KEY",
		},
		cli.StringFlag{
			Name:   "gw-client-jwt-token",
			Usage:  "jwt token used by the gateway-server client for authentication (issued by LoRa Server)",
			EnvVar: "GW_CLIENT_JWT_TOKEN",
		},
		cli.StringFlag{
			Name:   "base-config-file",
			Usage:  "path to the base configuration file",
			EnvVar: "BASE_CONFIG_FILE",
		},
		cli.StringFlag{
			Name:   "output-config-file",
			Usage:  "path to the output configuration file",
			EnvVar: "OUTPUT_CONFIG_FILE",
		},
		cli.StringFlag{
			Name:   "pf-restart-command",
			Usage:  "command which must be executed on configuration changes to restart the packet-forwarder",
			EnvVar: "PF_RESTART_COMMAND",
		},
		cli.DurationFlag{
			Name:   "config-poll-interval",
			Usage:  "interval between polling new configuration",
			Value:  time.Minute * 5,
			EnvVar: "CONFIG_POLL_INTERVAL",
		},
	}
	app.Run(os.Args)
}
