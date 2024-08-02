package config

import (
	"os"

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

// Init ...
func Init(opts ...Option) error {
	opt := &option{
		configFile: getDefaultConfigFile(),
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	out, err := os.ReadFile(opt.configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(out, &config)
}

// Option ...
type Option func(*option)

// WithConfigFile ...
func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

func getDefaultConfigFile() string {
	env := os.Getenv("ENVIRONMENT") // Example: "development", "staging", "production"
	if env == "" {
		env = envDevelopment // Default environment
	}

	var configFile string
	switch env {
	case envStaging:
		configFile = "./sekretariat.staging.yaml"
	case envProduction:
		configFile = "./sekretariat.production.yaml"
	default:
		configFile = "./sekretariat.development.yaml"
	}
	return configFile
}

// Get ...
func Get() *Config {
	return config
}
