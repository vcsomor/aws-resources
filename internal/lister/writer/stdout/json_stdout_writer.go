package stdout

import (
	"encoding/json"
	"fmt"
	"github.com/vcsomor/aws-resources/internal/lister/writer"
)

const (
	DefaultIndentation = "\t"
)

type Options struct {
	indent string
}

type OptionFnc func(*Options) error

func WithIndentation(indentation string) OptionFnc {
	return func(options *Options) error {
		options.indent = indentation
		return nil
	}
}

type jsonStdoutWriter struct {
	Options
}

var _ writer.Writer = (*jsonStdoutWriter)(nil)

func NewWriter(opts ...OptionFnc) (writer.Writer, error) {
	options := Options{
		indent: DefaultIndentation,
	}
	for _, fn := range opts {
		if err := fn(&options); err != nil {
			return nil, fmt.Errorf("unable to create the Json Stdout Writer: %w", err)
		}
	}

	return &jsonStdoutWriter{
		Options: options,
	}, nil
}

func (w jsonStdoutWriter) Write(obj any) error {
	bytes, err := w.serialize(obj)
	if err != nil {
		return err
	}

	w.writeResource(bytes)
	return nil
}

func (w jsonStdoutWriter) serialize(obj any) ([]byte, error) {
	return json.MarshalIndent(obj, "", w.indent)
}

func (w jsonStdoutWriter) writeResource(b []byte) {
	fmt.Printf("%s\n", string(b))
}
