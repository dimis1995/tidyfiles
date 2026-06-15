package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DryRun            bool     `mapstructure:"dry_run"`
	AllowedExtensions []string `mapstructure:"allowed_extensions"`
	Size              struct {
		MaxSize  int32  `mapstructure:"max_size"`
		SizeUnit string `mapstructure:"size_unit"`
	} `mapstructure:"size"`
}

func (service Config) PrintConfiguration() {
	jsonConfig, _ := json.Marshal(service)
	fmt.Println(string(jsonConfig))
}

func (service Config) SaveToFile(filename string) error {
	if filename == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("determining home directory: %w", err)
		}
		filename = filepath.Join(home, ".config", "tidyfiles", "config.toml")
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := viper.WriteConfigAs(filename); err != nil {
		return fmt.Errorf("writing config to %s: %w", filename, err)
	}

	return nil
}

var AppConfig *Config

func setDefaults() {
	viper.SetDefault("dry_run", false)
	viper.SetDefault("allowed_extensions", []string{"pdf", "jpg", "png"})
	viper.SetDefault("size.max_size", 100)
	viper.SetDefault("size.size_unit", "MB")
}

func LoadConfig(filename string) (*Config, error) {
	setDefaults()
	if filename != "" {
		viper.SetConfigFile(filename)
		fmt.Println("Configuration file: " + filename)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("$HOME/.config/tidyfiles")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()

	if err != nil {
		if fileNotFoundError, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			slog.Debug(fileNotFoundError.Error())
			slog.Warn("Cannot find a configuration file. Will use defaults")
		} else {
			slog.Error("Unexpected error when reading the configuration", "err", err)
			os.Exit(1)
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return AppConfig, nil

}
