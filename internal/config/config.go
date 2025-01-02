package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Response represents an HTTP response
type Response struct {
	Content    string            `yaml:"content"`
	StatusCode int               `yaml:"statusCode"`
	File       string            `yaml:"file"`
	Fail       string            `yaml:"fail"`
	Delay      Delay             `yaml:"delay"`
	Headers    map[string]string `yaml:"headers"`
	Template   bool              `yaml:"template"`
}

// Delay represents the delay configuration for a response
type Delay struct {
	Exact int `yaml:"exact"`
	Min   int `yaml:"min"`
	Max   int `yaml:"max"`
}

// Matcher represents anything that can be matched against a value
type Matcher interface {
	Match(actualValue string) bool
}

// StringMatcher is a simple string matcher that checks for exact equality
type StringMatcher string

func (s StringMatcher) Match(actualValue string) bool {
	return string(s) == actualValue
}

// MatchCondition represents a condition for matching requests
type MatchCondition struct {
	Value    string `yaml:"value"`
	Operator string `yaml:"operator"`
}

func (m MatchCondition) Match(actualValue string) bool {
	switch m.Operator {
	case "EqualTo", "":
		return actualValue == m.Value
	case "NotEqualTo":
		return actualValue != m.Value
	case "Exists":
		return actualValue != ""
	case "NotExists":
		return actualValue == ""
	case "Contains":
		return strings.Contains(actualValue, m.Value)
	case "NotContains":
		return !strings.Contains(actualValue, m.Value)
	case "Matches":
		matched, _ := regexp.MatchString(m.Value, actualValue)
		return matched
	case "NotMatches":
		matched, _ := regexp.MatchString(m.Value, actualValue)
		return !matched
	default:
		return false
	}
}

// BodyMatchCondition represents a condition for matching request bodies
type BodyMatchCondition struct {
	MatchCondition
	JSONPath      string            `yaml:"jsonPath,omitempty"`
	XPath         string            `yaml:"xPath,omitempty"`
	XMLNamespaces map[string]string `yaml:"xmlNamespaces"`
}

func (b BodyMatchCondition) Match(actualValue string) bool {
	return b.MatchCondition.Match(actualValue)
}

// RequestBody represents the request body matching configuration
type RequestBody struct {
	BodyMatchCondition
	AllOf []BodyMatchCondition `yaml:"allOf"`
	AnyOf []BodyMatchCondition `yaml:"anyOf"`
}

// Capture defines how to capture request data for later use in the response
type Capture struct {
	Enabled    bool       `yaml:"enabled,omitempty"`
	Store      string     `yaml:"store"`
	Key        CaptureKey `yaml:"key,omitempty"`
	CaptureKey `yaml:",inline"`
}

// CaptureKey represents the key configuration for capturing request data.
type CaptureKey struct {
	PathParam     string `yaml:"pathParam,omitempty"`
	QueryParam    string `yaml:"queryParam,omitempty"`
	FormParam     string `yaml:"formParam,omitempty"`
	RequestHeader string `yaml:"requestHeader,omitempty"`
	Expression    string `yaml:"expression,omitempty"`
	Const         string `yaml:"const,omitempty"`
	RequestBody   struct {
		JSONPath      string            `yaml:"jsonPath,omitempty"`
		XPath         string            `yaml:"xPath,omitempty"`
		XMLNamespaces map[string]string `yaml:"xmlNamespaces,omitempty"`
	} `yaml:"requestBody,omitempty"`
}

// RequestMatcher contains the common fields for matching requests
type RequestMatcher struct {
	Method      string                        `yaml:"method"`
	Path        string                        `yaml:"path"`
	QueryParams map[string]MatcherUnmarshaler `yaml:"queryParams"`
	Headers     map[string]MatcherUnmarshaler `yaml:"headers"`
	RequestBody RequestBody                   `yaml:"requestBody"`
	FormParams  map[string]MatcherUnmarshaler `yaml:"formParams"`
	PathParams  map[string]MatcherUnmarshaler `yaml:"pathParams"`
	Capture     map[string]Capture            `yaml:"capture,omitempty"`
}

// Resource represents an HTTP resource
type Resource struct {
	RequestMatcher `yaml:",inline"`
	Response       Response `yaml:"response"`
}

// Interceptor represents an HTTP interceptor that can be executed before resources
type Interceptor struct {
	RequestMatcher `yaml:",inline"`
	Response       *Response `yaml:"response,omitempty"`
	Continue       bool      `yaml:"continue"`
}

type System struct {
	Stores map[string]StoreDefinition `yaml:"stores"`
}

type StoreDefinition struct {
	PreloadFile string                 `yaml:"preloadFile,omitempty"`
	PreloadData map[string]interface{} `yaml:"preloadData,omitempty"`
}

type Config struct {
	Plugin       string `yaml:"plugin"`
	BasePath     string `yaml:"basePath"`
	Resources    []Resource
	Interceptors []Interceptor `yaml:"interceptors,omitempty"`
	System       *System       `yaml:"system,omitempty"`
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
	autoBasePath := (os.Getenv("IMPOSTER_AUTO_BASE_PATH") == "true")

	ignorePaths := loadIgnorePaths(configDir)

	err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored paths
		if shouldIgnorePath(path, ignorePaths) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
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

			// Set basePath if autoBasePath is enabled
			if autoBasePath && fileConfig.BasePath == "" {
				baseDir := filepath.Dir(path)
				relDir, err := filepath.Rel(configDir, baseDir)
				if err != nil {
					return err
				}
				fileConfig.BasePath = "/" + relDir
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
				// Prefix paths with basePath
				if fileConfig.BasePath != "" {
					fileConfig.Resources[i].Path = filepath.Join(fileConfig.BasePath, fileConfig.Resources[i].Path)
				}
			}
			configs = append(configs, *fileConfig)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// TODO remove guard for rest plugin support
	for _, cfg := range configs {
		if cfg.Plugin != "rest" {
			panic("Unsupported plugin type")
		}
	}

	return configs
}

// loadIgnorePaths loads ignore paths from .imposterignore file or uses default ignore paths
func loadIgnorePaths(configDir string) []string {
	defaultIgnorePaths := []string{".git", ".idea", ".svn", "node_modules"}
	ignoreFilePath := filepath.Join(configDir, ".imposterignore")

	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		return defaultIgnorePaths
	}

	data, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		fmt.Printf("Failed to read .imposterignore file: %v\n", err)
		return defaultIgnorePaths
	}

	lines := strings.Split(string(data), "\n")
	var ignorePaths []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			ignorePaths = append(ignorePaths, line)
		}
	}

	return ignorePaths
}

// shouldIgnorePath checks if a path should be ignored based on the ignore paths
func shouldIgnorePath(path string, ignorePaths []string) bool {
	for _, ignorePath := range ignorePaths {
		if strings.Contains(path, ignorePath) {
			return true
		}
	}
	return false
}

// parseConfig loads and parses a YAML configuration file
func parseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Substitute environment variables
	data = []byte(substituteEnvVars(string(data)))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &cfg, nil
}

// substituteEnvVars replaces ${env.VAR} and ${env.VAR:-default} with environment variable values
func substituteEnvVars(content string) string {
	re := regexp.MustCompile(`\$\{env\.([A-Z0-9_]+)(:-([^}]+))?\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		groups := re.FindStringSubmatch(match)
		envVar := groups[1]
		defaultValue := groups[3]
		if value, exists := os.LookupEnv(envVar); exists {
			return value
		}
		return defaultValue
	})
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Matcher
func (m *MatchCondition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// First try to unmarshal as a simple string
	var str string
	if err := unmarshal(&str); err == nil {
		*m = MatchCondition{Value: str}
		return nil
	}

	// If that fails, try to unmarshal as a MatchCondition struct
	type matchConditionAlias MatchCondition
	return unmarshal((*matchConditionAlias)(m))
}

// MatcherUnmarshaler is a helper type for unmarshaling Matcher from YAML
type MatcherUnmarshaler struct {
	Matcher Matcher
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for MatcherUnmarshaler
func (mu *MatcherUnmarshaler) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// First try to unmarshal as a simple string
	var str string
	if err := unmarshal(&str); err == nil {
		mu.Matcher = StringMatcher(str)
		return nil
	}

	// If that fails, try to unmarshal as a MatchCondition
	var mc MatchCondition
	if err := unmarshal(&mc); err == nil {
		mu.Matcher = mc
		return nil
	}

	return fmt.Errorf("failed to unmarshal as either string or MatchCondition")
}
