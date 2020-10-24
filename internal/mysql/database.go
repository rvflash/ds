// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

// Database represents a database.
type Database struct {
	Name    string
	Charset string
	Tables  []Table
}

const db = "database"

// Kind implements the ds.Data interface.
func (d Database) Kind() string {
	return db
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
