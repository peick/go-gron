package gron

import (
	"encoding/json"
	"fmt"
)

type defaultFormatter struct {
}

func NewDefaultFormatter() Formatter {
	return &defaultFormatter{}
}

func (f *defaultFormatter) FormatString(s string) string {
	b, _ := json.Marshal(s)
	s = string(b)
	s = s[1 : len(s)-1]
	return s
}

func (f *defaultFormatter) FormatNumber(n json.Number) string {
	return string(n)
}

func (f *defaultFormatter) FormatNull() string {
	return "null"
}

func (f *defaultFormatter) FormatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func (f *defaultFormatter) FormatObject(key string, _ bool) string {
	return "." + key
}

func (f *defaultFormatter) FormatArray(index int) string {
	return fmt.Sprintf("[%d]", index)
}

func (f *defaultFormatter) FormatStatement(key string, value string) string {
	if key == "" {
		return value
	}
	return fmt.Sprintf("%s: %s", key, value)
}

func (f *defaultFormatter) FormatEmptyObject() string {
	return "{}"
}

func (f *defaultFormatter) FormatEmptyArray() string {
	return "[]"
}

type originalGronFormatter struct {
	defaultFormatter
}

func NewOriginalGronFormatter() Formatter {
	return &originalGronFormatter{}
}

func (f *originalGronFormatter) FormatString(s string) string {
	b, _ := json.Marshal(s)
	s = string(b)
	return s
}

func (f *originalGronFormatter) FormatNumber(n json.Number) string {
	return string(n)
}

func (f *originalGronFormatter) FormatObject(key string, isFirst bool) string {
	if isFirst {
		return key
	}
	return "." + key
}

func (f *originalGronFormatter) FormatStatement(key string, value string) string {
	if key == "" {
		return value
	}
	return fmt.Sprintf("%s: %s", key, value)
}
