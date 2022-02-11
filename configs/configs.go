package configs

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/instill-ai/pipeline-backend/internal/logger"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"go.temporal.io/sdk/client"
)

// Config - Global variable to export
var Config AppConfig

// AppConfig defines
type AppConfig struct {
	Server       ServerConfig       `koanf:"server"`
	Database     DatabaseConfig     `koanf:"database"`
	Temporal     TemporalConfig     `koanf:"temporal"`
	Cache        CacheConfig        `koanf:"cache"`
	VDO          VDOConfig          `koanf:"vdo"`
	ModelBackend ModelBackendConfig `koanf:"modelbackend"`
}

// ServerConfig defines HTTP server configurations
type ServerConfig struct {
	Port  int `koanf:"port"`
	HTTPS struct {
		Enabled bool   `koanf:"enabled"`
		Cert    string `koanf:"cert"`
		Key     string `koanf:"key"`
	}
	CORSOrigins []string `koanf:"corsorigins"`
	Paginate    struct {
		Salt string `koanf:"salt"`
	}
}

// Configs related to database
type DatabaseConfig struct {
	Username     string `koanf:"username"`
	Password     string `koanf:"password"`
	Host         string `koanf:"host"`
	Port         int    `koanf:"port"`
	DatabaseName string `koanf:"databasename"`
	TimeZone     string `koanf:"timezone"`
	Pool         struct {
		IdleConnections int           `koanf:"idleconnections"`
		MaxConnections  int           `koanf:"maxconnections"`
		ConnLifeTime    time.Duration `koanf:"connlifetime"`
	}
}

// Configs related to Temporal
type TemporalConfig struct {
	ClientOptions client.Options `koanf:"clientoptions"`
}

type CacheConfig struct {
	Redis struct {
		RedisOptions redis.Options `koanf:"redisoptions"`
	}
}

type VDOConfig struct {
	Scheme string `koanf:"scheme"`
	Host   string `koanf:"host"`
	Port   int    `koanf:"port"`
	Path   string `koanf:"path"`
}

type ModelBackendConfig struct {
	Scheme string `koanf:"scheme"`
	Host   string `koanf:"host"`
	Port   int    `koanf:"port"`
}

// Init - Assign global config to decoded config struct
func Init() error {
	logger, _ := logger.GetZapLogger()

	k := koanf.New(".")
	parser := yaml.Parser()

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fileRelativePath := fs.String("file", "configs/config.yaml", "configuration file")
	flag.Parse()

	if err := k.Load(file.Provider(*fileRelativePath), parser); err != nil {
		logger.Fatal(err.Error())
	}

	k.Load(env.Provider("CFG_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CFG_")), "_", ".", -1)
	}), nil)

	k.Unmarshal("", &Config)

	return ValidateConfig(&Config)
}

// ValidateConfig is for custom validation rules for the configuration
func ValidateConfig(cfg *AppConfig) error {
	return nil
}
