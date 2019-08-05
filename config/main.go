package config

import (
	"io/ioutil"

	"github.com/astoria-arc/archie/msgs"
	yaml "gopkg.in/yaml.v2"
)

//SlackConfig configuration for Slack
type SlackConfig struct {
	Token string `yaml:"oauthToken"`
	Debug bool   `yaml:"debug"`
}

//Config configuration for slack bot
type Config struct {
	Messages    msgs.Messages `yaml:"messages"`
	SlackConfig SlackConfig   `yaml:"slack"`
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
