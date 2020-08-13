package config

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Default values for the server's settings
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
	defaultForbidCORS      = false
)

// Flag labels
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
	flagLabelForbidCORS      = "forbid-cors"
)

// Flag usage descriptions
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
	flagUsageForbidCORS      = "forbid cross-origin resource sharing"
)

// Label names
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
	fieldLabelForbidCORS      = "forbid-cors"
)

const envVarSnakeServerConfigPath = "SNAKE_SERVER_CONFIG_PATH"

func generateSeed() int64 {
	return time.Now().UnixNano()
}

// TLS structure represents TLS config
type TLS struct {
	Enable bool   `yaml:"enable"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
}

// Limits structure sets up server limits
type Limits struct {
	Groups int `yaml:"groups"`
	Conns  int `yaml:"conns"`
}

// Log structure defines preferences for logging
type Log struct {
	EnableJSON bool   `yaml:"enable_json"`
	Level      string `yaml:"level"`
}

// Server structure contains configurations for the server
type Server struct {
	Address string `yaml:"address"`

	TLS    TLS    `yaml:"tls"`
	Limits Limits `yaml:"limits"`
	Seed   int64  `yaml:"seed"`
	Log    Log    `yaml:"log"`

	EnableBroadcast bool `yaml:"enable_broadcast"`
	EnableWeb       bool `yaml:"enable_web"`
	ForbidCORS      bool `yaml:"forbid_cors"`
}

// Config is a base server configuration structure
type Config struct {
	Server Server `yaml:"server"`
}

// Fields returns a map of all configurations
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
		fieldLabelForbidCORS:      c.Server.ForbidCORS,
	}
}

// Default settings
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
		ForbidCORS:      defaultForbidCORS,
	},
}

// DefaultConfig returns configuration by default
func DefaultConfig() Config {
	return defaultConfig
}

// ParseYAML parses input byte slice and returns a config based on the default configuration
func ParseYAML(input []byte, defaults Config) (Config, error) {
	config := defaults

	if err := yaml.Unmarshal(input, &config); err != nil {
		return defaults, fmt.Errorf("cannot parse YAML: %s", err)
	}

	return config, nil
}

// ParseFlags parses flags and returns a config based on the default configuration
func ParseFlags(flagSet *flag.FlagSet, args []string, defaults Config) (Config, error) {
	if flagSet.Parsed() {
		panic("program composition error: the provided FlagSet has been parsed")
	}

	config := defaults

	// Address
	flagSet.StringVar(&config.Server.Address, flagLabelAddress, defaults.Server.Address, flagUsageAddress)

	// TLS
	flagSet.BoolVar(&config.Server.TLS.Enable, flagLabelTLSEnable, defaults.Server.TLS.Enable, flagUsageTLSEnable)
	flagSet.StringVar(&config.Server.TLS.Cert, flagLabelTLSCert, defaults.Server.TLS.Cert, flagUsageTLSCert)
	flagSet.StringVar(&config.Server.TLS.Key, flagLabelTLSKey, defaults.Server.TLS.Key, flagUsageTLSKey)

	// Limits
	flagSet.IntVar(&config.Server.Limits.Groups, flagLabelGroupsLimit, defaults.Server.Limits.Groups, flagUsageGroupsLimit)
	flagSet.IntVar(&config.Server.Limits.Conns, flagLabelConnsLimit, defaults.Server.Limits.Conns, flagUsageConnsLimit)

	// Random
	flagSet.Int64Var(&config.Server.Seed, flagLabelSeed, defaults.Server.Seed, flagUsageSeed)

	// Logging
	flagSet.BoolVar(&config.Server.Log.EnableJSON, flagLabelLogEnableJSON, defaults.Server.Log.EnableJSON, flagUsageLogEnableJSON)
	flagSet.StringVar(&config.Server.Log.Level, flagLabelLogLevel, defaults.Server.Log.Level, flagUsageLogLevel)

	// Flags
	flagSet.BoolVar(&config.Server.EnableBroadcast, flagLabelEnableBroadcast, defaults.Server.EnableBroadcast, flagUsageEnableBroadcast)
	flagSet.BoolVar(&config.Server.EnableWeb, flagLabelEnableWeb, defaults.Server.EnableWeb, flagUsageEnableWeb)
	flagSet.BoolVar(&config.Server.ForbidCORS, flagLabelForbidCORS, defaults.Server.ForbidCORS, flagUsageForbidCORS)

	if err := flagSet.Parse(args); err != nil {
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

// ReadYAMLConfig reads configurations from a reader and returns a config structure based on defaults
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

type errConfigurate struct {
	err error
}

func (e *errConfigurate) Error() string {
	return fmt.Sprintf("cannot configurate: %s", e.err)
}

// Configurate gathers a config from a config file and a flag set
func Configurate(fs afero.Fs, flagSet *flag.FlagSet, args []string) (Config, error) {
	defaults := DefaultConfig()
	config := defaults

	if configPath, ok := os.LookupEnv(envVarSnakeServerConfigPath); ok {
		f, err := fs.Open(configPath)
		if err != nil {
			return defaults, &errConfigurate{err}
		}

		config, err = ReadYAMLConfig(f, config)
		if err != nil {
			return defaults, &errConfigurate{err}
		}

		if err := f.Close(); err != nil {
			return defaults, &errConfigurate{err}
		}
	}

	config, err := ParseFlags(flagSet, args, config)
	if err != nil {
		return defaults, &errConfigurate{err}
	}

	return config, nil
}
