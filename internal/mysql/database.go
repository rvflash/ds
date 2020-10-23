// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"strings"

	"github.com/rvflash/ds/pkg/ds"
)

// DefaultDatabaseName is the name used by default for a database.
const DefaultDatabaseName = "unknown"

// notFound is the index value returned if the data is not found.
const notFound = -1

// Database represents a database.
type Database struct {
	Name    string
	Charset string
	Tables  []Table
}

// Kind implements the ds.Data interface.
func (Database) Kind() string {
	return "DATABASE"
}

// Size implements the ds.Data interface.
func (d Database) Size() (min, max uint64) {
	var n, x uint64
	for _, c := range d.Tables {
		n, x = c.Size()
		min += n
		max += x
	}
	return
}

// String implements the ds.Data interface.
func (d Database) String() string {
	return d.Name
}

// Table represents a table.
type Table struct {
	Name    string
	Columns []Column
	Keys    []Key
}

// AddKey tries to add this table's key.
func (t *Table) AddKey(name string, columns []string) error {
	cols := t.columnsNamed(columns)
	if len(cols) == 0 {
		return fmt.Errorf("key's columns: %w", ds.ErrInvalid)
	}
	t.Keys = append(t.Keys, Key{
		Name:    name,
		Columns: cols,
	})
	return nil
}

// Kind implements the ds.Data interface.
func (Table) Kind() string {
	return "TABLE"
}

// Size implements the ds.Data interface.
func (t Table) Size() (min, max uint64) {
	var n, x uint64
	for _, c := range t.Columns {
		n, x = c.Size()
		min += n
		max += x
	}
	for _, c := range t.Keys {
		n, x = c.Size()
		min += n
		max += x
	}
	return
}
func (t *Table) columnIndex(name string) int {
	for i, c := range t.Columns {
		if c.Name == name {
			return i
		}
	}
	return notFound
}

func (t *Table) columnsNamed(names []string) []Column {
	if len(names) == 0 {
		return nil
	}
	var res []Column
	for _, name := range names {
		i := t.columnIndex(name)
		if i != notFound {
			res = append(res, t.Columns[i])
		}
	}
	if len(res) != len(names) {
		return nil
	}
	return res
}

// String implements the ds.Data interface.
func (t Table) String() string {
	return t.Name
}

// Column is a table's column.
type Column struct {
	Name     string
	Charset  string
	DataSize uint64
	DataType DataType
}

// Size implements the ds.Data interface.
func (c Column) Size() (min, max uint64) {
	return c.DataType.Size(c.DataSize, c.Charset)
}

// Kind implements the ds.Data interface.
func (c Column) Kind() string {
	if c.DataSize == 0 {
		return c.DataType.String()
	}
	return fmt.Sprintf("%s(%d)", c.DataType.String(), c.DataSize)
}

// String implements the ds.Data interface.
func (c Column) String() string {
	return c.Name
}

// Key is a table's key.
type Key struct {
	Name    string
	Columns []Column
}

// Size implements the ds.Data interface.
func (k Key) Size() (min, max uint64) {
	var n, x uint64
	for _, c := range k.Columns {
		n, x = c.Size()
		min += n
		max += x
	}
	return
}

// Kind implements the ds.Data interface.
func (k Key) Kind() string {
	names := make([]string, len(k.Columns))
	for k, v := range k.Columns {
		names[k] = v.String()
	}
	if len(names) == 0 {
		return ""
	}
	return fmt.Sprintf("KEY(%s)", strings.Join(names, ", "))
}

// String implements the ds.Data interface.
func (k Key) String() string {
	return k.Name
}
