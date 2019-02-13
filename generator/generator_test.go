package generator

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	src := `
		package main

		import (
			"golang.org/x/net/context"
			"github.com/stretchr/testify/assert"
		)

		type Example interface {
			Init(ctx context.Context) (*assert.TestingT, error)
		}
`
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "demo", src, parser.ParseComments)
	assert.Nil(t, err)

	e := file.Scope.Lookup("Example")
	iType := e.Decl.(*ast.TypeSpec)

	// Type-check the package.
	// We create an empty map for each kind of input
	// we're interested in, and Check populates them.
	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{
		Importer: importer.Default(),
	}
	_, errCheck := conf.Check("main", fSet, []*ast.File{file}, &info)
	if errCheck != nil {
		log.Fatal(errCheck)
	}

	o := info.Defs[iType.Name]
	cmd := &Command{
		IFaceType:     o.Type().Underlying().(*types.Interface),
		NoopName:      "NoopExample",
		InterfaceName: "Example",
		Out:           os.Stdout,
		PackageName:   "demo",
	}
	assert.NoError(t, Process(cmd))
}
