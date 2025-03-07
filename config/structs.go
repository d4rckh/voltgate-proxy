package config

import (
	"voltgate-proxy/storage"
)

type RateLimitConfig struct {
	Rules []RateLimitRule `yaml:"rules"`
}

type CacheConfig struct {
	Rules []CacheRule `yaml:"rules"`
}

type AppRateLimitRules struct {
	EndpointRateLimitRules map[string][]RateLimitRule
	ServicesRateLimitRules map[string][]RateLimitRule
}

type AppCacheRules struct {
	EndpointCacheRules map[string][]CacheRule
}

type Endpoint struct {
	Host            string          `yaml:"host"`
	Service         string          `yaml:"service"`
	RateLimitConfig RateLimitConfig `yaml:"rate_limit"`
	CacheConfig     CacheConfig     `yaml:"cache"`
}

type Service struct {
	Url             string          `yaml:"url"`
	Name            string          `yaml:"name"`
	RateLimitConfig RateLimitConfig `yaml:"rate_limit"`
}

type MonitoringAppConfig struct {
	LokiUrl           string `yaml:"loki"`
	PrometheusEnabled bool   `yaml:"prometheus"`
}

type RateLimitRule struct {
	Path             string `yaml:"path"`
	NumberOfRequests int    `yaml:"requests"`
	WindowSeconds    int    `yaml:"window"`
	Action           string `yaml:"action"`
	Method           string `yaml:"method"`
}

type CacheRule struct {
	Path    string   `yaml:"path"`
	Ttl     int      `yaml:"ttl"`
	Params  []string `yaml:"params"`
	Methods []string `yaml:"methods"`
}

type RateLimitAppConfig struct {
	Storage string `yaml:"storage"`
}

type StorageAppConfig struct {
	Redis storage.RedisAppConfig `yaml:"redis"`
}

type AppProxyConfig struct {
	Address string `yaml:"address"`
}

type AppManagementConfig struct {
	Address string `yaml:"address"`
}

type CacheAppConfig struct {
	Storage string `yaml:"storage"`
}

type AppConfig struct {
	Services             []Service  `yaml:"services"`
	Endpoints            []Endpoint `yaml:"endpoints"`
	ReloadConfigInterval int        `yaml:"config.reload_interval"`

	ProxyConfig      AppProxyConfig      `yaml:"proxy"`
	ManagementConfig AppManagementConfig `yaml:"management"`

	Storage StorageAppConfig `yaml:"storage"`

	MonitoringAppConfig MonitoringAppConfig `yaml:"monitoring"`
	RateLimitAppConfig  RateLimitAppConfig  `yaml:"rate_limit"`
	CacheAppConfig      CacheAppConfig      `yaml:"cache"`
}
