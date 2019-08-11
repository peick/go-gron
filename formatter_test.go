package gron

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultFormatterFormatObject(t *testing.T) {
	formatter := NewDefaultFormatter()

	res := formatter.FormatObject("example", true)
	require.Equal(t, ".example", res)

	res = formatter.FormatObject("example", false)
	require.Equal(t, ".example", res)
}

func TestDefaultFormatterFormatString(t *testing.T) {
	formatter := NewDefaultFormatter()

	res := formatter.FormatString("example")
	require.Equal(t, "example", res)

	res = formatter.FormatString("example with spaces")
	require.Equal(t, "example with spaces", res)

	res = formatter.FormatString("example with\nnewline")
	require.Equal(t, "example with\\nnewline", res)

	res = formatter.FormatString(`example with "double-quote"`)
	require.Equal(t, `example with \"double-quote\"`, res)
}
