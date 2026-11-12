package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// LoadConfig carga app.env si existe y siempre prioriza las variables de entorno.
// Render proporciona PORT dinámicamente; por eso PORT tiene prioridad sobre SERVER_URL.
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
		API:       get("API_VERSION"),
		SECRET:    get("SECRET"),
		BackupDir: get("BACKUP_DIR"),
	}

	if config.DBDriver == "" {
		config.DBDriver = "sqlserver"
	}

	// Render define PORT en tiempo de ejecución (normalmente 10000).
	// Se fuerza 0.0.0.0 para que el proxy de Render pueda alcanzar el servicio.
	if port := strings.TrimSpace(get("PORT")); port != "" {
		config.ServerURL = "0.0.0.0:" + port
	} else if serverURL := strings.TrimSpace(get("SERVER_URL")); serverURL != "" {
		config.ServerURL = serverURL
	} else {
		config.ServerURL = "0.0.0.0:10000"
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
	if len(config.SECRET) != 32 {
		return config, fmt.Errorf("SECRET debe tener exactamente 32 bytes para PASETO v2 local; longitud recibida: %d", len(config.SECRET))
	}

	return config, nil
}
