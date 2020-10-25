// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import "strings"

// ToRowFormat returns a row format.
func ToRowFormat(s string) RowFormat {
	return RowFormat(strings.ToLower(s))
}

// RowFormat represents a row format.
type RowFormat string

// List of known row formats
const (
	UnknownRowFormat    = RowFormat("")
	CompactRowFormat    = RowFormat("compact")
	CompressedRowFormat = RowFormat("compressed")
	DynamicRowFormat    = RowFormat("dynamic")
	RedundantRowFormat  = RowFormat("redundant")
	StaticRowFormat     = RowFormat("static")
)

// String implements the fmt.Stringer interface.
func (f RowFormat) String() string {
	return string(f)
}
