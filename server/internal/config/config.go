package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Read the config from the yaml file

type Config struct {
	Server  ConfigServer  `yaml:"server"`
	Discord ConfigDiscord `yaml:"discord"`
	Redis   ConfigRedis   `yaml:"redis"`
	MariaDB ConfigMariaDB `yaml:"mariaDB"`
	Graylog ConfigGraylog `yaml:"graylog"`
	Cfx     ConfigCfx     `yaml:"cfx"`
}

type ConfigServer struct {
	Host      string  `yaml:"host"`
	Port      int     `yaml:"port"`
	Env       string  `yaml:"env"`
	ReqPerSeq float64 `yaml:"reqpsec"`
	Cors      struct {
		Origins []string `yaml:"origins"`
	} `yaml:"cors"`
	SessionSecret string `yaml:"sessionSecret"`
}

type ConfigDiscord struct {
	RedirectURL  string       `yaml:"redirectURL"`
	ClientID     string       `yaml:"clientID"`
	ClientSecret string       `yaml:"clientSecret"`
	Guild        string       `yaml:"guild"`
	Roles        []ConfigRole `yaml:"roles"`
}

type ConfigRole struct {
	Name string `yaml:"name"`
	Id   uint64 `yaml:"id"`
}

type ConfigRedis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type ConfigMariaDB struct {
	Panel ConfigMariaConn `yaml:"panel"`
	Cfx   ConfigMariaConn `yaml:"cfx"`
}

type ConfigMariaConn struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ConfigGraylog struct {
	URL           string `yaml:"url"`
	Token         string `yaml:"token"`
	StreamId      string `yaml:"streamId"`
	PanelStreamId string `yaml:"panelStreamId"`
	Gelf          string `yaml:"gelf"`
}

type ConfigCfx struct {
	Server string `yaml:"server"`
	ApiKey string `yaml:"apiKey"`
}

/*
Load the config from the given path
Path should start from the root of the project
*/
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprint("Error closing config file: ", err))
		}
	}(file)

	// New yaml decoder
	d := yaml.NewDecoder(file)

	// Decode the file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigFile(configPath string) (string, error) {
	// Check if path exists and is a file
	s, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Config file does not exist: %v", configPath)
	}

	if s.IsDir() {
		return "", fmt.Errorf("Config file is a directory: %v", configPath)
	}
	return configPath, nil
}

// Read the path from flags or the default
func GetConfigPath() string {
	// Get config flag without redefining it
	configFlag := flag.Lookup("config")
	var configPath string
	if (configFlag == nil) || (configFlag.Value.String() == "") {
		// Set the config flag
		flag.StringVar(&configPath, "config", "./config.yml", "Path to config file")
		configFlag = flag.Lookup("config")
	}
	configPath = configFlag.Value.String()

	flag.Parse()

	configPath, err := ValidateConfigFile(configPath)
	if err != nil {
		panic(err)
	}
	return configPath
}
