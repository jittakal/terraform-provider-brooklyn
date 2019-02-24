package config

import "github.com/apache/brooklyn-client/cli/net"

// Config structure holds Compose related configuration
type Config struct {
	AccessKey     string `json:"access_key,omitempty"`
	SecretKey     string `json:"secret_key,omitempty"`
	EndpointURL   string `json:"endpoint_url,omitempty"`
	SkipSslChecks bool   `json:"skip_ssl_checks,omitempty"`
}

// Client reference to *net.Network object
type Client = *net.Network

// Client configures and returns a fully initialized Client
func (c *Config) Client() (interface{}, error) {
	client := net.NewNetwork(c.EndpointURL, c.AccessKey, c.SecretKey, c.SkipSslChecks)
	return client, nil
}
