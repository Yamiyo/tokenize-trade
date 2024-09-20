package config

import (
	"os"
	"path/filepath"
	"time"
	"tokenize-trade/internal/binance"

	"github.com/spf13/viper"
)

type EnvType = string

const (
	EnvTypeLocal EnvType = "local"
	EnvTypeDev   EnvType = "dev"
	EnvTypeProd  EnvType = "prod"

	defaultEnv EnvType = EnvTypeLocal
)

// Config ...
var (
	config   ConfigSetup
	TimeZone *time.Location
)

func InitConfig() error {
	if err := LoadConfig("conf.d/config.yaml"); err != nil {
		return err
	}

	return nil
}

// LoadConfig ...
func LoadConfig(file string) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

// ConfigSetup
type ConfigSetup struct {
	LogConfig     LogConfig          `mapstructure:"log_config"`
	GinConfig     GinConfig          `mapstructure:"gin_config"`
	DBConfig      DBConfig           `mapstructure:"db_config"`
	BinanceConfig binance.BinanceCfg `mapstructure:"binance_config"`
}
