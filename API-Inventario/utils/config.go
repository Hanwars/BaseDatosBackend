package utils

import (
	"errors"
	"fmt"

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

// LoadConfig carga la configuración desde app.env (desarrollo local) y/o desde
// variables de entorno (Docker / Northflank). Si app.env no existe no es un error:
// las variables de entorno tienen prioridad sobre el archivo.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("DB_DRIVER", "sqlserver")
	viper.SetDefault("SERVER_URL", "0.0.0.0:8080")
	viper.SetDefault("API_VERSION", "api/v1")
	viper.SetDefault("BACKUP_DIR", "")

	// BindEnv es necesario para que Unmarshal vea las variables de entorno
	// cuando la clave no está en un archivo de configuración.
	viper.AutomaticEnv()
	for _, key := range []string{
		"DB_DRIVER", "DB_SOURCE", "SERVER_URL", "API_VERSION", "SECRET", "BACKUP_DIR",
	} {
		if err = viper.BindEnv(key); err != nil {
			return
		}
	}

	if err = viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return
		}
		err = nil
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	if config.DBSource == "" {
		err = fmt.Errorf("falta DB_SOURCE (variable de entorno o app.env)")
		return
	}
	if config.SECRET == "" {
		err = fmt.Errorf("falta SECRET (variable de entorno o app.env)")
		return
	}
	return
}
