package config

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
	configTest5.Server.EnableWeb = true
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
		msg: "r is nil",

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

		fieldLabelEnableBroadcast: true,
		fieldLabelEnableWeb:       false,
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

			EnableBroadcast: true,
			EnableWeb:       false,
		},
	}.Fields())
}

func Test_ReadYAMLConfig_ReadsConfigCorrectly(t *testing.T) {
	type Test struct {
		msg string

		r        io.Reader
		defaults Config

		expectConfig Config
		expectErr    bool
	}

	var tests = make([]*Test, 0)

	// Test case 1
	configTest2 := defaultConfig
	configTest2.Server.Seed = 0

	tests = append(tests, &Test{
		msg: "a valid config with seed equal 0",

		r:        bytes.NewBuffer(ConfigYAMLSampleDefault),
		defaults: defaultConfig,

		expectConfig: configTest2,
		expectErr:    false,
	})

	// Test case 2
	configTest3 := defaultConfig
	configTest3.Server.Address = ":9999"
	configTest3.Server.TLS.Enable = true
	configTest3.Server.TLS.Cert = "path/to/cert"
	configTest3.Server.TLS.Key = "path/to/key"

	tests = append(tests, &Test{
		msg: "a valid config with address and TLS settings",

		r:        bytes.NewBuffer(ConfigYAMLSampleAddressAndTLS),
		defaults: defaultConfig,

		expectConfig: configTest3,
		expectErr:    false,
	})

	// Test case 3
	tests = append(tests, &Test{
		msg: "bullshit YAML syntax of the config",

		r:        bytes.NewBuffer(ConfigYAMLSampleBullshitSyntax),
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 4
	tests = append(tests, &Test{
		msg: "empty reader",

		r:        bytes.NewBuffer(make([]byte, 0)),
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)

		config, err := ReadYAMLConfig(test.r, test.defaults)

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
	})

	err := os.Setenv(envVarSnakeServerConfigPath, configPath)
	require.Nil(t, err)

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, configPath, test.input, perm)
		require.Nil(t, err, label)

		flagSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
		flagSet.SetOutput(ioutil.Discard)

		config, err := Configurate(fs, flagSet, test.args)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)
	}

	os.Unsetenv(envVarSnakeServerConfigPath)
}
