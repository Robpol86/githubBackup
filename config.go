package main

import (
	"reflect"

	"gopkg.in/urfave/cli.v2"
)

// Config defines the application configuration.
type Config struct {
	LogFile *string `cli:"log"`
	Quiet   *bool   `cli:"quiet"`
	Verbose *bool   `cli:"verbose"`
}

// FromCLI is passed to cli.App{} in the Action field. It populates the GlobalConfig.
func (c *Config) FromCLI(ctx *cli.Context) error {
	if s := ctx.String("log"); c.LogFile == nil {
		c.LogFile = &s
	}
	if b := ctx.Bool("quiet"); c.Quiet == nil {
		c.Quiet = &b
	}
	if b := ctx.Bool("verbose"); c.Verbose == nil {
		c.Verbose = &b
	}
	return nil
}

// Finalize resolves remaining nil pointers to empty *values.
func (c *Config) Finalize() {
	emptyBool := false
	emptyString := ""
	structValue := reflect.ValueOf(c).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		if !field.IsNil() {
			continue
		}
		if field.Type() == reflect.TypeOf(&emptyBool) {
			field.Set(reflect.ValueOf(&emptyBool))
		} else if field.Type() == reflect.TypeOf(&emptyString) {
			field.Set(reflect.ValueOf(&emptyString))
		}
	}
}

// GlobalConfig will hold the config values for the entire application during runtime.
var GlobalConfig Config
