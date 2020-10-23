// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package ds provides methods and interfaces to estimate data size.
package ds

import (
	"fmt"
	"io"
)

// Data must be implemented by any data to estimate.
type Data interface {
	Size() (min, max uint64)
	Kind() string
	fmt.Stringer
}

// Estimator must be implemented by any data size estimator.
type Estimator func(io.Reader, io.Writer) error
