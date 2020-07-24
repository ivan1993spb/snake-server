package config

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	defaultAddress = ":8080"

	defaultTLSEnable = false
	defaultTLSCert   = ""
	defaultTLSKey    = ""

	defaultGroupsLimit = 100
	defaultConnsLimit  = 1000

	defaultLogEnableJSON = false
	defaultLogLevel      = "info"

	defaultEnableBroadcast = false
	defaultEnableWeb       = false
)

const (
	flagLabelAddress = "address"

	flagLabelTLSEnable = "tls-enable"
	flagLabelTLSCert   = "tls-cert"
	flagLabelTLSKey    = "tls-key"

	flagLabelGroupsLimit = "groups-limit"
	flagLabelConnsLimit  = "conns-limit"

	flagLabelSeed = "seed"

	flagLabelLogEnableJSON = "log-json"
	flagLabelLogLevel      = "log-level"

	flagLabelEnableBroadcast = "enable-broadcast"
	flagLabelEnableWeb       = "enable-web"
)

const (
	flagUsageAddress = "address to serve"

	flagUsageTLSEnable = "enable TLS"
	flagUsageTLSCert   = "path to certificate file"
	flagUsageTLSKey    = "path to key file"

	flagUsageGroupsLimit = "game groups limit"
	flagUsageConnsLimit  = "web-socket connections limit"

	flagUsageSeed = "random seed"

	flagUsageLogEnableJSON = "use json format for logger"
	flagUsageLogLevel      = "set log level: panic, fatal, error, warning (warn), info or debug"

	flagUsageEnableBroadcast = "enable broadcasting API method"
	flagUsageEnableWeb       = "enable web client"
)

const (
	fieldLabelAddress = "address"

	fieldLabelTLSEnable = "tls-enable"
	fieldLabelTLSCert   = "tls-cert"
	fieldLabelTLSKey    = "tls-key"

	fieldLabelGroupsLimit = "groups-limit"
	fieldLabelConnsLimit  = "conns-limit"

	fieldLabelSeed = "seed"

	fieldLabelLogEnableJSON = "log-json"
	fieldLabelLogLevel      = "log-level"

	fieldLabelEnableBroadcast = "enable-broadcast"
	fieldLabelEnableWeb       = "enable-web"
)

const envVarSnakeServerConfigPath = "SNAKE_SERVER_CONFIG_PATH"

func generateSeed() int64 {
	return time.Now().UnixNano()
}

type TLS struct {
	Enable bool   `yaml:"enable"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
}

type Limits struct {
	Groups int `yaml:"groups"`
	Conns  int `yaml:"conns"`
}

type Log struct {
	EnableJSON bool   `yaml:"enable_json"`
	Level      string `yaml:"level"`
}

type Server struct {
	Address string `yaml:"address"`

	TLS    TLS    `yaml:"tls"`
	Limits Limits `yaml:"limits"`
	Seed   int64  `yaml:"seed"`
	Log    Log    `yaml:"log"`

	EnableBroadcast bool `yaml:"enable_broadcast"`
	EnableWeb       bool `yaml:"enable_web"`
}

// Config is a server configuration structure
type Config struct {
	Server Server `yaml:"server"`
}

func (c Config) Fields() map[string]interface{} {
	return map[string]interface{}{
		fieldLabelAddress: c.Server.Address,

		fieldLabelTLSEnable: c.Server.TLS.Enable,
		fieldLabelTLSCert:   c.Server.TLS.Cert,
		fieldLabelTLSKey:    c.Server.TLS.Key,

		fieldLabelGroupsLimit: c.Server.Limits.Groups,
		fieldLabelConnsLimit:  c.Server.Limits.Conns,

		fieldLabelSeed: c.Server.Seed,

		fieldLabelLogEnableJSON: c.Server.Log.EnableJSON,
		fieldLabelLogLevel:      c.Server.Log.Level,

		fieldLabelEnableBroadcast: c.Server.EnableBroadcast,
		fieldLabelEnableWeb:       c.Server.EnableWeb,
	}
}

var defaultConfig = Config{
	Server: Server{
		Address: defaultAddress,

		TLS: TLS{
			Enable: defaultTLSEnable,
			Cert:   defaultTLSCert,
			Key:    defaultTLSKey,
		},

		Limits: Limits{
			Groups: defaultGroupsLimit,
			Conns:  defaultConnsLimit,
		},

		Seed: generateSeed(),

		Log: Log{
			EnableJSON: defaultLogEnableJSON,
			Level:      defaultLogLevel,
		},

		EnableBroadcast: defaultEnableBroadcast,
		EnableWeb:       defaultEnableWeb,
	},
}

func DefaultConfig() Config {
	return defaultConfig
}

func ParseYAML(input []byte, defaults Config) (Config, error) {
	config := defaults

	if err := yaml.Unmarshal(input, &config); err != nil {
		return defaults, fmt.Errorf("cannot parse YAML: %s", err)
	}

	return config, nil
}

func ParseFlags(fs *flag.FlagSet, args []string, defaults Config) (Config, error) {
	if fs.Parsed() {
		panic("program composition error: the provided FlagSet has been parsed")
	}

	config := defaults

	// Address
	fs.StringVar(&config.Server.Address, flagLabelAddress, defaults.Server.Address, flagUsageAddress)

	// TLS
	fs.BoolVar(&config.Server.TLS.Enable, flagLabelTLSEnable, defaults.Server.TLS.Enable, flagUsageTLSEnable)
	fs.StringVar(&config.Server.TLS.Cert, flagLabelTLSCert, defaults.Server.TLS.Cert, flagUsageTLSCert)
	fs.StringVar(&config.Server.TLS.Key, flagLabelTLSKey, defaults.Server.TLS.Key, flagUsageTLSKey)

	// Limits
	fs.IntVar(&config.Server.Limits.Groups, flagLabelGroupsLimit, defaults.Server.Limits.Groups, flagUsageGroupsLimit)
	fs.IntVar(&config.Server.Limits.Conns, flagLabelConnsLimit, defaults.Server.Limits.Conns, flagUsageConnsLimit)

	// Random
	fs.Int64Var(&config.Server.Seed, flagLabelSeed, defaults.Server.Seed, flagUsageSeed)

	// Logging
	fs.BoolVar(&config.Server.Log.EnableJSON, flagLabelLogEnableJSON, defaults.Server.Log.EnableJSON, flagUsageLogEnableJSON)
	fs.StringVar(&config.Server.Log.Level, flagLabelLogLevel, defaults.Server.Log.Level, flagUsageLogLevel)

	// Flags
	fs.BoolVar(&config.Server.EnableBroadcast, flagLabelEnableBroadcast, defaults.Server.EnableBroadcast, flagUsageEnableBroadcast)
	fs.BoolVar(&config.Server.EnableWeb, flagLabelEnableWeb, defaults.Server.EnableWeb, flagUsageEnableWeb)

	if err := fs.Parse(args); err != nil {
		return defaults, fmt.Errorf("cannot parse flags: %s", err)
	}

	return config, nil
}

type errReadConfigYAML struct {
	err error
}

func (e *errReadConfigYAML) Error() string {
	return fmt.Sprintf("cannot read YAML config: %s", e.err)
}

func ReadYAMLConfig(r io.Reader, defaults Config) (Config, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return defaults, &errReadConfigYAML{err}
	}

	config, err := ParseYAML(input, defaults)
	if err != nil {
		return defaults, &errReadConfigYAML{err}
	}

	return config, nil
}
