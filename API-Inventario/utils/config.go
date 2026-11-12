package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver  string `mapstructure:"DB_DRIVER"`
	DBSource  string `mapstructure:"DB_SOURCE"`
	ServerURL string `mapstructure:"SERVER_URL"`
	API       string `mapstructure:"API_VERSION"`
	SECRET    string `mapstructure:"SECRET"`
	BackupDir string `mapstructure:"BACKUP_DIR"`
}

// LoadConfig loads app.env when present and then applies environment variables.
// This supports local development as well as Docker/Northflank/Azure deployments.
func LoadConfig(path string) (config Config, err error) {
	v := viper.New()
	configPath := filepath.Join(path, "app.env")
	v.SetConfigFile(configPath)
	v.SetConfigType("env")

	if readErr := v.ReadInConfig(); readErr != nil {
		if _, statErr := os.Stat(configPath); statErr == nil {
			return config, fmt.Errorf("no se pudo leer app.env: %w", readErr)
		} else if !os.IsNotExist(statErr) {
			return config, fmt.Errorf("no se pudo acceder a app.env: %w", statErr)
		}
	}

	get := func(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return v.GetString(key)
	}

	config = Config{
		DBDriver:  get("DB_DRIVER"),
		DBSource:  get("DB_SOURCE"),
		ServerURL: get("SERVER_URL"),
		API:       get("API_VERSION"),
		SECRET:    get("SECRET"),
		BackupDir: get("BACKUP_DIR"),
	}

	if config.DBDriver == "" {
		config.DBDriver = "sqlserver"
	}
	if config.ServerURL == "" {
		config.ServerURL = "0.0.0.0:8080"
	}
	if config.API == "" {
		config.API = "api/v1"
	}
	if config.DBSource == "" {
		return config, fmt.Errorf("falta configurar DB_SOURCE")
	}
	if config.SECRET == "" {
		return config, fmt.Errorf("falta configurar SECRET")
	}

	return config, nil
}
