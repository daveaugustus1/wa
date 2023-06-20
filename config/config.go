package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/Expand-My-Business/go_windows_agent/constants"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

var lock = &sync.Mutex{}

type Config struct {
	CompanyCode         string `toml:"company_code"`
	InstructionEndpoint string `toml:"instruction_endpoint"`
	Port                string `toml:"port"`
}

var SingleConfigInstance *Config
var err error

func GetConfigInstance() *Config {
	if SingleConfigInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if SingleConfigInstance == nil {
			fmt.Println("Creating single instance now.")
			SingleConfigInstance, err = loadConfigFromFile(constants.ConfigFilePath)
			if err != nil {
				logrus.Fatalf("Error occured while loading config file, error: %+v", err)
			}
		} else {
			logrus.Info("Single instance already created.")
		}
	} else {
		logrus.Info("Single instance already created.")
	}

	return SingleConfigInstance
}

func loadConfigFromFile(filePath string) (*Config, error) {
	// Read the TOML file content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Create a Config struct instance
	config := Config{}

	// Unmarshal the TOML data into the Config struct
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
