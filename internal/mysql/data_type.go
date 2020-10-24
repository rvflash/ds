// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"math"
	"strings"
)

// ToDataType returns a MySQL DataType based on the given data name.
func ToDataType(s string) DataType {
	return DataType(strings.ToLower(s))
}

// DataType represents a MySQL data type.
type DataType string

// List of supported MySQL data types.
const (
	Bit        DataType = "bit"
	TinyInt    DataType = "tinyint"
	SmallInt   DataType = "smallint"
	MediumInt  DataType = "mediumint"
	Int        DataType = "int"
	Integer    DataType = "integer"
	BigInt     DataType = "bigint"
	Float      DataType = "float"
	Double     DataType = "double"
	Decimal    DataType = "decimal"
	Numeric    DataType = "numeric"
	Real       DataType = "real"
	Year       DataType = "year"
	Date       DataType = "date"
	Time       DataType = "time"
	Timestamp  DataType = "timestamp"
	DateTime   DataType = "datetime"
	Char       DataType = "char"
	Binary     DataType = "binary"
	VarChar    DataType = "varchar"
	VarBinary  DataType = "varbinary"
	TinyBlob   DataType = "tinyblob"
	TinyText   DataType = "tinytext"
	Blob       DataType = "blob"
	Text       DataType = "test"
	MediumBlob DataType = "mediumblob"
	MediumText DataType = "mediumtext"
	LongBlob   DataType = "longblob"
	LongText   DataType = "longtext"
	JSON       DataType = "json"
	Enum       DataType = "enum"
	Set        DataType = "set"
)

// Kind implements the ds.Data interface.
func (DataType) Kind() string {
	return ""
}

// IsInt returns true if the data type is an integer.
func (d DataType) IsInt() bool {
	switch d {
	case
		TinyInt, SmallInt, MediumInt, Int, BigInt:
		return true
	default:
		return false
	}
}

// IsString returns true if the data type is a string.
func (d DataType) IsString() bool {
	switch d {
	case
		Char, Binary,
		VarChar, VarBinary,
		TinyBlob, MediumBlob, Blob, LongBlob,
		TinyText, MediumText, Text, LongText,
		Enum, Set,
		JSON:
		return true
	default:
		return false
	}
}

// IsVar returns true if the data type is a variable one.
func (d DataType) IsVar() bool {
	switch d {
	case
		VarChar, VarBinary,
		TinyBlob, MediumBlob, Blob, LongBlob,
		TinyText, MediumText, Text, LongText,
		JSON:
		return true
	default:
		return false
	}
}

const maxMediumSize = 16777216

// Size returns the required storage of the data type for this requested size in bytes and charset.
// It implements the ds.Data interface.
// todo To improve (decimal, numeric, etc.)
// See https://dev.mysql.com/doc/refman/8.0/en/storage-requirements.html
func (d DataType) Size(size uint64, charset string) (min, max uint64) {
	switch d {
	case Bit:
		return both((size + 7) / 8)
	case Binary, Char:
		return both(size)
	case TinyInt, Year:
		return both(1)
	case SmallInt:
		return both(2)
	case MediumInt, Date:
		return both(3)
	case Int, Integer:
		return both(4)
	case BigInt:
		return both(8)
	case Float:
		return float(size)
	case Double, Real:
		return both(8)
	case Decimal, Numeric:
		return 4, 8
	case Time:
		return both(3 + fsp(size))
	case Timestamp:
		return both(4 + fsp(size))
	case DateTime:
		return both(5 + fsp(size))
	case VarChar:
		return variable(bytes(size, charset))
	case VarBinary:
		return variable(size)
	case TinyBlob, TinyText:
		return blob(bytes(size, charset), 1, math.MaxUint8)
	case Blob, Text:
		return blob(bytes(size, charset), 2, math.MaxUint16)
	case MediumBlob, MediumText:
		return blob(bytes(size, charset), 3, maxMediumSize)
	case LongBlob, LongText, JSON:
		return blob(bytes(size, charset), 4, math.MaxUint32)
	case Enum:
		return enum(size)
	case Set:
		return set(size)
	default:
		return 0, math.MaxUint64
	}
}

// String implements the ds.Data interface.
func (d DataType) String() string {
	return string(d)
}

func blob(size, reserved, max uint64) (uint64, uint64) {
	if size > 0 {
		return reserved, size + reserved
	}
	return reserved, max - 1 + reserved
}

func both(i uint64) (uint64, uint64) {
	return i, i
}

func bytes(size uint64, charset string) uint64 {
	char, set := charsets[charset]
	if !set {
		return 0
	}
	return size * uint64(char)
}

func enum(size uint64) (uint64, uint64) {
	if size > math.MaxUint8 {
		return both(2)
	}
	return both(1)
}

func float(size uint64) (uint64, uint64) {
	if size > 24 {
		return both(8)
	}
	return both(4)
}

// fsp aka fractional seconds precision.
func fsp(size uint64) uint64 {
	switch {
	case size > 4:
		return 3
	case size > 2:
		return 2
	case size > 0:
		return 1
	default:
		return 0
	}
}

func set(size uint64) (uint64, uint64) {
	if size > 64 {
		return both(8)
	}
	if size == 0 {
		return both(1)
	}
	return both((size + 7) / 8)
}

func variable(size uint64) (uint64, uint64) {
	if size > math.MaxUint8 {
		return 2, size + 2
	}
	if size == 0 {
		return 1, math.MaxUint8
	}
	return 1, size + 1
}
