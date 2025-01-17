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
	var (
		repoPath     = filepath.Join(os.Getenv("GOPATH"), "src/sekretariat")
		configPath   = filepath.Join(repoPath, "files/etc/sekretariat/sekretariat.development.yaml")
		namespace, _ = os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	)

	env := string(namespace)
	if os.Getenv("GOPATH") == "" {
		configPath = "files/etc/sekretariat/sekretariat.development.yaml"
	}

	if env != "" {
		if env == envStaging {
			configPath = "files/etc/sekretariat/sekretariat.staging.yaml"
		} else if env == envProduction {
			configPath = "files/etc/sekretariat/sekretariat.production.yaml"
		}
	}
	return configPath
}

// Get ...
func Get() *Config {
	return config
}
