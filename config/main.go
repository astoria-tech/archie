package config

import (
	"io/ioutil"

	"github.com/astoria-arc/archie/msgs"
	yaml "gopkg.in/yaml.v2"
)

//Config configuration for slack bot
type Config struct {
	Messages msgs.Messages `yaml:"messages"`
}

//Load config
func Load(path string) (*Config, error) {

	// Load file
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Load config
	config := &Config{}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
