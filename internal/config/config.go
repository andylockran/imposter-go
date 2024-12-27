package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Response represents an HTTP response
type Response struct {
	Content    string `yaml:"content"`
	StatusCode int    `yaml:"statusCode"`
	File       string `yaml:"file"`
}

// Resource represents an HTTP resource
type Resource struct {
	Method      string            `yaml:"method"`
	Path        string            `yaml:"path"`
	QueryParams map[string]string `yaml:"queryParams"`
	Headers     map[string]string `yaml:"headers"`
	RequestBody map[string]string `yaml:"requestBody"`
	FormParams  map[string]string `yaml:"formParams"`
	PathParams  map[string]string `yaml:"pathParams"` // new field
	Response    Response          `yaml:"response"`
}

type Config struct {
	Plugin    string `yaml:"plugin"`
	Resources []Resource
}

// Application-wide configuration
type ImposterConfig struct {
	ServerPort string
}

// LoadImposterConfig loads configurations from environment variables
func LoadImposterConfig() *ImposterConfig {
	port := os.Getenv("IMPOSTER_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	return &ImposterConfig{
		ServerPort: port,
	}
}

// LoadConfig loads all config files in the specified directory
func LoadConfig(configDir string) []Config {
	var configs []Config

	scanRecursive := (os.Getenv("IMPOSTER_CONFIG_SCAN_RECURSIVE") == "true")

	err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip subdirectories if not scanning recursively
		if info.IsDir() && info.Name() != filepath.Base(configDir) && !scanRecursive {
			return filepath.SkipDir
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), "-config.json") || strings.HasSuffix(info.Name(), "-config.yaml") || strings.HasSuffix(info.Name(), "-config.yml")) {
			fmt.Printf("Loading config file: %s\n", path)
			fileConfig, err := parseConfig(path)
			if err != nil {
				return err
			}
			// Prefix 'File' properties if in a subdirectory
			baseDir := filepath.Dir(path)
			relDir, err := filepath.Rel(configDir, baseDir)
			if err != nil {
				return err
			}
			for i := range fileConfig.Resources {
				if fileConfig.Resources[i].Response.File != "" && relDir != "." {
					fileConfig.Resources[i].Response.File = filepath.Join(relDir, fileConfig.Resources[i].Response.File)
				}
			}
			configs = append(configs, *fileConfig)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return configs
}

// parseConfig loads and parses a YAML configuration file
func parseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &cfg, nil
}
