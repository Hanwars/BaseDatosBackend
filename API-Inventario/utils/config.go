package utils

import "github.com/spf13/viper"

type Config struct {
	DBDriver  string `mapstructure:"DB_DRIVER"`
	DBSource  string `mapstructure:"DB_SOURCE"`
	ServerURL string `mapstructure:"SERVER_URL"`
	API       string `mapstructure:"API_VERSION"`
	SECRET    string `mapstructure:"SECRET"`
	BackupDir string `mapstructure:"BACKUP_DIR"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
