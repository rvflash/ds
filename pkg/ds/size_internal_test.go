// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package ds

import (
	"math"
	"testing"

	"github.com/matryer/is"
)

func TestMin(t *testing.T) {
	are := is.New(t)
	are.Equal(min(math.MinInt8, math.MaxUint8), math.MinInt8) // mismatch left
	are.Equal(min(math.MaxInt8, math.MinInt8), math.MinInt8)  // mismatch right
}
