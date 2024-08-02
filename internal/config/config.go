package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
)

const (
	envDevelopment = "development"
	envStaging     = "jx-staging"
	envProduction  = "jx-production"
)

type option struct {
	configFile string
}

// Init initializes the configuration by reading environment variables or a file
func Init(opts ...Option) error {
	opt := &option{
		configFile: getDefaultConfigFile(),
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	// Check if environment variables are set
	if os.Getenv("DATABASE_URL") != "" {
		config = &Config{
			Server: ServerConfig{
				Port: os.Getenv("PORT"),
			},
			Database: DatabaseConfig{
				Master: os.Getenv("DATABASE_URL"),
			},
		}
		return nil
	}

	// Otherwise, read from the config file
	out, err := os.ReadFile(opt.configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(out, &config)
}

// Option defines a functional option for the configuration
type Option func(*option)

// WithConfigFile allows specifying a custom configuration file
func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

func getDefaultConfigFile() string {
	if os.Getenv("GOPATH") == "" {
		return "./sekretariat.development.yaml"
	}

	namespace, _ := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	env := string(namespace)

	switch env {
	case envStaging:
		return "./sekretariat.staging.yaml"
	case envProduction:
		return "./sekretariat.production.yaml"
	default:
		return "./sekretariat.development.yaml"
	}
}

// Get returns the current configuration
func Get() *Config {
	return config
}
