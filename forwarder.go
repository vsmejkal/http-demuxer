package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "net/http"
    "net/http/httputil"
	"strings"
	"encoding/json"
)


type Config struct {
	Port        int
	Forwards    map[string]string
	Redirects 	map[string]string
}

func (c *Config) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open config file '%s'", path)
	}

	decoder := json.NewDecoder(file)
	if err:= decoder.Decode(c); err != nil {
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

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: %s config.json\n", os.Args[0])
    }

    // Load forward rules from file
    var conf Config
    if err := conf.Load(os.Args[1]); err != nil {
        log.Fatal(err)
    }

    // HTTP proxy
    proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
		host, _, _ := net.SplitHostPort(req.Host)
		req.URL.Scheme = "http"
		req.URL.Host = conf.Forwards[host]
	}}

	// HTTP handler
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		host, _, _ := net.SplitHostPort(req.Host)
		forward := conf.Forwards[host]
		redirect := conf.Redirects[host]

		if forward != "" {
			proxy.ServeHTTP(rw, req)
		} else if redirect != "" {
			http.Redirect(rw, req, redirect, http.StatusMovedPermanently)
		} else {
			http.NotFound(rw, req)
		}
	})

    // Start server
    addr := fmt.Sprintf(":%d", conf.Port)
    log.Fatal(http.ListenAndServe(addr, nil))
}