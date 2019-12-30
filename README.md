# sql-prom-stats [![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/sql-prom-stats)](https://goreportcard.com/report/github.com/saromanov/sql-prom-stats)
[![Coverage Status](https://coveralls.io/repos/github/saromanov/sql-prom-stats/badge.svg?branch=master)](https://coveralls.io/github/saromanov/sql-prom-stats?branch=master)
[![GoDoc](https://godoc.org/github.com/saromanov/sql-prom-stats?status.svg)](https://godoc.org/github.com/saromanov/sql-prom-stats) 

Collecting metrics from SQL to Prometheus

## Install

```sh
go get github.com/saromanov/sql-prom-stats
```

## Usage
This tool using standard list of metrics which exposing from [DB.Stats](https://golang.org/pkg/database/sql/#DB.Stats) and then pushing to Prometheus

See `example` directory

## Licence
MIT


