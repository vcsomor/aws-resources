package jsonfile

import (
	"encoding/json"
	"fmt"
	"github.com/vcsomor/aws-resources/internal/lister/writer"
	"os"
	"path/filepath"
)

const (
	DefaultOutputFile  = "resource.json"
	DefaultIndentation = "\t"
)

type Options struct {
	outputFile string
	indent     string
}

type OptionFnc func(*Options) error

func WithOutputFile(outputFile string) OptionFnc {
	return func(options *Options) error {
		if !filepath.IsLocal(outputFile) {
			return fmt.Errorf("the file is not a relative filename %s", outputFile)
		}
		options.outputFile = outputFile
		return nil
	}
}

func WithIndentation(indentation string) OptionFnc {
	return func(options *Options) error {
		options.indent = indentation
		return nil
	}
}

type jsonFileWriter struct {
	to   string
	opts Options
}

var _ writer.Writer = (*jsonFileWriter)(nil)

func NewWriter(to string, opts ...OptionFnc) (writer.Writer, error) {
	options := Options{
		outputFile: DefaultOutputFile,
		indent:     DefaultIndentation,
	}
	for _, fn := range opts {
		if err := fn(&options); err != nil {
			return nil, fmt.Errorf("unable to create the Json File Writer: %w", err)
		}
	}

	return &jsonFileWriter{
		to:   to,
		opts: options,
	}, nil
}

func (w jsonFileWriter) Write(obj any) error {
	bytes, err := w.serialize(obj)
	if err != nil {
		return err
	}

	err = w.makeDirectory(w.to)
	if err != nil {
		return err
	}

	return w.writeResource(bytes)
}

func (w jsonFileWriter) serialize(obj any) ([]byte, error) {
	return json.MarshalIndent(obj, "", w.opts.indent)
}

func (w jsonFileWriter) writeResource(b []byte) error {
	return writeFile(filepath.Join(w.to, w.opts.outputFile), b)
}

func (w jsonFileWriter) makeDirectory(to string) error {
	return os.MkdirAll(to, os.ModePerm)
}

func writeFile(to string, b []byte) error {
	f, err := os.Create(to)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.Write(b)
	return err
}
