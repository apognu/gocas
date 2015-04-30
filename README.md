# GoCAS

Minimalist CAS server in Go. Here what currently works:

* Basic workflow (/login, /validate, /serviceValidate)
* Trust authentication (disabled by default)
* /logout (no SLO for now)
* simple whitelisting of exact service hosts (wildcard might come some day)
* REST API
* CAS proxy simple implementation
* Middleware system
 * Failed login attempts throttling (quite native for now)

GoCAS requires a MongoDB service to be available. The available authenticators are :

* Dummy (username should equal password, for testing purposes)
* LDAP

Also, the following server protocols are supported:

* CAS (duh!)
* OAuth2

## Configuration

Exhaustive example of configuration can be found in _gocas.yaml.example_. Location of configuration file can be given with switch _-config_.

## Build and run

```
$ cd $GOPATH
$ go get -u github.com/apognu/gocas
$ go install github.com/apognu/gocas
$ $GOPATH/bin/gocas [-config /etc/gocas.yaml]
```

For now, the _template/_ directory must be copied/symlinked in the same directory the binary is located. This might change in the future.

## CAS specification

This is a work in progress, we might or might not cover the whole CAS specification, for now, here is what we do:

*TODO:* include specification coverage stats.
