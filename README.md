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

With the same file, we can estimate its size per one, and a million rows.
See the following command line which also prints details per column or key (verbose mode):

```
$ ds mysql -v -n 1000000 testdata/mysql/sample.sql 
+--------------+------------------------+---------------+---------------+-----------------+-----------------+
| DATA         | TYPE                   | PER ROW (MIN) | PER ROW (MAX) | X 1000000 (MIN) | X 1000000 (MAX) |
+--------------+------------------------+---------------+---------------+-----------------+-----------------+
| id           | int                    |        4.00 B |        4.00 B |         4.00 MB |         4.00 MB |
| name         | char(35, utf8mb4)      |       35.00 B |       35.00 B |        35.00 MB |        35.00 MB |
| country_code | char(3, utf8mb4)       |        3.00 B |        3.00 B |         3.00 MB |         3.00 MB |
| district     | char(20, utf8mb4)      |       20.00 B |       20.00 B |        20.00 MB |        20.00 MB |
| population   | int                    |        4.00 B |        4.00 B |         4.00 MB |         4.00 MB |
| PRIMARY      | key(id)                |        4.00 B |        4.00 B |         4.00 MB |         4.00 MB |
| country_code | key(country_code)      |        7.00 B |        7.00 B |         7.00 MB |         7.00 MB |
| city         | table(InnoDB, dynamic) |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
|              |                        |               |               |                 |                 |
| country      | database               |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
+--------------+------------------------+---------------+---------------+-----------------+-----------------+
```

Without the verbose mode, only tables with aggregates by database are displayed.

```
$ cat testdata/mysql/sample.sql | ds mysql -n 1000000 
+---------+------------------------+---------------+---------------+-----------------+-----------------+
| DATA    | TYPE                   | PER ROW (MIN) | PER ROW (MAX) | X 1000000 (MIN) | X 1000000 (MAX) |
+---------+------------------------+---------------+---------------+-----------------+-----------------+
| city    | table(InnoDB, dynamic) |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
| country | database               |       77.00 B |       77.00 B |        77.00 MB |        77.00 MB |
+---------+------------------------+---------------+---------------+-----------------+-----------------+
```


### Features

- Supports MyISAM engine with Static (Fixed-Length), Dynamic and Compressed table characteristics.
- Supports InnoDB engine with Redundant, Compact, Dynamic and Compressed row formats.
- Supports various statements `CREATE DATABASE`, `DROP DATABASE` or `CREATE TABLE`. More incoming!
- The charset is takes account in the computation. 
If no one is defined on the table, the database's charset is used as failover, otherwise `utf8mb4` is used.  
- Display the minimum and maximum sizes estimations to handle variable data types.
- Data sizes are calculated per column, per key, per table and per database.
- Results are aggregated by database. If not specified, `unknown` name is used by default.


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