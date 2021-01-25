package xstrings

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
)

type Arguments struct {
	Comment rune
}

func (a *Arguments) Parse(str string) ([]string, error) {
	cr := csv.NewReader(bytes.NewBufferString(str))
	cr.Comma = ' '
	cr.Comment = a.Comment
	cr.FieldsPerRecord = -1
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true
	cr.ReuseRecord = false

	args, err := cr.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return []string{}, nil
		}
		return nil, newParseError(err)
	}

	if l := len(args); l > 0 && args[l-1] == "" {
		args = args[:l-1]
	}

	return args, nil
}

func (a *Arguments) Format(args ...string) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	cw := csv.NewWriter(buf)
	cw.Comma = ' '

	err := cw.WriteAll([][]string{args})
	if err != nil {
		return "", newFormatError(err)
	}

	return string(bytes.TrimSuffix(buf.Bytes(), []byte("\n"))), nil
}
