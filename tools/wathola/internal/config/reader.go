package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"knative.dev/test-infra/tools/wathola/internal/ensure"
	"os"
)

var location = "~/.config/wathola/config.toml"
var logFatal = log.Fatal

// ReadIfPresent read a configuration file if it exists
func ReadIfPresent() {
	configFile, err := homedir.Expand(location)
	ensure.NoError(err)
	if fileExists(configFile) {
		log.Infof("Reading config file: %v", configFile)
		err := Read(configFile)
		if err != nil {
			logFatal(err)
		}
	} else {
		log.Infof("Define config file to be taken into account: %v", configFile)
	}
}

// Read a config file and update configuration object
func Read(configFile string) error {
	r, err := os.Open(configFile)
	if err != nil {
		return err
	}
	d := toml.NewDecoder(r)
	err = d.Decode(Instance)
	if err == nil {
		log.SetLevel(Instance.LogLevel)
	}
	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
