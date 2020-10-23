# ds

[![GoDoc](https://godoc.org/github.com/rvflash/ds?status.svg)](https://godoc.org/github.com/rvflash/ds)
[![Build Status](https://api.travis-ci.org/rvflash/ds.svg?branch=master)](https://travis-ci.org/rvflash/ds?branch=master)
[![Code Coverage](https://codecov.io/gh/rvflash/ds/branch/master/graph/badge.svg)](https://codecov.io/gh/rvflash/ds)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/ds?)](https://goreportcard.com/report/github.com/rvflash/ds)


`ds` estimates the data size in bytes.

## Features

1. ...


## Demo with mysql as data source

```shell
$ 
```

## Installation

### Go

```shell
GO111MODULE=on go get github.com/rvflash/ds@v0.1.0
```

### Docker

```shell
docker run --rm -v $(pwd):/pkg rvflash/ds:v0.1.0 ds
```

## Usage

### MySQL data size

```shell
ds mysql [flags] [file]
```

or by parsing the standard input:

```shell
cat database.sql | ds mysql -v
```

The `ds mysql` sub command estimates the data size of a database or one or more tables, based on a MySQL statements.
It supports the following flags:

* `-v`: verbose output