package walker

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalker(t *testing.T) {
	w := &Walker{
		InterfaceName:  "Example",
		WriterProvider: &StdoutProvider{},
		Ext:            ".go.sample",
	}
	assert.NoError(t, filepath.Walk("../samples/with_error.go.sample", w.Walk))
}
