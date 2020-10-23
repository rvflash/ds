// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package ds

import (
	"fmt"
	"strconv"
)

// List of supported size ranges.
const (
	Byte     Unit = 1
	KiloByte      = 1000 * Byte
	MegaByte      = KiloByte * KiloByte
	GigaByte      = MegaByte * KiloByte
	TeraByte      = GigaByte * KiloByte
)

// HumanSize converts the given number of bytes into a human size in bytes.
func HumanSize(size uint64, decimal uint8) string {
	return Unit(size).Format(decimal)
}

// Unit is the unit of measure.
type Unit uint64

// Format returns a human size readable of the Unit.
func (u Unit) Format(decimal uint8) string {
	r, p := interval(u)
	if r < KiloByte {
		if decimal == 0 {
			return u.String() + " " + p
		}
		return fmt.Sprintf("%s.%0*d %s", u.String(), decimal, 0, p)
	}
	i := u / r
	if decimal == 0 {
		return i.String() + " " + p
	}
	return fmt.Sprintf("%s.%s %s", i.String(), fractional(u-(i*r), int(decimal)), p)
}

const base10 = 10

// String returns the Unit in string.
func (u Unit) String() string {
	return strconv.FormatUint(uint64(u), base10)
}

func interval(u Unit) (Unit, string) {
	switch {
	case u > TeraByte:
		return TeraByte, "TB"
	case u > GigaByte:
		return GigaByte, "GB"
	case u > MegaByte:
		return MegaByte, "MB"
	case u > KiloByte:
		return KiloByte, "KB"
	default:
		return 1, "B"
	}
}

func fractional(u Unit, n int) string {
	var size int
	switch {
	case u > GigaByte:
		size = 12
	case u > MegaByte:
		size = 9
	case u > KiloByte:
		size = 6
	default:
		size = 3
	}
	s := fmt.Sprintf("%0*d", size, u)
	return s[:min(n, len(s))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
