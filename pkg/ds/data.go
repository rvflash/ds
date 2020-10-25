// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package ds provides methods and interfaces to estimate data size.
package ds

import (
	"fmt"
	"io"
)

//go:generate mockgen -source ${GOFILE} -destination ../../testdata/mock/${GOPACKAGE}/${GOFILE}

// NewDataSize allows to overload the minimum and maximum data size.
func NewDataSize(d Data, min, max uint64) Data {
	return &data{Data: d, min: min, max: max}
}

// Data must be implemented by any data to estimate.
type Data interface {
	Size() (min, max uint64)
	Kind() string
	fmt.Stringer
}

// Estimator must be implemented by any data size estimator.
type Estimator func(io.Reader, io.Writer) error

type data struct {
	Data
	min, max uint64
}

// Size overloads the Size method to force new min and max sizes.
func (d data) Size() (min, max uint64) {
	return d.min, d.min
}
