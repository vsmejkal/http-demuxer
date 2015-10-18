package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "net/http"
    "net/http/httputil"
    "encoding/json"
)

type Config struct {
    Port        int
    Forwards    map[string]string
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

    // Set up HTTP proxy
    proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
        host, _, _ := net.SplitHostPort(req.Host)
        gate := conf.Forwards[host]

        req.URL.Scheme = "http"
        req.URL.Host = gate
    }}

    // Start server
    addr := fmt.Sprintf(":%d", conf.Port)
    log.Fatal(http.ListenAndServe(addr, proxy))
}