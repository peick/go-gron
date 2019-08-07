package gron

type Option func(*gronImpl)

func DefaultFormatter() Option {
	return func(gr *gronImpl) {
		gr.formatter = NewDefaultFormatter()
	}
}

func OriginalGronFormatter() Option {
	return func(gr *gronImpl) {
		gr.formatter = NewOriginalGronFormatter()
	}
}
