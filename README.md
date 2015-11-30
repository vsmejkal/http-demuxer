# HTTP Demuxer
Configurable HTTP traffic forwarder written in Go.

This application works as HTTP reverse proxy. It aggregates traffic from several sub/domains and forwards it to specified host and port. Main purpose is to route traffic to different services running on the same server.

## Configuration
Routing is based on the rules specified in JSON file:
```json
{
    "port": 80,

    "forwards": {
        "www.example.com": ":8000",
        "blog.example.com": ":8001",
        "service.example.com": ":8002"
    },

    "redirects": {
        "example.com": "www.example.com"
    }
}
```
With this config, demuxer will listen  on port 80 and forward all requests
to `www.example.com` to local port `8000`, `blog.example.com` to `8001`,
and `service.example.com` to `8002`. It will also redirect all requests to `example.com` to
`www.example.com`.

If demuxer receives request not specified by any _forward_ or _redirect_, it returns 404
Not Found error.

## Usage
```bash
./http-demuxer routes.json
```
