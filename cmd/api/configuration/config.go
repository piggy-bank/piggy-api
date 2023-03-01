package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	Development              = "dev"
	Production               = "prod"
	Test                     = "test"
	configurationPackagePath = "configfiles"
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

func BuildConfig(profile string) *Config {
	path := configurationPackagePath + "/properties-" + profile + ".yml"
	return GetConfig(path)

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

const flowConfigPath = "./flow.json"

type FlowServiceAccount struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

type FlowConfig struct {
	Accounts struct {
		Emulator FlowServiceAccount `json:"emulator"`
		Testnet  FlowServiceAccount `json:"testnet"`
		Mainnet  FlowServiceAccount `json:"mainnet"`
	} `json:"accounts"`
	Contracts map[string]string `json:"contracts"`
	Networks  struct {
		Emulator string `json:"emulator"`
		Testnet  string `json:"testnet"`
		Mainnet  string `json:"mainnet"`
	} `json:"networks"`
}

func ReadFlowConfig() *FlowConfig {
	f, err := os.Open(flowConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Needs to include flow config file")
		} else {
			fmt.Printf("Failed to load config from %s: %s\n", flowConfigPath, err.Error())
		}

		os.Exit(1)
	}

	var conf *FlowConfig
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		msg := "error decoding flow.json, file properties and Config struct does not match, error: " + err.Error()
		log.Printf("[package:configuration][method:GetFlowConfig]" + msg)
		panic(msg)
	}

	return conf
}

type ProjectConfig struct {
	ProjectID string      `yaml:"project_id"`
	URLS      *URLSConfig `yaml:"urls"`
}

type URLSConfig struct {
	EventsURL      string `yaml:"events"`
	TicketTypesURL string `yaml:"ticket_types"`
}

func BuildProjectConfig(profile string) *ProjectConfig {
	path := configurationPackagePath + "/properties-" + profile + ".yml"
	return GetProjectConfig(path)

}

func GetProjectConfig(path string) *ProjectConfig {
	propertiesFile, err := os.Open(path)
	if err != nil {
		msg := "error getting properties file for path'" + path + "', error: " + err.Error()
		log.Print("[package:configuration][method:GetConfig]" + msg)
		panic(msg)
	}

	defer propertiesFile.Close()

	decoder := yaml.NewDecoder(propertiesFile)
	var cfg ProjectConfig

	err = decoder.Decode(&cfg)
	if err != nil {
		msg := "error decoding yaml, file properties and Config struct does not match, error: " + err.Error()
		log.Printf("[package:configuration][method:GetConfig]" + msg)
		panic(msg)
	}
	return &cfg
}
