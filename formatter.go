package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Formatter interface {
	Write(Slicer) error
}

type PlainText struct {
	output io.Writer
}

func NewPlainText(output io.Writer) *PlainText {
	return &PlainText{
		output: output,
	}
}

func (p *PlainText) Write(s Slicer) error {
	_, err := p.output.Write([]byte(fmt.Sprintf("%+v\n", s.Slice())))
	return err
}

type CSVFormatter struct {
	output *csv.Writer
}

func NewCSVFormatter(output io.Writer) (*CSVFormatter, error) {
	formatter := &CSVFormatter{
		output: csv.NewWriter(output),
	}
	err := formatter.output.Write(Headers())
	return formatter, err
}

func (c *CSVFormatter) Write(s Slicer) error {
	defer c.output.Flush()
	return c.output.Write(s.Slice())
}
