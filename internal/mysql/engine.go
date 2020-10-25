// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"strings"

	"github.com/rvflash/ds/pkg/ds"
)

// ToEngine returns a MySQL engine based on the given engine name.
// By default the default MySQL engine is used if no one is provided.
func ToEngine(s string) Engine {
	switch strings.ToLower(s) {
	case "", strings.ToLower(InnoDB.String()):
		return InnoDB
	case strings.ToLower(MyISAM.String()):
		return MyISAM
	default:
		return ""
	}
}

// Engine is an engine.
type Engine string

// List of supported engines.
const (
	InnoDB = Engine("InnoDB")
	MyISAM = Engine("MyISAM")
)

// Fields returns the columns with their sizes updated with engine and row format constrains.
func (e Engine) Fields(cols []Column) []ds.Data {
	res := make([]ds.Data, len(cols))
	for p, c := range cols {
		res[p] = c
	}
	return res
}

// Keys returns the keys with their sizes updated with engine constrains.
func (e Engine) Keys(indexes []Index, primary int) []ds.Data {
	switch e {
	case InnoDB:
		return innoDBKeys(indexes, primary)
	case MyISAM:
		return myISAMKeys(indexes)
	default:
		return nil
	}
}

const defaultClusteredIndexSize = 6

// https://dev.mysql.com/doc/refman/8.0/en/innodb-index-types.html
func innoDBKeys(keys []Index, primary int) []ds.Data {
	var (
		n, x uint64
		res  = make([]ds.Data, len(keys))
		pks  = func() (min, max uint64) {
			if primary != notFound && len(keys) > primary {
				return keys[primary].Size()
			}
			return both(defaultClusteredIndexSize)
		}
	)
	pkn, pkx := pks()
	for p, k := range keys {
		n, x = k.Size()
		if k.Primary {
			res[p] = k
		} else {
			res[p] = ds.NewDataSize(k, n+pkn, x+pkx)
		}
	}
	return res
}

// https://dev.mysql.com/doc/refman/8.0/en/key-space.html
func myISAMKeys(keys []Index) []ds.Data {
	var (
		n, x uint64
		res  = make([]ds.Data, len(keys))
		fml  = func(size uint64) uint64 {
			i := (float64(size) + 4) / 0.67
			return uint64(i)
		}
	)
	for p, c := range keys {
		n, x = c.Size()
		res[p] = ds.NewDataSize(c, fml(n), fml(x))
	}
	return res
}

// RowFormat defines the row format to use based on columns or the current value.
// sql: SELECT row_format FROM information_schema.tables WHERE table_schema="dbName" AND table_name="tbName";
func (e Engine) RowFormat(cols []Column, cur RowFormat) RowFormat {
	switch e {
	case InnoDB:
		return innoDBRowFormat(cur)
	case MyISAM:
		return myISAMRowFormat(cols, cur)
	default:
		return UnknownRowFormat
	}
}

func innoDBRowFormat(cur RowFormat) RowFormat {
	if cur != UnknownRowFormat {
		return cur
	}
	return DynamicRowFormat
}

func myISAMRowFormat(cols []Column, cur RowFormat) RowFormat {
	if cur == CompressedRowFormat {
		return cur
	}
	for _, c := range cols {
		if c.DataType.IsVar() {
			return DynamicRowFormat
		}
	}
	return StaticRowFormat
}

// RowSize returns the estimates row length.
func (e Engine) RowSize(cols []Column, cur RowFormat) (min, max uint64) {
	switch e {
	case InnoDB:
		return innoDBRowSize(cols, cur)
	case MyISAM:
		return myISAMRowSize(cols, cur)
	default:
		return both(0)
	}
}

func innoDBRowSize(cols []Column, _ RowFormat) (min, max uint64) {
	var n, x uint64
	for _, c := range cols {
		n, x = c.Size()
		min += n
		max += x
	}
	return
}

const (
	staticDeleteFlag = 1
	staticHeader     = 1
	dynamicHeader    = 3
)

func myISAMRowSize(cols []Column, cur RowFormat) (min, max uint64) {
	var n, x, nn, nv, ns, ni uint64
	for _, c := range cols {
		n, x = c.Size()
		min += n
		max += x
		if !c.NotNull {
			nn++
		}
		if c.DataType.IsVar() {
			nv++
		}
		if c.DataType.IsString() {
			ns++
		}
		if c.DataType.IsInt() {
			ni++
		}
	}
	num := uint64(len(cols))
	switch cur {
	case StaticRowFormat:
		// Variable data type are space-padded to the specified column width.
		// Formula:
		// 1 as header
		// + (sum of column lengths)
		// + (number of NULL columns + delete_flag + 7)/8
		// + (number of variable-length columns)
		return both(staticHeader + max + ((nn + staticDeleteFlag + 7) / 8) + nv)
	case DynamicRowFormat:
		// Formula:
		// 3 as header
		// + (number of columns + 7) / 8
		// + (number of char columns)
		// + (packed size of numeric columns)
		// + (length of strings)
		// + (number of NULL columns + 7) / 8
		var fml = func(size uint64) uint64 {
			return dynamicHeader + ((num + 7) / 8) + ns + size + ((nn + 7) / 8)
		}
		return fml(min), fml(max)
	case CompressedRowFormat:
		// Formula:
		// 1 or 3 as header
		// + data size
		return staticHeader + min, dynamicHeader + max
	default:
		return both(0)
	}
}

// String implements the fmt.Stringer interface.
func (e Engine) String() string {
	return string(e)
}
