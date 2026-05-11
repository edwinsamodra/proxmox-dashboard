package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Proxmox  ProxmoxConfig  `yaml:"proxmox"`
	Auth     AuthConfig     `yaml:"auth"`
	OpenTofu OpenTofuConfig `yaml:"opentofu"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ProxmoxConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Insecure bool   `yaml:"insecure"`
	// Credentials are loaded from env vars at runtime.
}

type AuthConfig struct {
	// Secret used to sign session tokens.
	Secret string `yaml:"secret"`
}

type OpenTofuConfig struct {
	// Directory where .tf workspace files are stored.
	WorkspaceDir string `yaml:"workspace_dir"`
	// Path to the opentofu binary (defaults to "tofu" on PATH).
	BinaryPath string `yaml:"binary_path"`
}

// Load reads configuration from a YAML file and overrides values with
// environment variables where applicable.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Proxmox: ProxmoxConfig{
			Port:     8006,
			Insecure: true,
		},
		Auth: AuthConfig{
			Secret: "change-me-in-production",
		},
		OpenTofu: OpenTofuConfig{
			WorkspaceDir: "/tmp/pmo-workspaces",
			BinaryPath:   "tofu",
		},
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		if err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config file: %w", err)
			}
		}
	}

	// Environment variable overrides.
	if v := os.Getenv("PMO_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("PMO_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = p
		}
	}
	if v := os.Getenv("PMO_PROXMOX_HOST"); v != "" {
		cfg.Proxmox.Host = v
	}
	if v := os.Getenv("PMO_PROXMOX_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Proxmox.Port = p
		}
	}
	if v := os.Getenv("PMO_SECRET"); v != "" {
		cfg.Auth.Secret = v
	}
	if v := os.Getenv("PMO_TOFU_DIR"); v != "" {
		cfg.OpenTofu.WorkspaceDir = v
	}

	return cfg, nil
}

// Addr returns the listen address for the HTTP server.
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// BaseURL returns the Proxmox API base URL.
func (c *ProxmoxConfig) BaseURL() string {
	scheme := "https"
	return fmt.Sprintf("%s://%s:%d/api2/json", scheme, c.Host, c.Port)
}
