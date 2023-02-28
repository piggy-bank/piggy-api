package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL string
	DB      *DBConfig     `yaml:"data_base"`
	Sender  *SenderConfig `yaml:"sender"`
}

type SenderConfig struct {
	CallbackURL string `yaml:"callback_url"`
}
type DBConfig struct {
	Dialect      string `yaml:"dialect"`
	Username     string `yaml:"user_name"`
	Password     string `yaml:"password"`
	Project      string `yaml:"project"`
	Zone         string `yaml:"zone"`
	InstanceName string `yaml:"instance_name"`
	DatabaseName string `yaml:"database_name"`
}

func GetConfig(path string) *Config {
	propertiesFile, err := os.Open(path)
	if err != nil {
		msg := "error getting properties file for path'" + path + "', error: " + err.Error()
		log.Print("[package:configuration][method:GetConfig]" + msg)
		panic(msg)
	}

	defer propertiesFile.Close()

	decoder := yaml.NewDecoder(propertiesFile)
	var cfg Config

	err = decoder.Decode(&cfg)
	if err != nil {
		msg := "error decoding yaml, file properties and Config struct does not match, error: " + err.Error()
		log.Printf("[package:configuration][method:GetConfig]" + msg)
		panic(msg)
	}
	return &cfg
}
