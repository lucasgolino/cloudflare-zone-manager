# Cloudflare Zone Manager


### Enviroment Variables

```bash
export CONFIG_MAP_PATH=path/to/condig/map.yaml
export CONFIG_MOD_PATH=../modules # Default path to built-in modules
```

### Config map
 This is a config example for CZM to identify your zones and configurations

```yaml
cloudflare:
    email: "your@email.com"
    api_key: "YourCloudflareAPIKey"
zones:
    -   id: "0000000000000000000"
        hostname: "golinux.network"
        dns:
            -   name: "subdomain.golinux.network"
                dtype: "A"
                proxied: false
                ttl: 120
                rules:
                    not-exist: "create"
                    update: "always"
                module:
                    name: "external-ip"
                    metadata:
                        -   key: "route"
                            data: "eno1"
            -   name: "something.golinux.network"
                dtype: "A"
                content: "10.0.0.1"
                proxied: false
                ttl: 120
                rules:
                    not-exist: "create"
                    update: "always"
```

#### Rules

Existent rules are:
 - `not-exist` with keys `"create" | "skip"`
 - `update` with `"always" | "never"`
 
 Update its for when CZM found a diff over DNS on Cloudflare and Config Map
 

#### Modules
Modules its a tools to create DNSRecord content procedurally

Included modules
 - `external-ip` - This module fetch your external ip to create DDNS Service
 
##### How to create a Module

For your module can be compatible they need some things

- Main function named as `Resolve(args interface{})`
- a export named Plugin as type of your module name

Example: **modulename.go**
```go
package main

type modulename string

func (e modulename) Resolve(args interface{}) (string) {
	return "192.168.0.1"
}

var Plugin modulename
```

How to build:
`go build -buildmode=plugin -o modulename.so -o modulename.go`


### Contributing

Feel free to fork this project and send pull request or open a issue.

Thanks!