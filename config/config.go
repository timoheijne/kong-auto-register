package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Route struct {
	Name        string   `yaml:"name"`
	Paths       []string `yaml:"paths"`
	Methods     []string `yaml:"methods"`
	Description string   `yaml:"description"`
}

type Service struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Path string `yaml:"path"`
}

type KARConfig struct {
	Service       Service `yaml:"service"`
	AdminEndpoint string  `yaml:"adminEndpoint"`
	Debug         bool    `yaml:"debug"`
	Routes        []Route `yaml:"routes"`
}

var config *KARConfig

func SetConfig(new_config KARConfig) {
	config = &new_config

	// Environment variables overrule programatically set config variables
	// TODO: Rewrite this to make use of another package to handle reading and setting of env vars.
	val, ok := os.LookupEnv("KAR_Debug")
	if ok {
		bo, err := strconv.ParseBool(val)
		check(err)

		config.Debug = bo
	}

	val, ok = os.LookupEnv("KAR_AdminEndpoint")
	if ok {
		config.AdminEndpoint = val
	}

	if config.Debug == false {
		log.SetOutput(ioutil.Discard)
	}
}

func GetConfig() *KARConfig {
	return config
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
