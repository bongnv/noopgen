package walker

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bongnv/noopgen/generator"
	"github.com/bongnv/noopgen/logger"
)

// DefaultWalker ...
type Walker struct {
	BaseDir       string
	Ext           string
	InterfaceName string

	WriterProvider WriterProvider
}

// Walk ...
func (w *Walker) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		logger.Info("walker", "Error while processing file, err: %v", err)
		return err
	}

	if info.IsDir() {
		return nil
	}

	if !strings.HasSuffix(path, w.Ext) {
		return nil
	}

	return w.parseAndProcess(path)
}

func (w *Walker) parseAndProcess(filePath string) (err error) {
	logger.Info("walker", "Processing %v", filePath)
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	fSet := token.NewFileSet() // positions are relative to fset
	fNode, err := parser.ParseFile(fSet, absFilePath, nil, parser.ParseComments)

	obj := fNode.Scope.Lookup(w.InterfaceName)
	if obj == nil {
		return nil
	}
	iType := obj.Decl.(*ast.TypeSpec)

	// Type-check the package.
	// We create an empty map for each kind of input
	// we're interested in, and Check populates them.
	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}

	conf := types.Config{
		Importer:         importer.Default(),
		IgnoreFuncBodies: true,
	}
	_, errCheck := conf.Check(filePath, fSet, []*ast.File{fNode}, &info)
	if errCheck != nil {
		log.Fatal(errCheck)
	}

	o := info.Defs[iType.Name]

	writer, err := w.WriterProvider.New()
	if err != nil {
		logger.Error("walker", "Failed to create writer, err: %v", err)
		return err
	}

	cmd := &generator.Command{
		IFaceType:     o.Type().Underlying().(*types.Interface),
		NoopName:      "Noop" + w.InterfaceName,
		InterfaceName: w.InterfaceName,
		Out:           writer,
		PackageName:   fNode.Name.String(),
	}
	return generator.Process(cmd)
}
