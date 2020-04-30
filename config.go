package faraday

import (
	"fmt"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lightningnetwork/lnd/build"
)

const (
	defaultRPCPort        = "10009"
	defaultRPCHostPort    = "localhost:" + defaultRPCPort
	defaultMacaroon       = "admin.macaroon"
	defaultNetwork        = "mainnet"
	defaultMinimumMonitor = time.Hour * 24 * 7 * 4 // four weeks in hours
	defaultDebugLevel     = "info"
	defaultRPCListen      = "localhost:8465"
)

type config struct {
	// RPCServer is host:port that lnd's RPC server is listening on.
	RPCServer string `long:"rpcserver" description:"host:port that LND is listening for RPC connections on"`

	// MacaroonDir is the directory containing macaroons.
	MacaroonDir string `long:"macaroondir" description:"Dir containing macaroons"`

	// MacaroonFile is the file name of the macaroon to use.
	MacaroonFile string `long:"macaroonfile" description:"Macaroon file to use."`

	// TLSCertPath is the path to the tls cert that faraday should use.
	TLSCertPath string `long:"tlscertpath" description:"Path to TLS cert"`

	// TestNet is set to true when running on testnet.
	TestNet bool `long:"testnet" description:"Use the testnet network"`

	// Simnet is set to true when using btcd's simnet.
	Simnet bool `long:"simnet" description:"Use simnet"`

	// Simnet is set to true when using bitcoind's regtest.
	Regtest bool `long:"regtest" description:"Use regtest"`

	// MinimumMonitored is the minimum amount of time that a channel must be monitored for before we consider it for termination.
	MinimumMonitored time.Duration `long:"min_monitored" description:"The minimum amount of time that a channel must be monitored for before recommending termination. Valid time units are {s, m, h}."`

	// network is a string containing the network we're running on.
	network string

	// DebugLevel is a string defining the log level for the service either
	// for all subsystems the same or individual level by subsystem.
	DebugLevel string `long:"debuglevel" description:"Debug level for faraday and its subsystems."`

	// RPCListen is the listen address for the faraday rpc server.
	RPCListen string `long:"rpclisten" description:"Address to listen on for gRPC clients."`

	// RESTListen is the listen address for the faraday REST server.
	RESTListen string `long:"restlisten" description:"Address to listen on for REST clients. If not specified, no REST listener will be started."`

	// CORSOrigin specifies the CORS header that should be set on REST responses. No header is added if the value is empty.
	CORSOrigin string `long:"corsorigin" description:"The value to send in the Access-Control-Allow-Origin header. Header will be omitted if empty."`
}

// DefaultConfig returns all default values for the config struct.
func DefaultConfig() config {
	return config{
		RPCServer:        defaultRPCHostPort,
		network:          defaultNetwork,
		MacaroonFile:     defaultMacaroon,
		MinimumMonitored: defaultMinimumMonitor,
		DebugLevel:       defaultDebugLevel,
		RPCListen:        defaultRPCListen,
	}
}

// LoadConfig starts with a skeleton default config, and reads in user provided
// configuration from the command line. It does not provide a full set of
// defaults or validate user input because validation and sensible default
// setting are performed by the lndclient package.
func LoadConfig() (*config, error) {
	// Start with a default config.
	config := DefaultConfig()

	// Parse command line options to obtain user specified values.
	if _, err := flags.Parse(&config); err != nil {
		return nil, err
	}

	var netCount int
	if config.TestNet {
		config.network = "testnet"
		netCount++
	}
	if config.Regtest {
		config.network = "regtest"
		netCount++
	}
	if config.Simnet {
		config.network = "simnet"
		netCount++
	}

	if netCount > 1 {
		return nil, fmt.Errorf("do not specify more than one network flag")
	}

	if err := build.ParseAndSetDebugLevels(config.DebugLevel, logWriter); err != nil {
		return nil, err
	}

	return &config, nil
}
