# HTTP demuxer
Configurable HTTP traffic forwarder written in Go.

This application works as HTTP reverse proxy. It aggregates traffic from several sub/domains and forwards it to specified host and port. Main purpose is to route traffic to different services running on the same server.

## Configuration
Routing is based on rules specified in JSON file. Example configuration can look like this:
```javascript
{
    "port": 80,
    "forwards": {
            "www.example.com": ":8000",
            "blog.example.com": ":8001",
            "service.example.com": ":8002"
    }
}
```
With this config, demuxer will listen on port 80 a forward all HTTP requests to `www.example.com` to `localhost:8000`,
`blog.example.com` to `localhost:8001`, and `service.example.com` to `localhost:8002`.

## Usage
```bash
./forwarder path_to_config.json
```
