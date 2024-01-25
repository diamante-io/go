package config

import (
	"go/types"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfigOption_UsageText(t *testing.T) {
	configOpt := ConfigOption{
		Usage:  "Port to listen and serve on",
		EnvVar: "PORT",
	}
	assert.Equal(t, "Port to listen and serve on (PORT)", configOpt.UsageText())
}

type testOptions struct {
	String string
	Int    int
	Bool   bool
	Uint   uint
	Uint32 uint32
}

// Test that optional flags are set to nil if they are not configured explicitly.
func TestConfigOption_optionalFlags_defaults(t *testing.T) {
	var optUint *uint
	var optString *string
	configOpts := ConfigOptions{
		{Name: "uint", OptType: types.Uint, ConfigKey: &optUint, CustomSetValue: SetOptionalUint, FlagDefault: uint(0)},
		{Name: "string", OptType: types.String, ConfigKey: &optString, CustomSetValue: SetOptionalString},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	cmd.SetArgs([]string{})
	cmd.Execute()
	assert.Equal(t, (*string)(nil), optString)
	assert.Equal(t, (*uint)(nil), optUint)
}

// Test that optional flags are set to non nil values when they are configured explicitly.
func TestConfigOption_optionalFlags_set(t *testing.T) {
	var optUint *uint
	var optString *string
	configOpts := ConfigOptions{
		{Name: "uint", OptType: types.Uint, ConfigKey: &optUint, CustomSetValue: SetOptionalUint, FlagDefault: uint(0)},
		{Name: "string", OptType: types.String, ConfigKey: &optString, CustomSetValue: SetOptionalString},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	cmd.SetArgs([]string{"--uint", "6", "--string", "test-string"})
	cmd.Execute()
	assert.Equal(t, "test-string", *optString)
	assert.Equal(t, uint(6), *optUint)
}

// Test that optional flags are set to non nil values when they are configured explicitly.
func TestConfigOption_optionalFlags_env_set_empty(t *testing.T) {
	var optUint *uint
	var optString *string
	configOpts := ConfigOptions{
		{Name: "uint", OptType: types.Uint, ConfigKey: &optUint, CustomSetValue: SetOptionalUint, FlagDefault: uint(0)},
		{Name: "string", OptType: types.String, ConfigKey: &optString, CustomSetValue: SetOptionalString},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	prev := envVars
	envVars = map[string]bool{
		"STRING": true,
	}
	defer func() {
		envVars = prev
	}()

	cmd.Execute()
	assert.Equal(t, "", *optString)
	assert.Equal(t, (*uint)(nil), optUint)
}

// Test that optional flags are set to non nil values when they are configured explicitly.
func TestConfigOption_optionalFlags_env_set(t *testing.T) {
	var optUint *uint
	var optString *string
	configOpts := ConfigOptions{
		{Name: "uint", OptType: types.Uint, ConfigKey: &optUint, CustomSetValue: SetOptionalUint, FlagDefault: uint(0)},
		{Name: "string", OptType: types.String, ConfigKey: &optString, CustomSetValue: SetOptionalString},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	prev := envVars
	envVars = map[string]bool{
		"STRING": true,
		"UINT":   true,
	}
	defer func() {
		envVars = prev
	}()

	defer os.Setenv("STRING", os.Getenv("STRING"))
	defer os.Setenv("UINT", os.Getenv("UINT"))
	os.Setenv("STRING", "str")
	os.Setenv("UINT", "6")

	cmd.Execute()
	assert.Equal(t, "str", *optString)
	assert.Equal(t, uint(6), *optUint)
}

// Test that when there are no args the defaults in the config options are
// used.
func TestConfigOption_getSimpleValue_defaults(t *testing.T) {
	opts := testOptions{}
	configOpts := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts.String, FlagDefault: "default"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts.Int, FlagDefault: 1},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts.Bool, FlagDefault: true},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts.Uint, FlagDefault: uint(2)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts.Uint32, FlagDefault: uint32(3)},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	cmd.SetArgs([]string{})
	cmd.Execute()
	assert.Equal(t, "default", opts.String)
	assert.Equal(t, 1, opts.Int)
	assert.Equal(t, true, opts.Bool)
	assert.Equal(t, uint(2), opts.Uint)
	assert.Equal(t, uint32(3), opts.Uint32)
	for _, opt := range configOpts {
		assert.False(t, opt.flag.Changed)
	}
}

// Test that when args are given, their values are used.
func TestConfigOption_getSimpleValue_setFlag(t *testing.T) {
	opts := testOptions{}
	configOpts := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts.String, FlagDefault: "default"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts.Int, FlagDefault: 1},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts.Bool, FlagDefault: false},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts.Uint, FlagDefault: uint(2)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts.Uint32, FlagDefault: uint32(3)},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	cmd.SetArgs([]string{
		"--string", "value",
		"--int", "10",
		"--bool",
		"--uint", "20",
		"--uint32", "30",
	})
	cmd.Execute()
	assert.Equal(t, "value", opts.String)
	assert.Equal(t, 10, opts.Int)
	assert.Equal(t, true, opts.Bool)
	assert.Equal(t, uint(20), opts.Uint)
	assert.Equal(t, uint32(30), opts.Uint32)
	for _, opt := range configOpts {
		assert.True(t, opt.flag.Changed)
	}
}

// Test that when args are not given but env vars are, their values are used.
func TestConfigOption_getSimpleValue_setEnv(t *testing.T) {
	opts := testOptions{}
	configOpts := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts.String, FlagDefault: "default"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts.Int, FlagDefault: 1},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts.Bool, FlagDefault: false},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts.Uint, FlagDefault: uint(2)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts.Uint32, FlagDefault: uint32(3)},
	}
	cmd := &cobra.Command{
		Use: "doathing",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
		},
	}
	configOpts.Init(cmd)

	defer os.Setenv("STRING", os.Getenv("STRING"))
	defer os.Setenv("INT", os.Getenv("INT"))
	defer os.Setenv("BOOL", os.Getenv("BOOL"))
	defer os.Setenv("UINT", os.Getenv("UINT"))
	defer os.Setenv("UINT32", os.Getenv("UINT32"))
	os.Setenv("STRING", "value")
	os.Setenv("INT", "10")
	os.Setenv("BOOL", "true")
	os.Setenv("UINT", "20")
	os.Setenv("UINT32", "30")
	cmd.Execute()
	assert.Equal(t, "value", opts.String)
	assert.Equal(t, 10, opts.Int)
	assert.Equal(t, true, opts.Bool)
	assert.Equal(t, uint(20), opts.Uint)
	assert.Equal(t, uint32(30), opts.Uint32)
}

// Test that when multiple commands register the same option, they can be set
// with flags.
func TestConfigOption_getSimpleValue_setMultipleFlag(t *testing.T) {
	opts1 := testOptions{}
	configOpts1 := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts1.String, FlagDefault: "default1"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts1.Int, FlagDefault: 11},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts1.Bool, FlagDefault: false},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts1.Uint, FlagDefault: uint(12)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts1.Uint32, FlagDefault: uint32(13)},
	}
	cmd1 := &cobra.Command{
		Use: "doathing1",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts1.Require()
			configOpts1.SetValues()
		},
	}
	configOpts1.Init(cmd1)

	opts2 := testOptions{}
	configOpts2 := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts2.String, FlagDefault: "default2"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts2.Int, FlagDefault: 21},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts2.Bool, FlagDefault: false},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts2.Uint, FlagDefault: uint(22)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts2.Uint32, FlagDefault: uint32(23)},
	}
	cmd2 := &cobra.Command{
		Use: "doathing2",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts2.Require()
			configOpts2.SetValues()
		},
	}
	configOpts2.Init(cmd2)

	cmd1.SetArgs([]string{
		"--string", "value1",
		"--int", "110",
		"--bool",
		"--uint", "120",
		"--uint32", "130",
	})
	cmd1.Execute()
	assert.Equal(t, "value1", opts1.String)
	assert.Equal(t, 110, opts1.Int)
	assert.Equal(t, true, opts1.Bool)
	assert.Equal(t, uint(120), opts1.Uint)
	assert.Equal(t, uint32(130), opts1.Uint32)

	cmd2.SetArgs([]string{
		"--string", "value2",
		"--int", "210",
		"--bool",
		"--uint", "220",
		"--uint32", "230",
	})
	cmd2.Execute()
	assert.Equal(t, "value2", opts2.String)
	assert.Equal(t, 210, opts2.Int)
	assert.Equal(t, true, opts2.Bool)
	assert.Equal(t, uint(220), opts2.Uint)
	assert.Equal(t, uint32(230), opts2.Uint32)
}

// Test that when multiple commands register the same option, they can be set
// with environment variables.
func TestConfigOption_getSimpleValue_setMultipleEnv(t *testing.T) {
	opts1 := testOptions{}
	configOpts1 := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts1.String, FlagDefault: "default1"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts1.Int, FlagDefault: 11},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts1.Bool, FlagDefault: false},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts1.Uint, FlagDefault: uint(12)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts1.Uint32, FlagDefault: uint32(13)},
	}
	cmd1 := &cobra.Command{
		Use: "doathing1",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts1.Require()
			configOpts1.SetValues()
		},
	}
	configOpts1.Init(cmd1)

	opts2 := testOptions{}
	configOpts2 := ConfigOptions{
		{Name: "string", OptType: types.String, ConfigKey: &opts2.String, FlagDefault: "default2"},
		{Name: "int", OptType: types.Int, ConfigKey: &opts2.Int, FlagDefault: 21},
		{Name: "bool", OptType: types.Bool, ConfigKey: &opts2.Bool, FlagDefault: false},
		{Name: "uint", OptType: types.Uint, ConfigKey: &opts2.Uint, FlagDefault: uint(22)},
		{Name: "uint32", OptType: types.Uint32, ConfigKey: &opts2.Uint32, FlagDefault: uint32(23)},
	}
	cmd2 := &cobra.Command{
		Use: "doathing2",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts2.Require()
			configOpts2.SetValues()
		},
	}
	configOpts2.Init(cmd2)

	defer os.Setenv("STRING", os.Getenv("STRING"))
	defer os.Setenv("INT", os.Getenv("INT"))
	defer os.Setenv("BOOL", os.Getenv("BOOL"))
	defer os.Setenv("UINT", os.Getenv("UINT"))
	defer os.Setenv("UINT32", os.Getenv("UINT32"))

	os.Setenv("STRING", "value1")
	os.Setenv("INT", "110")
	os.Setenv("BOOL", "true")
	os.Setenv("UINT", "120")
	os.Setenv("UINT32", "130")

	cmd1.Execute()
	assert.Equal(t, "value1", opts1.String)
	assert.Equal(t, 110, opts1.Int)
	assert.Equal(t, true, opts1.Bool)
	assert.Equal(t, uint(120), opts1.Uint)
	assert.Equal(t, uint32(130), opts1.Uint32)

	cmd2.Execute()
	assert.Equal(t, "value1", opts2.String)
	assert.Equal(t, 110, opts2.Int)
	assert.Equal(t, true, opts2.Bool)
	assert.Equal(t, uint(120), opts2.Uint)
	assert.Equal(t, uint32(130), opts2.Uint32)
}
