// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"strings"

	"github.com/rvflash/ds/pkg/ds"
	"github.com/xwb1989/sqlparser"
)

// Table represents a table.
type Table struct {
	Name      string
	Engine    Engine
	Columns   []Column
	Indexes   []Index
	RowFormat RowFormat
}

// Analyze rechallenges any table properties to validate them, to define the primary key or the row format.
func (t *Table) Analyze() error {
	switch {
	case t.Name == "":
		return ds.WrapErr("table name", ds.ErrMissing)
	case t.Engine == "":
		return ds.WrapErr("table engine", ds.ErrInvalid)
	case len(t.Columns) == 0:
		return ds.WrapErr("table column", ds.ErrMissing)
	default:
		t.RowFormat = t.Engine.RowFormat(t.Columns, t.RowFormat)
		return nil
	}
}

const table = "table"

// Kind implements the ds.Data interface.
func (t Table) Kind() string {
	var a []string
	if s := t.Engine.String(); s != "" {
		a = append(a, s)
	}
	if s := t.RowFormat.String(); s != "" {
		a = append(a, s)
	}
	return fmt.Sprintf("%s(%s)", table, strings.Join(a, ", "))
}

// Fields returns columns properties.
func (t Table) Fields() []ds.Data {
	return t.Engine.Fields(t.Columns)
}

// Keys returns keys properties.
func (t Table) Keys() []ds.Data {
	return t.Engine.Keys(t.Indexes, t.primaryKeyIndex())
}

func (t *Table) primaryKeyIndex() int {
	for p, k := range t.Indexes {
		if k.Primary {
			return p
		}
	}
	return notFound
}

// Size implements the ds.Data interface.
func (t Table) Size() (min, max uint64) {
	min, max = t.Engine.RowSize(t.Columns, t.RowFormat)
	var n, x uint64
	for _, k := range t.Keys() {
		n, x = k.Size()
		min += n
		max += x
	}
	return
}

// String implements the ds.Data interface.
func (t Table) String() string {
	return t.Name
}

func (t *Table) addKeys(spec *sqlparser.TableSpec) (err error) {
	if spec == nil || len(spec.Indexes) == 0 {
		return
	}
	for _, k := range spec.Indexes {
		var (
			name string
			cols = make([]string, len(k.Columns))
		)
		for p, c := range k.Columns {
			cols[p] = c.Column.String()
		}
		if k.Info != nil {
			name = k.Info.Name.String()
		}
		err = t.addKey(name, cols, primary(k.Info))
		if err != nil {
			return
		}
	}
	return
}

func primary(info *sqlparser.IndexInfo) bool {
	if info == nil {
		return false
	}
	return info.Primary
}

func (t *Table) addKey(name string, columns []string, primary bool) error {
	cols := t.columnsNamed(columns)
	if len(cols) == 0 {
		return ds.WrapErr("key column", ds.ErrInvalid)
	}
	t.Indexes = append(t.Indexes, Index{
		Name:    name,
		Columns: cols,
		Primary: primary,
	})
	return nil
}

// notFound is the index value returned if the data is not found.
const notFound = -1

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
