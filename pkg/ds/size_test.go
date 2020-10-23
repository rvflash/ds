// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package ds_test

import (
	"math"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/ds/pkg/ds"
)

func TestHumanSize(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			in      uint64
			decimal uint8
			out     string
		}{
			{out: "0 B"},
			{in: math.MaxUint16, out: "65 KB"},
			{in: math.MaxUint32, out: "4 GB"},
			{in: math.MaxUint64, decimal: 3, out: "18446744.073 TB"},
			{in: 425005, decimal: 2, out: "425.00 KB"},
			{in: 8741208, decimal: 2, out: "8.74 MB"},
			{in: 114448910, decimal: 2, out: "114.44 MB"},
			{in: 557891, decimal: 2, out: "557.89 KB"},
			{in: 557, decimal: 2, out: "557.00 B"},
			{in: 114710, decimal: 2, out: "114.71 KB"},
			{in: 8933578, decimal: 2, out: "8.93 MB"},
			{in: 5849684981, decimal: 2, out: "5.84 GB"},
			{in: 12033687, decimal: 2, out: "12.03 MB"},
			{in: 742289, decimal: 2, out: "742.28 KB"},
			{in: 678007439, decimal: 2, out: "678.00 MB"},
		}
	)
	for _, tt := range dt {
		tt := tt
		t.Run(tt.out, func(t *testing.T) {
			out := ds.HumanSize(tt.in, tt.decimal)
			are.Equal(tt.out, out) // mismatch result
		})
	}
}
