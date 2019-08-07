package gron

type Statement struct {
	key       string
	value     string
	formatter Formatter
}

func newStatement(key, value string, formatter Formatter) *Statement {
	return &Statement{
		key:       key,
		value:     value,
		formatter: formatter,
	}
}

func (s *Statement) String() string {
	return s.formatter.FormatStatement(s.key, s.value)
}

type element interface {
	String() string
	Empty() bool
}

type array struct {
	empty     bool
	index     int
	formatter Formatter
}

func newArray(formatter Formatter) element {
	return &array{
		empty:     true,
		index:     -1,
		formatter: formatter,
	}
}

func (a *array) Empty() bool {
	return a.index == -1
}

func (a *array) Inc() {
	a.index += 1
}

func (a *array) String() string {
	return a.formatter.FormatArray(a.index)
}

type doc struct {
	empty     bool
	formatter Formatter
	key       string
	next      int
	root      bool
}

func newDoc(formatter Formatter, root bool) element {
	return &doc{
		empty:     true,
		formatter: formatter,
		next:      expectDocKey,
		root:      root,
	}
}

func (d *doc) Empty() bool {
	return d.empty
}

func (d *doc) SetKey(key string) {
	d.empty = false
	d.key = key
}

func (d *doc) String() string {
	return d.formatter.FormatObject(d.key, false)
}
