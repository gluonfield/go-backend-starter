package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/augustdev/autoclip/internal/storage/pg"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewConfig() (Config, error) {
	if err := InitConfig(); err != nil {
		return Config{}, err
	}

	var cfg struct {
		GraphQLAPI GraphQLAPIConfig
		Db         pg.DbConfig
		OpenRouter OpenRouterConfig
		OpenAI     OpenAIConfig
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to decode config: %w", err)
	}

	return Config{
		GraphQLAPI: cfg.GraphQLAPI,
		Db:         cfg.Db,
	}, nil
}

type Config struct {
	fx.Out

	GraphQLAPI GraphQLAPIConfig
	Db         pg.DbConfig
	OpenRouter OpenRouterConfig
	OpenAI     OpenAIConfig
}

type OpenRouterConfig struct {
	APIKey string
}

type OpenAIConfig struct {
	APIKey string
}

type AuditConfig struct {
	OrderMatchingEnabled bool
}

// TracingConfig holds configuration for OpenTelemetry tracing
type TracingConfig struct {
	Enabled             bool     `mapstructure:"enabled"`
	ServiceName         string   `mapstructure:"serviceName"`
	Endpoint            string   `mapstructure:"endpoint"`
	SamplingRate        float64  `mapstructure:"samplingRate"`
	Insecure            bool     `mapstructure:"insecure"`
	Provider            string   `mapstructure:"provider"`
	ParentBasedSampling bool     `mapstructure:"parentBasedSampling"`
	AlwaysOnUserIds     []string `mapstructure:"alwaysOnUserIds"`
	AlwaysOnOperations  []string `mapstructure:"alwaysOnOperations"`
}

func InitConfig() error {
	_ = godotenv.Load()

	// Check for APPLICATION_CONFIG environment variable
	if configFile := os.Getenv("APPLICATION_CONFIG"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("application")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}
