[![Build Status](https://travis-ci.org/Ashish-Bansal/redis-spectacles.svg?branch=master)](https://travis-ci.org/Ashish-Bansal/redis-spectacles)
[![Apache](https://img.shields.io/badge/license-APACHE-blue.svg)](https://raw.githubusercontent.com/Ashish-Bansal/redis-spectacles/master/LICENSE)

# Redis Spectacles

Redis Spectacles is a CLI tool which scans your whole redis keyspace and helps you visualise key patterns which are most used in your database.


## Build

#### Pre-requisites:
1. GoLang >= 1.13

#### Steps to build:
1. `git clone https://github.com/Ashish-Bansal/redis-spectacles`
2. `go build ./...`

## Usage

To explore redis keyspace, you can run interactive version using
```
./cmd/cmd interactive --url "redis://localhost/0"
```

You explore more available options you can run `./cmd/cmd help`.
