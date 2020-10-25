// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"strings"

	"github.com/rvflash/ds/pkg/ds"
)

// Column is a table's column.
type Column struct {
	Name     string
	Charset  string
	DataSize uint64
	DataType DataType
	NotNull  bool
}

// Size implements the ds.Data interface.
func (c Column) Size() (min, max uint64) {
	return c.DataType.Size(c.DataSize, c.Charset)
}

// Kind implements the ds.Data interface.
func (c Column) Kind() string {
	if c.DataType.IsInt() {
		return c.DataType.String()
	}
	var a []string
	if c.DataSize > 0 {
		a = append(a, ds.Unit(c.DataSize).String())
	}
	if c.DataType.IsString() {
		a = append(a, c.Charset)
	}
	if s := strings.Join(a, ", "); s != "" {
		return fmt.Sprintf("%s(%s)", c.DataType.String(), s)
	}
	return c.DataType.String()
}

// String implements the ds.Data interface.
func (c Column) String() string {
	return c.Name
}
