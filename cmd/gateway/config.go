package main

import (
	"strings"

	"github.com/mniak/ps121/pkg/dynpb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadConfig() Config {
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	cobra.CheckErr(viper.ReadInConfig())
	cobra.CheckErr(viper.Unmarshal(&config))
	return config
}

type Config struct {
	Endpoints []ConfigEndpoint `mapstructure:"endpoints"`
}

type ConfigEndpoint struct {
	Method string        `mapstructure:"method"`
	Fields []ConfigField `mapstructure:"fields"`
}

type ConfigField struct {
	Index int    `mapstructure:"index"`
	Alias string `mapstructure:"alias"`
	Type  string `mapstructure:"type"`
}

func (cfg Config) AliasesPerMethod() map[string]AliasTree {
	result := make(map[string]AliasTree)
	for _, endpoint := range cfg.Endpoints {

		tree := make(AliasTree)
		result[endpoint.Method] = tree

		for _, field := range endpoint.Fields {
			tree[field.Index] = AliasNode{
				Alias: field.Alias,
			}
		}

	}

	return result
}

func (cfg Config) HintsPerMethod() map[string]dynpb.HintMap {
	result := make(map[string]dynpb.HintMap)
	for _, endpoint := range cfg.Endpoints {

		hints := make(dynpb.HintMap)
		result[endpoint.Method] = hints

		for _, field := range endpoint.Fields {
			switch strings.ToLower(field.Type) {
			case "string":
				hints[field.Index] = dynpb.HintString
			// case "int":
			// 	hints[field.Index] = dynpb.HintInt32
			case "struct":
				hints[field.Index] = dynpb.HintObject(nil)
			}
		}
	}

	return result
}
