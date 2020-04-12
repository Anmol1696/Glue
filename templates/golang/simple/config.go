package main

import (
	"fmt"
	"reflect"

	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
)

// Defualt configiguration setting defined here, overridden by env and command line
// Returns a config struct
func NewDefaultConfig() *Config {
	return &Config{
		Addr: ":8080",
	}
}

// Struct that holds all the configurable info
// Make config as large as possible, default values set in func NewDefaultConfig()
// Config is also available as enviornment variables with a prefix `envPrefix`
type Config struct {
	Addr string `name:"addr" json:"addr" usage:"IP address and port to listen on" env:"ADDRESS"`
	// Sample config that holds url of a backend api
	BackendURL string `name:"backend-url" json:"backend-url" usage:"Backend client url" env:"BACKEND_URL"`
	Verbose    bool   `name:"verbose" json:"verbose" usage:"switch on debug / verbose logging"`
}

func (c *Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("addr", c.Addr)
	enc.AddBool("verbose", c.Verbose)
	return nil
}

func getCommandLineOptions() []cli.Flag {
	defaults := NewDefaultConfig()
	var flags []cli.Flag
	count := reflect.TypeOf(Config{}).NumField()
	for i := 0; i < count; i++ {
		field := reflect.TypeOf(Config{}).Field(i)
		usage, found := field.Tag.Lookup("usage")
		if !found {
			continue
		}
		envName := field.Tag.Get("env")
		if envName != "" {
			envName = envPrefix + envName
		}
		optName := field.Tag.Get("name")

		switch t := field.Type; t.Kind() {
		case reflect.Bool:
			dv := reflect.ValueOf(defaults).Elem().FieldByName(field.Name).Bool()
			msg := fmt.Sprintf("%s (default: %t)", usage, dv)
			flags = append(flags, cli.BoolTFlag{
				Name:   optName,
				Usage:  msg,
				EnvVar: envName,
			})
		case reflect.String:
			defaultValue := reflect.ValueOf(defaults).Elem().FieldByName(field.Name).String()
			flags = append(flags, cli.StringFlag{
				Name:   optName,
				Usage:  usage,
				EnvVar: envName,
				Value:  defaultValue,
			})
		}
	}

	return flags
}

func parseCLIOptions(ctx *cli.Context, config *Config) (err error) {
	// iterate the Config and grab command line options via reflection
	count := reflect.TypeOf(config).Elem().NumField()
	for i := 0; i < count; i++ {
		field := reflect.TypeOf(config).Elem().Field(i)
		name := field.Tag.Get("name")

		if ctx.IsSet(name) {
			switch field.Type.Kind() {
			case reflect.Bool:
				reflect.ValueOf(config).Elem().FieldByName(field.Name).SetBool(ctx.Bool(name))
			case reflect.String:
				reflect.ValueOf(config).Elem().FieldByName(field.Name).SetString(ctx.String(name))
			}
		}
	}
	return nil
}
