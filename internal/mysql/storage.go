// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package mysql provides methods to parse SQL and estimate MySQL data sizes.
package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rvflash/ds/pkg/ds"
	"github.com/xwb1989/sqlparser"
)

// Storage represents a storage. It can contain many MySQL database.
type Storage []Database

// addDatabase adds a database to the storage.
func (s Storage) addDatabase(name, charset string) (Storage, string) {
	_, err := s.get(name)
	if err != nil {
		return append(s, Database{
			Name:    name,
			Charset: charset,
		}), name
	}
	return s, name
}

// dropDatabase drops a database by its name.
func (s Storage) dropDatabase(name string) Storage {
	i, err := s.get(name)
	if err != nil {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

// createTable tries to create a table inside the given database.
func (s Storage) createTable(dbName string, stmt *sqlparser.DDL) error {
	i, err := s.get(dbName)
	if err != nil {
		return err
	}
	opts := options(stmt.TableSpec)
	t := Table{
		Columns:   columns(stmt.TableSpec, s[i].Charset),
		Engine:    ToEngine(opts[engine]),
		Name:      stmt.NewName.Name.String(),
		RowFormat: ToRowFormat(opts[rowFormat]),
	}
	err = t.addKeys(stmt.TableSpec)
	if err != nil {
		return err
	}
	err = t.Analyze()
	if err != nil {
		return err
	}
	s[i].Tables = append(s[i].Tables, t)
	return nil
}

const (
	engine    = "engine"
	rowFormat = "row_format"

	equal = "="
	pair  = 2
	space = " "
)

// parses table spec (ex: engine=InnoDB default charset=latin1) to build kv options.
func options(spec *sqlparser.TableSpec) map[string]string {
	if spec == nil || spec.Options == "" {
		return nil
	}
	res := make(map[string]string)
	for _, s := range strings.Split(spec.Options, space) {
		kv := strings.SplitN(s, equal, pair)
		if len(kv) == pair {
			res[strings.ToLower(kv[0])] = kv[1]
		}
	}
	return res
}

func columns(spec *sqlparser.TableSpec, dbCharset string) []Column {
	if spec == nil || len(spec.Columns) == 0 {
		return nil
	}
	res := make([]Column, len(spec.Columns))
	for k, v := range spec.Columns {
		c := Column{
			Name:     v.Name.String(),
			Charset:  Charset(v.Type.Charset, dbCharset),
			DataType: ToDataType(v.Type.Type),
			NotNull:  bool(v.Type.NotNull),
		}
		if v.Type.Length != nil {
			c.DataSize, _ = strconv.ParseUint(string(v.Type.Length.Val), base10, bits64)
		}
		res[k] = c
	}
	return res
}

// alterTable tries to alter this database's table.
func (s Storage) alterTable(dbName string, _ *sqlparser.DDL) error {
	_, err := s.get(dbName)
	if err != nil {
		return err
	}
	return nil
}

// dropTable tries to drop this database's table.
func (s Storage) dropTable(dbName string, _ *sqlparser.DDL) error {
	_, err := s.get(dbName)
	if err != nil {
		return err
	}
	return nil
}

// renameTable tries to rename this database's table.
func (s Storage) renameTable(dbName string, _ *sqlparser.DDL) error {
	_, err := s.get(dbName)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) get(name string) (pos int, err error) {
	for i, d := range s {
		if d.Name == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("database: %s: %w", name, ds.ErrInvalid)
}
