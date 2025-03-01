package config

import (
	"errors"
	"fmt"

	"github.com/ardanlabs/conf/v3"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Environment         string        `conf:"env:ENVIRONMENT,default:production"`
	DatabaseName        string        `conf:"env:DATABASE_NAME,default:postgres"`
	DatabaseUser        string        `conf:"env:DATABASE_USER,default:postgres"`
	DatabasePassword    string        `conf:"env:DATABASE_PASSWORD,required"`
	DatabaseHost        string        `conf:"env:DATABASE_HOST,default:localhost"`
	DatabasePort        string        `conf:"env:DATABASE_PORT,default:5432"`
	DatabaseSSLMode     string        `conf:"env:DATABASE_SSLMODE,default:disable"`
	ApiAddress          string        `conf:"env:API_ADDRESS,default:0.0.0.0:8000"`
	AuthSecretKey       string        `conf:"env:AUTH_SECRET_KEY,required"`
	DatabasePoolMinSize int32         `conf:"env:DATABASE_POOL_MIN_SIZE,default:2"`
	DatabasePoolMaxSize int32         `conf:"env:DATABASE_POOL_MAX_SIZE,default:10"`
	BundlerRpcUrl       string        `conf:"env:BUNDLER_RPC_URL"`
	RpcUrl              string        `conf:"env:RPC_URL"`
	TestPrivateKey      string        `conf:"env:PRIVATE_KEY"`
	ChainID             int64         `conf:"env:CHAIN_ID"`
	Logging             LoggingConfig `conf:"env:LOGGING"`
}

// LoggingConfig configures logging.
type LoggingConfig struct {
	Level  string `conf:"env:LOGGING_LEVEL"`
	Type   string `conf:"env:LOGGING_TYPE"`
	Stderr bool   `conf:"env:LOGGING_STDERR,default:false"`
}

func (c *Config) Load() error {
	if help, err := conf.Parse("", c); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return err
		}
		return err
	}
	return nil
}

func (c *Config) ConnectionString() string {
	if c.DatabaseSSLMode == "" {
		c.DatabaseSSLMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName,
		c.DatabaseSSLMode)
}
