package config

import (
"os"
"path/filepath"

"gopkg.in/yaml.v3"
"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
homeDir, err := os.UserHomeDir()
if err != nil {
return "", err
}

configDir := filepath.Join(homeDir, ".config", "openapi-tui")
if err := os.MkdirAll(configDir, 0755); err != nil {
return "", err
}

return filepath.Join(configDir, "config.yaml"), nil
}

// LoadConfig loads configuration from the config file
func LoadConfig() models.Config {
cfg := models.Config{
VerboseMode: false,
}

configPath, err := GetConfigPath()
if err != nil {
return cfg
}

data, err := os.ReadFile(configPath)
if err != nil {
return cfg
}

var fileConfig models.ConfigFile
if err := yaml.Unmarshal(data, &fileConfig); err != nil {
return cfg
}

cfg.BaseURL = fileConfig.BaseURL
cfg.SpecPath = fileConfig.SpecPath
cfg.VerboseMode = fileConfig.VerboseMode

if fileConfig.Auth != nil {
cfg.Auth = &models.AuthConfig{
AuthType:   fileConfig.Auth.Type,
Token:      fileConfig.Auth.Token,
APIKeyIn:   fileConfig.Auth.APIKeyIn,
APIKeyName: fileConfig.Auth.APIKeyName,
Username:   fileConfig.Auth.Username,
Password:   fileConfig.Auth.Password,
}
}

return cfg
}

// SaveConfig saves the current configuration to the config file
func SaveConfig(cfg models.Config) error {
configPath, err := GetConfigPath()
if err != nil {
return err
}

fileConfig := models.ConfigFile{
BaseURL:     cfg.BaseURL,
SpecPath:    cfg.SpecPath,
VerboseMode: cfg.VerboseMode,
}

if cfg.Auth != nil {
fileConfig.Auth = &struct {
Type       string `yaml:"type"`
Token      string `yaml:"token,omitempty"`
APIKeyIn   string `yaml:"apiKeyIn,omitempty"`
APIKeyName string `yaml:"apiKeyName,omitempty"`
Username   string `yaml:"username,omitempty"`
Password   string `yaml:"password,omitempty"`
}{
Type:       cfg.Auth.AuthType,
Token:      cfg.Auth.Token,
APIKeyIn:   cfg.Auth.APIKeyIn,
APIKeyName: cfg.Auth.APIKeyName,
Username:   cfg.Auth.Username,
Password:   cfg.Auth.Password,
}
}

data, err := yaml.Marshal(fileConfig)
if err != nil {
return err
}

return os.WriteFile(configPath, data, 0644)
}
