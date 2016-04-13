package config

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Website Website
}

type Website struct {
	// URL from which the site can be accessed.
	URL string
	// Port over which HTTP is served
	HTTPPort string
	// Port over which HTTPS is served
	HTTPSPort string
	// Path to the TLS certificate
	Cert string
	// Path to the TLS key
	Key string
	// Path to the website directory
	Directory string
}

// New constructs a new config file with default values set where applicable.
func New() *Config {
	return &Config{
		Website: Website{
			URL:       "localhost",
			HTTPPort:  ":80",
			HTTPSPort: ":443",
			Cert:      "cert.pem",
			Key:       "key.pem",
			Directory: "app/public",
		},
	}
}

// ReadFile reads in config information from the file with the given name,
// assumed to be in gcfg format, overwriting any existing values. Returns an
// error if the file cannot be read, nil otherwise.
func ReadFile(filename string) (*Config, error) {
	config := New()
	err := gcfg.ReadFileInto(config, filename)
	if err != nil {
		return nil, err
	}
	return config, nil
}
