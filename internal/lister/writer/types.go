package writer

type Writer interface {
	Write(data any) error
}
