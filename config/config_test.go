package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

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

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n)
		fs := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
		fs.SetOutput(ioutil.Discard)

		config, err := ParseFlags(fs, test.args, test.defaults)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)
	}
}
