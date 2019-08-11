package gron

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringArray(t *testing.T) {
	testCases := []struct {
		input  string
		expect []string
	}{
		{`""`, []string{""}},

		{`"sample"`, []string{"sample"}},

		{`null`, []string{"null"}},

		{`42`, []string{"42"}},

		{`4.21`, []string{"4.21"}},

		{`true`, []string{"true"}},

		{`{}`, []string{"{}"}},

		{`{"": 1}`, []string{".: 1"}},

		{`{"": ""}`, []string{".: "}},

		{`{"sample": "example"}`, []string{".sample: example"}},

		{`{"sample": null}`, []string{".sample: null"}},

		{`{"sample": "null"}`, []string{".sample: null"}},

		{`{"sample": 42}`, []string{".sample: 42"}},

		{`{"sample": "42"}`, []string{".sample: 42"}},

		{`{"sample": 4.21}`, []string{".sample: 4.21"}},

		{`{"sample": "4.21"}`, []string{".sample: 4.21"}},

		{`{"sample": true}`, []string{".sample: true"}},

		{`{"sample": "true"}`, []string{".sample: true"}},

		{`{"key": {"sub": "value"}}`, []string{".key.sub: value"}},

		{`{"key1": {"sub1": "sub-value1", "sub2": "sub-value2"}, "key2": "value2"}`,
			[]string{".key1.sub1: sub-value1", ".key1.sub2: sub-value2", ".key2: value2"}},

		{`{"key1": [], "key2": 42}`, []string{".key1: []", ".key2: 42"}},

		{`{"key1": {}, "key2": 42}`, []string{".key1: {}", ".key2: 42"}},

		{`{"sample": "Lorem\nipsum\ndolor"}`, []string{`.sample: Lorem\nipsum\ndolor`}},

		{`{"sample": {}}`, []string{".sample: {}"}},

		{`{"sample": []}`, []string{".sample: []"}},

		{`[]`, []string{"[]"}},

		{`[1]`, []string{"[0]: 1"}},

		{`[[]]`, []string{"[0]: []"}},

		{`[{}]`, []string{"[0]: {}"}},
	}

	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			gr := New(reader)
			res, err := gr.StringArray()
			require.NoError(t, err)
			require.Equal(t, tt.expect, res)
		})
	}
}

func TestNewFromMarshalable(t *testing.T) {
	type sampleStruct struct {
		A int
		B string `json:"b"`
	}

	testCases := []struct {
		name   string
		input  interface{}
		expect []string
	}{
		{"struct",
			sampleStruct{43, "test"},
			[]string{".A: 43", ".b: test"},
		},

		{"map",
			map[string]interface{}{"A": 43, "b": "test"},
			[]string{".A: 43", ".b: test"},
		},

		{"string",
			"test",
			[]string{"test"},
		},

		{"int",
			44,
			[]string{"44"},
		},

		{"array",
			[]interface{}{45, "test"},
			[]string{"[0]: 45", "[1]: test"},
		},

		{"nil",
			nil,
			[]string{"null"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gr, err := NewFromMarshalable(tt.input)
			require.NoError(t, err)
			require.NotNil(t, gr)
			res, err := gr.StringArray()
			require.NoError(t, err)
			require.Equal(t, tt.expect, res)
		})
	}
}
