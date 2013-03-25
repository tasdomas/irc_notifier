/*
 Parse IRC notifier configuration.
*/
package config

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type ChannelConfig struct {
	Name  string
	Watch string
	Nick  string
}

type Config struct {
	BotNick  string
	Network  string
	Port     int
	SSL      bool
	Password string
	Channels []ChannelConfig
}

func LoadConfig(filename string) (*Config, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load config file %s: %v",
			filename, err)
	}
	var c Config
	err = goyaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config file: %s: %v",
			filename, err)
	}
	if c.Port == 0 {
		c.Port = 6697
	}
	return &c, nil
}
