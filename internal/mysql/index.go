// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"strings"
)

// Index is a table's key.
type Index struct {
	Name    string
	Columns []Column
	Primary bool
}

// Size implements the ds.Data interface.
func (i Index) Size() (min, max uint64) {
	var n, x uint64
	for _, c := range i.Columns {
		n, x = c.Size()
		min += n
		max += x
	}
	return
}

// Kind implements the ds.Data interface.
func (i Index) Kind() string {
	names := make([]string, len(i.Columns))
	for p, v := range i.Columns {
		names[p] = v.String()
	}
	if len(names) == 0 {
		return ""
	}
	return fmt.Sprintf("key(%s)", strings.Join(names, ", "))
}

// String implements the ds.Data interface.
func (i Index) String() string {
	return i.Name
}
