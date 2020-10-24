// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package ds

import "fmt"

// Error represents an error.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

// List of known errors.
const (
	// ErrProcess is returned when the estimation process has failed.
	ErrProcess = Error("invalid estimator")
	// ErrInvalid is returned  when the data is invalid.
	ErrInvalid = Error("invalid data")
	// ErrMissing is returned  when the data is missing.
	ErrMissing = Error("missing data")
)

// WrapErr wraps the error with the given message.
func WrapErr(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
