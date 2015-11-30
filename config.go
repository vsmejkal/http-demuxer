package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port      int
	Forwards  map[string]string
	Redirects map[string]string
}

func (c *Config) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open config file '%s'", path)
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf("Error parsing '%s': %s", path, err)
	}

	return c.checkRedirects()
}

func (c *Config) checkRedirects() error {
	for k, v := range c.Redirects {
		if !strings.HasPrefix(v, "http://") {
			c.Redirects[k] = "http://" + v
		}
	}

	return nil
}
