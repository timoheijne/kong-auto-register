package kar

import (
	"io/ioutil"

	config "gitlab.liquidstudio.nl/deeplens/kong-auto-registration/config"
	"gitlab.liquidstudio.nl/deeplens/kong-auto-registration/kong"
	"gopkg.in/yaml.v2"
)

func kar() {
	// 1. Will read config from YAML file OR programatic input
	// 2. Config variables might be overriden by environment variables
	// 3. Register routes using the various methods
	// 4. Run kong sync
}

// Init will simply set the config you pass along
func Init(new_config config.KARConfig) {
	config.SetConfig(new_config)

	// Run KONG check
	kong.Run()
}

// InitFromFile will load the yaml file provided and parse the config from there
func InitFromYaml(file string) {
	conf := config.KARConfig{}

	dat, err := ioutil.ReadFile(file)
	check(err)

	err = yaml.Unmarshal([]byte(dat), &conf)
	check(err)

	Init(conf)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
