package config

// LogConfig
type LogConfig struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Level string `mapstructure:"level"`
}

// GinConfig
type GinConfig struct {
	Address string `mapstructure:"address"`
}

// DBConfig
type DBConfig struct {
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`

	LogMode        bool `mapstructure:"log_mode"`
	MaxIdle        int  `mapstructure:"max_idle"`
	MaxOpen        int  `mapstructure:"max_open"`
	ConnMaxLifeMin int  `mapstructure:"conn_max_life_min"`
}

func GetConfig() ConfigSetup {
	return config
}

func GetLogConfig() LogConfig {
	return config.LogConfig
}

func GetGinConfig() GinConfig {
	return config.GinConfig
}

func GetDBConfig() DBConfig {
	return config.DBConfig
}
