package gron

import (
	"encoding/json"
	"io"
)

type Gron interface {
	NextString() (string, error)
	NextStatement() (*Statement, error)

	String() (string, error)
	StringArray() ([]string, error)
	StatementArray() ([]Statement, error)
}

func New(reader io.Reader, options ...Option) Gron {
	return newGronImpl(reader, options...)
}

type Formatter interface {
	FormatString(string) string
	FormatNumber(json.Number) string
	FormatNull() string
	FormatBool(bool) string
	FormatObject(key string, isFirst bool) string
	FormatArray(index int) string
	FormatStatement(key, value string) string
	FormatEmptyObject() string
	FormatEmptyArray() string
}
