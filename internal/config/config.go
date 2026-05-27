package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var defaultCapturePaths = []string{
	"/v1/chat/completions",
	"/v1/responses",
	"/v1/completions",
	"/v1/embeddings",
	"/v1/messages",
}

type Config struct {
	ListenAddr        string   `yaml:"listen_addr"`
	UpstreamBase      string   `yaml:"upstream_base"`
	// Deprecated: retained so older config files with web_addr still load.
	WebAddr           string   `yaml:"web_addr"`
	WebBasePath       string   `yaml:"web_base_path"`
	PostgresDSN       string   `yaml:"postgres_dsn"`
	HMACSecret        string   `yaml:"hmac_secret"`
	AdminUsername     string   `yaml:"admin_username"`
	AdminPasswordHash string   `yaml:"admin_password_hash"`
	MaxBodyBytes      int64    `yaml:"max_body_bytes"`
	CapturePaths      []string `yaml:"capture_paths"`
}

func Load(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	cfg.applyDefaults()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) applyDefaults() {
	if c.ListenAddr == "" {
		c.ListenAddr = "127.0.0.1:3007"
	}
	if c.MaxBodyBytes <= 0 {
		c.MaxBodyBytes = 10 * 1024 * 1024
	}
	if len(c.CapturePaths) == 0 {
		c.CapturePaths = append([]string(nil), defaultCapturePaths...)
	}

	c.ListenAddr = strings.TrimSpace(c.ListenAddr)
	c.UpstreamBase = strings.TrimRight(c.UpstreamBase, "/")
	c.WebBasePath = normalizeBasePath(c.WebBasePath)
}

func (c Config) Validate() error {
	switch {
	case c.ListenAddr == "":
		return fmt.Errorf("listen_addr is required")
	case c.UpstreamBase == "":
		return fmt.Errorf("upstream_base is required")
	case c.WebBasePath == "":
		return fmt.Errorf("web_base_path is required")
	case c.PostgresDSN == "":
		return fmt.Errorf("postgres_dsn is required")
	case c.HMACSecret == "":
		return fmt.Errorf("hmac_secret is required")
	case c.AdminUsername == "":
		return fmt.Errorf("admin_username is required")
	case c.AdminPasswordHash == "":
		return fmt.Errorf("admin_password_hash is required")
	}

	parsed, err := url.Parse(c.UpstreamBase)
	if err != nil {
		return fmt.Errorf("invalid upstream_base: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("upstream_base must use http or https")
	}
	if parsed.Host == "" {
		return fmt.Errorf("upstream_base host is required")
	}
	if c.WebBasePath == "/" {
		return fmt.Errorf("web_base_path cannot be / when audit web shares the proxy port")
	}
	if !strings.HasPrefix(c.WebBasePath, "/") {
		return fmt.Errorf("web_base_path must start with /")
	}
	if strings.ContainsAny(c.WebBasePath, "?#") {
		return fmt.Errorf("web_base_path cannot contain query or fragment")
	}

	return nil
}

func normalizeBasePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "/audit"
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	if value != "/" {
		value = strings.TrimRight(value, "/")
	}
	return value
}
