// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package ds_test

import (
	"testing"

	"github.com/matryer/is"

	"github.com/rvflash/ds/pkg/ds"
)

func TestError_Error(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			err ds.Error
			msg string
		}{
			"Default": {},
			"Invalid": {err: ds.ErrInvalid, msg: "invalid data"},
			"Missing": {err: ds.ErrMissing, msg: "missing data"},
			"Process": {err: ds.ErrProcess, msg: "invalid estimator"},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			are.Equal(tt.err.Error(), tt.msg) // mismatch message
		})
	}
}
