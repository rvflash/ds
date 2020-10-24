# ds

[![GoDoc](https://godoc.org/github.com/rvflash/ds?status.svg)](https://godoc.org/github.com/rvflash/ds)
[![Build Status](https://api.travis-ci.org/rvflash/ds.svg?branch=master)](https://travis-ci.org/rvflash/ds?branch=master)
[![Code Coverage](https://codecov.io/gh/rvflash/ds/branch/master/graph/badge.svg)](https://codecov.io/gh/rvflash/ds)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/ds?)](https://goreportcard.com/report/github.com/rvflash/ds)


`ds` estimates the data size in bytes.

## ds mysql

The subcommand `mysql` allows estimating data size of databases, tables, by columns or keys, based on SQL statements.

`ds` estimates data sizes by parsing the SQL statements from the named files (or standard input) and 
prints the result in ASCII or CVS format, if the batch mode is enabled. 

For example, based on the SQL statements in the [sample.sql](testdata/mysql/sample.sql) file, 
which generates the following table:

```
mysql> desc city;
+--------------+----------+------+-----+---------+----------------+
| Field        | Type     | Null | Key | Default | Extra          |
+--------------+----------+------+-----+---------+----------------+
| id           | int(11)  | NO   | PRI | NULL    | auto_increment |
| name         | char(35) | NO   |     |         |                |
| country_code | char(3)  | NO   | MUL |         |                |
| district     | char(20) | NO   |     |         |                |
| population   | int(11)  | NO   |     | 0       |                |
+--------------+----------+------+-----+---------+----------------+
5 rows in set (0.00 sec)
 ```

With the same file, we can estimate its size per one, and a million rows and print details per column or key:

```
$ ds mysql -v -n 1000000 testdata/mysql/sample.sql 
+--------------+-------------------+---------------+---------------+-----------------+-----------------+
| DATA         | TYPE              | PER ROW (MIN) | PER ROW (MAX) | X 1000000 (MIN) | X 1000000 (MAX) |
+--------------+-------------------+---------------+---------------+-----------------+-----------------+
| id           | INT(11)           |        4.00 B |        4.00 B |         4.00 MB |         4.00 MB |
| name         | CHAR(35)          |       35.00 B |       35.00 B |        35.00 MB |        35.00 MB |
| country_code | CHAR(3)           |        3.00 B |        3.00 B |         3.00 MB |         3.00 MB |
| district     | CHAR(20)          |       20.00 B |       20.00 B |        20.00 MB |        20.00 MB |
| population   | INT(11)           |        4.00 B |        4.00 B |         4.00 MB |         4.00 MB |
| PRIMARY      | KEY(id)           |        4.00 B |        4.00 B |         4.00 MB |         4.00 MB |
| country_code | KEY(country_code) |        7.00 B |        7.00 B |         7.00 MB |         7.00 MB |
| city         | TABLE             |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
| country      | DATABASE          |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
+--------------+-------------------+---------------+---------------+-----------------+-----------------+
```

Without the verbose mode, only tables with aggregates by database are displayed.

```
$ cat testdata/mysql/sample.sql | ds mysql -n 1000000 
+---------+----------+---------------+---------------+-----------------+-----------------+
| DATA    | TYPE     | PER ROW (MIN) | PER ROW (MAX) | X 1000000 (MIN) | X 1000000 (MAX) |
+---------+----------+---------------+---------------+-----------------+-----------------+
| city    | TABLE    |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
| country | DATABASE |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
+---------+----------+---------------+---------------+-----------------+-----------------+
```


### Features

- Supports the MySQL and InnoDB engines specificities.
- Supports various statements `CREATE DATABASE`, `DROP DATABASE` or `CREATE TABLE`.
- The charset is takes account in the computation. 
If no one is defined on the table, the database's charset is used as failover, otherwise `utf8mb4` is used.  
- Display the minimum and maximum sizes estimations to handle variable data types.
- Data sizes are calculated per column, per key, per table and per database.
- Even if not specified, results are aggregated by database, named `unknown` by default.


### Usage

By default, it uses ASCII table to display the report.
Size per row, and for 100 rows are calculated and displayed using 2 decimals.

```
ds mysql [flags] [file.sql]
```

It supports the following flags:

* `-B`: batch mode, print results using comma as the column separator, with each row on a new line.
* `-n`: number of lines to considerate by table (default 100).
* `-p`: number of decimals to display (default 2).
* `-v`: verbose output, produce more output about what the program does.


## Installation

### Go

```shell
GO111MODULE=on go get github.com/rvflash/ds@v0.1.0
```

### Docker

```shell
docker run --rm -v $(pwd):/pkg rvflash/ds:v0.1.0 ds
```
