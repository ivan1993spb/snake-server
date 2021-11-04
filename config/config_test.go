package config

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_DefaultConfig_ReturnsDefaultConfig(t *testing.T) {
	require.Equal(t, defaultConfig, DefaultConfig())
}

func Test_ParseFlags_ParsesFlagsCorrectly(t *testing.T) {
	const flagSetName = "test"

	type Test struct {
		msg string

		args     []string
		defaults Config

		expectConfig Config
		expectErr    bool
	}

	var tests = make([]*Test, 0)

	// Test case 1
	tests = append(tests, &Test{
		msg: "run without arguments",

		args:     []string{},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	// Test case 2
	configTest2 := defaultConfig
	configTest2.Server.Address = ":7070"

	tests = append(tests, &Test{
		msg: "change address",

		args:     []string{"-address", ":7070"},
		defaults: defaultConfig,

		expectConfig: configTest2,
		expectErr:    false,
	})

	// Test case 3
	configTest3 := defaultConfig
	configTest3.Server.Address = "localhost:6670"
	configTest3.Server.Seed = 0

	tests = append(tests, &Test{
		msg: "change address and seed",

		args:     []string{"-address", "localhost:6670", "-seed", "0"},
		defaults: defaultConfig,

		expectConfig: configTest3,
		expectErr:    false,
	})

	// Test case 4
	configTest4 := defaultConfig
	configTest4.Server.Address = "snakeonline.xyz:7986"
	configTest4.Server.Seed = 0
	configTest4.Server.Log.EnableJSON = true

	tests = append(tests, &Test{
		msg: "change address, seed and logging",

		args: []string{
			"-address", "snakeonline.xyz:7986",
			"-seed", "0",
			"-log-json",
		},
		defaults: defaultConfig,

		expectConfig: configTest4,
		expectErr:    false,
	})

	// Test case 5
	configTest5 := defaultConfig
	configTest5.Server.Address = "snakeonline.xyz:3211"
	configTest5.Server.Seed = 32
	configTest5.Server.Log.EnableJSON = true
	configTest5.Server.Flags.EnableWeb = true
	configTest5.Server.Limits.Conns = 321

	tests = append(tests, &Test{
		msg: "change address, seed, logging, connection limit and enable web",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-seed", "32",
			"-log-json",
			"-enable-web",
			"-conns-limit", "321",
		},
		defaults: defaultConfig,

		expectConfig: configTest5,
		expectErr:    false,
	})

	// Test case 6
	tests = append(tests, &Test{
		msg: "change address, connection limit, enable web and make 1 mistake",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-enable-web",
			"-conns-limit", "321",
			"-foobar",
		},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 7
	tests = append(tests, &Test{
		msg: "change address, connection limit, enable web and make 2 mistakes",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-enable-web",
			"-groups-limit", "error",
			"-conns-limit", "321",
			"-foobar",
		},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 8
	tests = append(tests, &Test{
		msg: "args is nil",

		args:     nil,
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	// Test case 9
	configTest9 := defaultConfig
	configTest9.Server.Address = "snakeonline.xyz:3211"
	configTest9.Server.Sentry.Enable = true
	configTest9.Server.Sentry.DSN = "https://public@sentry.example.com/44"

	tests = append(tests, &Test{
		msg: "change address, sentry settings",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-sentry-enable",
			"-sentry-dsn", "https://public@sentry.example.com/44",
		},
		defaults: defaultConfig,

		expectConfig: configTest9,
		expectErr:    false,
	})

	// Test case 10
	configTest10 := defaultConfig
	configTest10.Server.Address = "snakeonline.xyz:3211"
	configTest10.Server.Flags.Debug = true

	tests = append(tests, &Test{
		msg: "change address, debug trace",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-debug",
		},
		defaults: defaultConfig,

		expectConfig: configTest10,
		expectErr:    false,
	})

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)
		flagSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
		flagSet.SetOutput(ioutil.Discard)

		config, err := ParseFlags(flagSet, test.args, test.defaults)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)
	}
}

func Test_ParseYAML_ParsesYAMLCorrectly(t *testing.T) {
	type Test struct {
		msg string

		input    []byte
		defaults Config

		expectConfig Config
		expectErr    bool
	}

	var tests = make([]*Test, 0)

	// Test case 1
	tests = append(tests, &Test{
		msg: "input is nil",

		input:    nil,
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	// Test case 2
	configTest2 := defaultConfig
	configTest2.Server.Seed = 0

	tests = append(tests, &Test{
		msg: "a valid config with seed equal 0",

		input:    ConfigYAMLSampleDefault,
		defaults: defaultConfig,

		expectConfig: configTest2,
		expectErr:    false,
	})

	// Test case 3
	configTest3 := defaultConfig
	configTest3.Server.Address = ":9999"
	configTest3.Server.TLS.Enable = true
	configTest3.Server.TLS.Cert = "path/to/cert"
	configTest3.Server.TLS.Key = "path/to/key"

	tests = append(tests, &Test{
		msg: "a valid config with address and TLS settings",

		input:    ConfigYAMLSampleAddressAndTLS,
		defaults: defaultConfig,

		expectConfig: configTest3,
		expectErr:    false,
	})

	// Test case 4
	tests = append(tests, &Test{
		msg: "bullshit YAML syntax of the config",

		input:    ConfigYAMLSampleBullshitSyntax,
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 5
	tests = append(tests, &Test{
		msg: "empty config",

		input:    []byte{},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	// Test case 6
	configTest6 := defaultConfig
	configTest6.Server.Limits.Groups = 144
	configTest6.Server.Limits.Conns = 4123
	configTest6.Server.Sentry.Enable = true
	configTest6.Server.Sentry.DSN = "https://public@sentry.example.com/1"

	tests = append(tests, &Test{
		msg: "limits and sentry",

		input:    ConfigYAMLSampleLimitsAndSentry,
		defaults: defaultConfig,

		expectConfig: configTest6,
		expectErr:    false,
	})

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)

		config, err := ParseYAML(test.input, test.defaults)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)
	}
}

func Test_Config_Fields_ReturnsFieldsOfTheConfig(t *testing.T) {
	require.Equal(t, map[string]interface{}{
		fieldLabelAddress: ":9999",

		fieldLabelTLSEnable: true,
		fieldLabelTLSCert:   "path/to/cert",
		fieldLabelTLSKey:    "path/to/key",

		fieldLabelGroupsLimit: 1000,
		fieldLabelConnsLimit:  10000,

		fieldLabelSeed: int64(321),

		fieldLabelLogEnableJSON: false,
		fieldLabelLogLevel:      "warning",

		fieldLabelFlagsEnableBroadcast: true,
		fieldLabelFlagsEnableWeb:       false,
		fieldLabelFlagsForbidCORS:      true,
		fieldLabelFlagsDebug:           true,

		fieldLabelSentryEnable: true,
		fieldLabelSentryDSN:    "https://public@sentry.example.com/1",
	}, Config{
		Server: Server{
			Address: ":9999",

			TLS: TLS{
				Enable: true,
				Cert:   "path/to/cert",
				Key:    "path/to/key",
			},

			Limits: Limits{
				Groups: 1000,
				Conns:  10000,
			},

			Seed: 321,

			Log: Log{
				EnableJSON: false,
				Level:      "warning",
			},

			Flags: Flags{
				EnableBroadcast: true,
				EnableWeb:       false,
				ForbidCORS:      true,
				Debug:           true,
			},

			Sentry: Sentry{
				Enable: true,
				DSN:    "https://public@sentry.example.com/1",
			},
		},
	}.Fields())
}

func Test_ReadYAMLConfig_ReadsConfigCorrectly(t *testing.T) {
	type Test struct {
		msg string

		input    []byte
		defaults Config

		expectConfig Config
		expectErr    bool
	}

	var tests = make([]*Test, 0)

	// Test case 1
	configTest1 := defaultConfig
	configTest1.Server.Seed = 0

	tests = append(tests, &Test{
		msg: "a valid config with seed equal 0",

		input:    ConfigYAMLSampleDefault,
		defaults: defaultConfig,

		expectConfig: configTest1,
		expectErr:    false,
	})

	// Test case 2
	configTest2 := defaultConfig
	configTest2.Server.Address = ":9999"
	configTest2.Server.TLS.Enable = true
	configTest2.Server.TLS.Cert = "path/to/cert"
	configTest2.Server.TLS.Key = "path/to/key"

	tests = append(tests, &Test{
		msg: "a valid config with address and TLS settings",

		input:    ConfigYAMLSampleAddressAndTLS,
		defaults: defaultConfig,

		expectConfig: configTest2,
		expectErr:    false,
	})

	// Test case 3
	tests = append(tests, &Test{
		msg: "bullshit YAML syntax of the config",

		input:    ConfigYAMLSampleBullshitSyntax,
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 4
	tests = append(tests, &Test{
		msg: "empty reader",

		input:    make([]byte, 0),
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	// Test case 5
	configTest5 := defaultConfig
	configTest5.Server.Address = ":9999"
	configTest5.Server.TLS.Enable = true
	configTest5.Server.TLS.Cert = "path/to/cert"
	configTest5.Server.TLS.Key = "path/to/key"
	configTest5.Server.Limits.Groups = 144
	configTest5.Server.Limits.Conns = 4123
	configTest5.Server.Flags.EnableBroadcast = true
	configTest5.Server.Flags.ForbidCORS = true

	tests = append(tests, &Test{
		msg: "a valid config with address, TLS settings, limits, and flags of broadcast and CORS",

		input:    ConfigYAMLSampleAddressAndTLSAndLimitsAndCORS,
		defaults: defaultConfig,

		expectConfig: configTest5,
		expectErr:    false,
	})

	// Test case 6
	configTest6 := defaultConfig
	configTest6.Server.Sentry.Enable = true
	configTest6.Server.Sentry.DSN = "https://public@sentry.example.com/1"
	configTest6.Server.Flags.Debug = true

	tests = append(tests, &Test{
		msg: "a valid config with sentry and debug settings",

		input:    ConfigYAMLSampleSentryAndDebug,
		defaults: defaultConfig,

		expectConfig: configTest6,
		expectErr:    false,
	})

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)

		config, err := ReadYAMLConfig(bytes.NewBuffer(test.input), test.defaults)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)
	}
}

func Test_Configurate_ReturnsCorrectConfig(t *testing.T) {
	const (
		configPath = "/etc/server/test/config.yaml"

		perm = 0666

		flagSetName = "test"
	)

	type Test struct {
		msg string

		input []byte
		args  []string

		expectConfig Config
		expectErr    bool

		setEnv     bool
		saveConfig bool
	}

	var tests = make([]*Test, 0)

	// Test case 1
	configTest1 := defaultConfig
	configTest1.Server.Seed = 0

	tests = append(tests, &Test{
		msg: "a valid config with seed equal 0",

		input: ConfigYAMLSampleDefault,
		args:  []string{},

		expectConfig: configTest1,
		expectErr:    false,

		setEnv:     true,
		saveConfig: true,
	})

	// Test case 2
	tests = append(tests, &Test{
		msg: "environment variable hasn't been set and flags are empty",

		input: ConfigYAMLSampleDefault,
		args:  []string{},

		expectConfig: defaultConfig,
		expectErr:    false,

		setEnv:     false,
		saveConfig: true,
	})

	// Test case 3
	configTest3 := defaultConfig
	configTest3.Server.Log.EnableJSON = true
	configTest3.Server.Limits.Groups = 120

	tests = append(tests, &Test{
		msg: "environment variable hasn't been set, json logging is enabled, group limits is 120",

		input: ConfigYAMLSampleDefault,
		args:  []string{"-log-json", "-groups-limit", "120"},

		expectConfig: configTest3,
		expectErr:    false,

		setEnv:     false,
		saveConfig: true,
	})

	// Test case 4
	tests = append(tests, &Test{
		msg: "environment variable has been set, config is empty, flags are incorrect",

		input: make([]byte, 0),
		args:  []string{"-log-json", "-invalid-flag"},

		expectConfig: defaultConfig,
		expectErr:    true,

		setEnv:     true,
		saveConfig: true,
	})

	// Test case 5
	tests = append(tests, &Test{
		msg: "environment variable has been set, config not found",

		input: make([]byte, 0),
		args:  []string{"-log-json"},

		expectConfig: defaultConfig,
		expectErr:    true,

		setEnv:     true,
		saveConfig: false,
	})

	// Test case 5
	tests = append(tests, &Test{
		msg: "environment variable has been set, config exists, YAML syntax is invalid",

		input: ConfigYAMLSampleBullshitSyntax,
		args:  []string{"-log-json", "-enable-web"},

		expectConfig: defaultConfig,
		expectErr:    true,

		setEnv:     true,
		saveConfig: true,
	})

	// Test case 6
	configTest6 := defaultConfig
	configTest6.Server.Address = ":9999"
	configTest6.Server.TLS.Enable = true
	configTest6.Server.TLS.Cert = "/etc/path/cert"
	configTest6.Server.TLS.Key = "path/to/key"
	configTest6.Server.Limits.Groups = 422
	configTest6.Server.Limits.Conns = 4123
	configTest6.Server.Flags.EnableBroadcast = false
	configTest6.Server.Log.EnableJSON = true
	configTest6.Server.Flags.EnableWeb = true

	tests = append(tests, &Test{
		msg: "overwrite configuration in file config and with flags",

		input: ConfigYAMLSampleAddressAndTLSAndLimits,
		args: []string{
			"-log-json",
			"-enable-web",
			"-groups-limit", "422",
			"-tls-cert", "/etc/path/cert",
			"-enable-broadcast=false",
		},

		expectConfig: configTest6,
		expectErr:    false,

		setEnv:     true,
		saveConfig: true,
	})

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)

		fs := afero.NewMemMapFs()

		if test.saveConfig {
			err := afero.WriteFile(fs, configPath, test.input, perm)
			require.Nil(t, err, label)
		}

		if test.setEnv {
			err := os.Setenv(envVarSnakeServerConfigPath, configPath)
			require.Nil(t, err)
		}

		flagSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
		flagSet.SetOutput(ioutil.Discard)

		config, err := Configurate(fs, flagSet, test.args)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)

		if test.setEnv {
			err := os.Unsetenv(envVarSnakeServerConfigPath)
			require.Nil(t, err)
		}
	}
}
