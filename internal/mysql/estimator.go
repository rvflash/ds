// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/rvflash/ds/pkg/ds"
)

// Columns names.
const (
	dataName = "Data"
	dataType = "Type"
	minRow   = "Per row (min)"
	maxRow   = "Per row (max)"
)

// Default values used to configure the estimator.
const (
	Command          = "mysql"
	DefaultPerN      = 100
	DefaultPrecision = 2
)

// Config lists any customizable settings.
type Config struct {
	Batch,
	Verbose bool
	Precision,
	PerN uint64
}

// Configurator is implemented by any method exposing cursor to adjust the estimator.
type Configurator func(*Estimator) error

// SetPerN defines the number of data to take account in the estimation.
func SetPerN(i uint64) Configurator {
	return func(e *Estimator) error {
		if i == 0 {
			return ds.WrapErr("per N value", ds.ErrMissing)
		}
		e.perN = i
		return nil
	}
}

// SetPrecision defines the decimal precision used to print data size.
func SetPrecision(i uint64) Configurator {
	return func(e *Estimator) error {
		if i > math.MaxUint8 {
			return ds.WrapErr("precision", ds.ErrInvalid)
		}
		e.precision = uint8(i)
		return nil
	}
}

// SetVerbose defines the verbose mode to use to print the report.
func SetVerbose(verbose bool) Configurator {
	return func(e *Estimator) error {
		e.verbose = verbose
		return nil
	}
}

// SetBatchMode defines if the batch mode must be used to export the report in CSV format.
func SetBatchMode(enabled bool) Configurator {
	return func(e *Estimator) error {
		e.batch = enabled
		return nil
	}
}

// Estimate tries to instantiate a new estimator based on this configuration.
func Estimate(opts ...Configurator) (*Estimator, error) {
	opts = append([]Configurator{
		SetPerN(DefaultPerN),
		SetPrecision(DefaultPrecision),
	}, opts...)
	cnf := new(Estimator)
	for _, opt := range opts {
		err := opt(cnf)
		if err != nil {
			return nil, err
		}
	}
	return cnf, nil
}

// Estimator represents an MySQL data estimator.
type Estimator struct {
	batch,
	verbose bool
	precision uint8
	perN      uint64
}

// Run runs the estimator.
func (e *Estimator) Run(r io.Reader, w io.Writer) error {
	if e.perN == 0 {
		return ds.ErrProcess
	}
	dbs, err := Parse(r)
	if err != nil {
		return err
	}
	if len(dbs) == 0 {
		return ds.ErrMissing
	}
	res := make([][]string, 0)
	for p, d := range dbs {
		if p > 0 && e.verbose {
			res = append(res, e.blank())
		}
		for _, t := range d.Tables {
			if e.verbose {
				for _, c := range t.Fields() {
					res = append(res, e.row(c))
				}
				for _, k := range t.Keys() {
					res = append(res, e.row(k))
				}
			}
			res = append(res, e.row(t))
			if e.verbose {
				res = append(res, e.blank())
			}
		}
		res = append(res, e.row(d))
	}
	if e.batch {
		return e.batchRender(w, res)
	}
	return e.render(w, res)
}

// batchRender prints the results using comma as the column separator.
func (e *Estimator) batchRender(writer io.Writer, data [][]string) error {
	var (
		w   = csv.NewWriter(writer)
		err = w.Write(e.header())
	)
	if err != nil {
		return err
	}
	err = w.WriteAll(data)
	if err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func (e *Estimator) blank() []string {
	return make([]string, len(e.header()))
}

func (e *Estimator) header() []string {
	return []string{dataName, dataType, minRow, maxRow, xRow(e.perN, false), xRow(e.perN, true)}
}

// render prints results inside a ASCII-table format.
func (e *Estimator) render(writer io.Writer, data [][]string) error {
	w := tablewriter.NewWriter(writer)
	w.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	w.SetHeader(e.header())
	w.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT,
	})
	w.AppendBulk(data)
	w.Render()
	return nil
}

func (e *Estimator) row(data ds.Data) []string {
	min, max := data.Size()
	return []string{
		data.String(),
		data.Kind(),
		ds.HumanSize(min, e.precision),
		ds.HumanSize(max, e.precision),
		ds.HumanSize(min*e.perN, e.precision),
		ds.HumanSize(max*e.perN, e.precision),
	}
}

const (
	base10 = 10
	bits64 = 64
)

func xRow(i uint64, max bool) string {
	var (
		m = "min"
		s = strconv.FormatUint(i, base10)
	)
	if max {
		m = "max"
	}
	return fmt.Sprintf("X %s (%s)", s, m)
}
