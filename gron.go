package gron

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	unknown = iota
	expectDocKey
	expectDocValue
	expectArrayValue
)

type gronImpl struct {
	dec       *json.Decoder
	eof       bool
	formatter Formatter
	stack     *stack
}

func newGronImpl(reader io.Reader, options ...Option) *gronImpl {
	dec := json.NewDecoder(reader)
	dec.UseNumber()

	inst := gronImpl{
		dec:       dec,
		stack:     newStack(),
		formatter: &defaultFormatter{},
	}

	for _, opt := range options {
		opt(&inst)
	}

	return &inst
}

func newGronImplFromMarshal(v interface{}, options ...Option) (*gronImpl, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(b)
	inst := newGronImpl(reader, options...)

	return inst, nil
}

func newGronImplFromMap(m map[string]interface{}, options ...Option) (*gronImpl, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(b)
	inst := newGronImpl(reader, options...)

	return inst, nil
}

func (gr *gronImpl) NextStatement() (*Statement, error) {
	for {
		if gr.eof {
			return nil, io.EOF
		}
		s, err := gr.next()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return nil, err
		}
		if s == nil {
			continue
		}

		return s, nil
	}
}

func (gr *gronImpl) NextString() (string, error) {
	s, err := gr.NextStatement()
	if err != nil {
		return "", err
	}

	return s.String(), nil
}

func (gr *gronImpl) StatementArray() ([]Statement, error) {
	var result []Statement
	for {
		s, err := gr.NextStatement()
		if gr.eof {
			return result, nil
		}
		if err != nil {
			return result, err
		}
		result = append(result, *s)
	}

	return result, nil
}

func (gr *gronImpl) StringArray() ([]string, error) {
	sx, err := gr.StatementArray()
	if err != nil {
		return nil, err
	}

	var result []string

	for _, s := range sx {
		result = append(result, s.String())
	}

	return result, nil
}

func (gr *gronImpl) String() (string, error) {
	sx, err := gr.StringArray()
	if err != nil {
		return "", err
	}

	return strings.Join(sx, "\n"), nil
}

func (gr *gronImpl) next() (*Statement, error) {
	if gr.eof {
		return nil, io.EOF
	}

	t, err := gr.dec.Token()
	if err == io.EOF {
		gr.eof = true
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}

	switch t.(type) {
	case json.Delim:
		v := t.(json.Delim).String()

		if v == "{" || v == "[" {
			a, ok := gr.stack.Peek().(*array)
			if ok {
				a.Inc()
			}
		}

		switch v {
		case "{":
			gr.stack.Push(newDoc(gr.formatter, gr.stack.Empty()))
		case "}":
			return gr.handleDocEnd(), nil
		case "[":
			gr.stack.Push(newArray(gr.formatter))
		case "]":
			return gr.handleArrayEnd(), nil
		}
	case bool:
		return gr.handlePrimitive(t), nil
	case float64:
		return gr.handlePrimitive(t), nil
	case json.Number:
		return gr.handlePrimitive(t), nil
	case string:
		return gr.handlePrimitive(t), nil
	case nil:
		return gr.handlePrimitive(t), nil
	default:
		return nil, errors.New("unexpected value")
	}

	return nil, nil
}

func (gr *gronImpl) handleDocEnd() *Statement {
	var s *Statement
	if gr.stack.Peek().Empty() {
		gr.stack.Pop()
		emptyDoc := map[interface{}]interface{}{}
		s = newStatement(gr.stack.String(), gr.formatter.FormatEmptyObject(), emptyDoc, gr.formatter)
	} else {
		gr.stack.Pop()
	}
	d, ok := gr.stack.Peek().(*doc)
	if ok {
		d.next = expectDocKey
	}

	return s
}

func (gr *gronImpl) handleArrayEnd() *Statement {
	var s *Statement
	if gr.stack.Peek().Empty() {
		gr.stack.Pop()
		var k string
		if !gr.stack.Empty() {
			k = gr.stack.String()
		}
		emptyArray := []interface{}{}
		s = newStatement(k, gr.formatter.FormatEmptyArray(), emptyArray, gr.formatter)
	} else {
		gr.stack.Pop()
	}
	d, ok := gr.stack.Peek().(*doc)
	if ok {
		d.next = expectDocKey
	}

	return s
}

func (gr *gronImpl) handlePrimitive(t interface{}) *Statement {
	e := gr.stack.Peek()
	if e == nil {
		return newStatement("", gr.formatPrimitive(t), t, gr.formatter)
	}

	switch e.(type) {
	case *doc:
		d := e.(*doc)
		switch d.next {
		case expectDocKey:
			// a json object key must be a string
			d.SetKey(t.(string))
			d.next = expectDocValue
		case expectDocValue:
			s := gr.buildStatement(t)
			d.next = expectDocKey
			return s
		}
	case *array:
		a := e.(*array)
		a.Inc()
		s := gr.buildStatement(t)
		return s
	}

	return nil
}

func (gr *gronImpl) buildStatement(t interface{}) *Statement {
	return newStatement(gr.stack.String(), gr.formatPrimitive(t), t, gr.formatter)
}

func (gr *gronImpl) formatPrimitive(p interface{}) string {
	switch p.(type) {
	case bool:
		return gr.formatter.FormatBool(p.(bool))
	case string:
		return gr.formatter.FormatString(p.(string))
	case json.Number:
		return gr.formatter.FormatNumber(p.(json.Number))
	case nil:
		return gr.formatter.FormatNull()
	default:
		return fmt.Sprintf("<%s>", p)
	}
}

func (gr *gronImpl) Map() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for {
		s, err := gr.NextStatement()
		if gr.eof {
			return result, nil
		}
		if err != nil {
			return result, err
		}
		result[s.key] = s.rawValue
	}
}

func (gr *gronImpl) MarshalJSON() ([]byte, error) {
	m, err := gr.Map()
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}
