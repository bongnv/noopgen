package walker

import (
	"io"
	"os"
)

// WriterProvider ...
type WriterProvider interface {
	New() (io.WriteCloser, error)
}

// StdoutProvider ...
type StdoutProvider struct{}

// New ...
func (p *StdoutProvider) New() (io.WriteCloser, error) {
	return os.Stdout, nil
}

type FileWriterProvider struct {
	Path string
}

// New ...
func (p *FileWriterProvider) New() (io.WriteCloser, error) {
	f, err := os.Create(p.Path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
