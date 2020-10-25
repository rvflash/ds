// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package ds_test

import (
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"github.com/rvflash/ds/pkg/ds"
	ds_mock "github.com/rvflash/ds/testdata/mock/ds"
)

func TestNewDataSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const a, b, c, d = math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64
	d0 := ds_mock.NewMockData(ctrl)
	d0.EXPECT().Size().Return(uint64(a), uint64(b)).AnyTimes()
	d0.EXPECT().Kind().Return("kind").AnyTimes()
	d0.EXPECT().String().Return("name").AnyTimes()

	var (
		are = is.New(t)
		d1  = ds.NewDataSize(d0, c, d)
	)
	oa, ob := d0.Size()
	na, nb := d1.Size()
	are.True(oa != na)                  // expected different minimum value
	are.True(ob != nb)                  // expected different maximum value
	are.Equal(d0.Kind(), d1.Kind())     // mismatch kind
	are.Equal(d0.String(), d1.String()) // mismatch name
}
