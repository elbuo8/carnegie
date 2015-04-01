# Carnegie

[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/elbuo8/carnegie/carnegie) [![Build Status](https://travis-ci.org/elbuo8/carnegie.svg?branch=master)](https://travis-ci.org/elbuo8/carnegie) [![Coverage Status](https://coveralls.io/repos/elbuo8/carnegie/badge.svg?branch=master)](https://coveralls.io/r/elbuo8/carnegie?branch=master)

![carnegie](http://explorepahistory.com/kora/files/1/2/1-2-A0F-25-ExplorePAHistory-a0j4b6-a_349.jpg)

Carnegie is a distributed reverse proxy inspired by [Hipache](https://github.com/hipache/hipache). Currently on *active* development.

## Installing

```bash
go get github.com/elbuo8/carnegie
```

## Running

### 1. Configuration file

Carnegie understands `JSON` & `YAML` thanks to [viper](https://godoc.org/github.com/spf13/viper). For the purposes of this example, `JSON` will be used.

```js
{
  "backend": "consul"
}
```

### 2. Go!

```bash
$ carnegie -c config.json
```

Carnegie will be started on port `8181` and will expect a local installation of `consul`.

## Configuration

The following are fields Carnegie will use when starting up.

* `address`: address that will be listened on (defaults to `:8181`).
* `interval`: frequency that new backends will refreshed (defaults to `1m0s`).
* `cert`: path of certificate file.
* `key`: path of key file.
* `backend`: type of backend to use (`consul` is the only supported one at the moment)

This might change as development progresses.

## Health Checks

Carnegie will remove a VHOST if a backend returns `5xx`. It is recommended to have an active health check system along side Carnegie. Currently, `consul` performs health checks and will only return `healthy` hosts.
