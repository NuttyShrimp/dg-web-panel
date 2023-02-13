package config

import (
	"flag"
	"testing"
)

const (
	goodConfigPath = "../../config.yml"
	badConfigPath  = "/this/does/not/exist"
	dirConfigPath  = "/"
)

func TestLoadConfig(t *testing.T) {
	configPath, err := ValidateConfigFile(goodConfigPath)
	if err != nil {
		t.Error(err)
	}
	cPtr, err := LoadConfig(configPath)
	c := *cPtr
	switch {
	case err != nil:
		t.Errorf("Error loading config: %v", err)
	case c.Server.Host == "":
		t.Error("Server host is empty")
	case c.Server.Port == 0:
		t.Error("Server port must be greater than 0")
	case c.Server.Env == "":
		t.Error("Server env is empty")
	// Env must match production or development
	case c.Server.Env != "production" && c.Server.Env != "development":
		t.Error("Server env must be production or development")
	case c.Server.Cors.Origins == nil:
		t.Error("Server cors allowed origins is nil")
	case c.Server.SessionSecret == "":
		t.Error("Server session secret is empty")
	// Session secret must be at least 32 characters
	case len(c.Server.SessionSecret) < 32:
		t.Error("Server session secret must be at least 32 characters")
	case c.Discord.RedirectURL == "":
		t.Error("Discord redirectURL is empty")
	case c.Discord.ClientID == "":
		t.Error("Discord clientID is empty")
	case c.Discord.ClientSecret == "":
		t.Error("Discord clientSecret is empty")
	case c.Discord.Guild == "":
		t.Error("Discord guild is empty")
	case c.Discord.Roles == nil:
		t.Error("Discord roles are empty")
	case c.Redis.Host == "":
		t.Error("Redis host is empty")
	case c.Redis.Port == 0:
		t.Error("Redis port must be greater than 0")
	case c.Redis.Password == "":
		t.Error("Redis password is empty")
	case c.MariaDB.Panel.Host == "":
		t.Error("MariaDB.Panel host is empty")
	case c.MariaDB.Panel.Port == 0:
		t.Error("MariaDB.Panel port must be greater than 0")
	case c.MariaDB.Panel.User == "":
		t.Error("MariaDB.Panel user is empty")
	case c.MariaDB.Panel.Password == "":
		t.Error("MariaDB.Panel password is empty")
	case c.MariaDB.Panel.Database == "":
		t.Error("MariaDB.Panel database is empty")
	case c.MariaDB.Cfx.Host == "":
		t.Error("MariaDB.Cfx host is empty")
	case c.MariaDB.Cfx.Port == 0:
		t.Error("MariaDB.Cfx port must be greater than 0")
	case c.MariaDB.Cfx.User == "":
		t.Error("MariaDB.Cfx user is empty")
	case c.MariaDB.Cfx.Password == "":
		t.Error("MariaDB.Cfx password is empty")
	case c.MariaDB.Cfx.Database == "":
		t.Error("MariaDB.Cfx database is empty")
	case c.Graylog.URL == "":
		t.Error("Graylog endpoint is empty")
	case c.Graylog.StreamId == "":
		t.Error("Graylog targeted stream Id is empty")
	case c.Graylog.Token == "":
		t.Error("Graylog authentication token is empty")
	// CFX Config
	case c.Cfx.Server == "":
		t.Error("Cfx server can not be empty")
	case c.Cfx.ApiKey == "":
		t.Error("Cfx api key should not be empty")
	}
}

// TestLoadConfig2 test loading with an invalid config path
func TestLoadConfig2(t *testing.T) {
	_, err := LoadConfig(badConfigPath)
	if err == nil {
		t.Error("Error should have been thrown for bad config path")
	}
}

func TestValidateConfigFile(t *testing.T) {
	// Good config path
	if _, err := ValidateConfigFile(goodConfigPath); err != nil {
		t.Errorf("Error validating config file: %v", err)
	}

	// Bad config path
	if _, err := ValidateConfigFile(badConfigPath); err == nil {
		t.Error("Error should have been thrown for bad config path")
	}

	// Dir config path
	if _, err := ValidateConfigFile(dirConfigPath); err == nil {
		t.Error("Error should have been thrown for dir config path")
	}
}

func TestGetConfigPath(t *testing.T) {
	// Change the config path
	lookupFlag := flag.Lookup("config")
	if lookupFlag == nil {
		var configFlag string
		flag.StringVar(&configFlag, "config", "../../config.yml", "Path to config file")
	}
	err := flag.CommandLine.Set("config", "../../config.yml")
	if err != nil {
		t.Errorf("Error setting config flag: %v", err)
		return
	}

	if path := GetConfigPath(); path != goodConfigPath {
		t.Errorf("Config path should have been %s but was %s", goodConfigPath, path)
	}
}
