# Block User-Agent

[![Build Status](https://github.com/agence-gaya/traefik-plugin-blockuseragent/workflows/Main/badge.svg?branch=master)](https://github.com/agence-gaya/traefik-plugin-blockuseragent/actions)

Block User-Agent is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which sends an HTTP `403 Forbidden` 
response when the requested HTTP User-Agent header matches one the configured [regular expressions](https://github.com/google/re2/wiki/Syntax).

## Configuration

## StaticUpdate 

```toml
[pilot]
    token="xxx"

[experimental.plugins.blockuseragent]
    modulename = "github.com/agence-gaya/traefik-plugin-blockuseragent"
    version = "vX.Y.Z"
```

## Dynamic

To configure the `Block User-Agent` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `blockuseragent` middleware plugin to block all HTTP requests with a User-Agent like `\bTheAgent\b`.
You can use regexAllow to make exception on blocking regex.

```toml
[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["block-foo"]
    service = "my-service"

# Block all user agent containing TheAgent except if containing Allowed word
[http.middlewares]
  [http.middlewares.block-foo.plugin.blockuseragent]
    regexAllow = ["\bAllowed\b"]    
    regex = ["\bTheAgent\b"]

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```
