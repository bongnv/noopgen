package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bongnv/noopgen/logger"
	"github.com/bongnv/noopgen/walker"
)

type config struct {
	baseDir       string
	interfaceName string
}

func main() {
	c := parseConfigFromArgs(os.Args)
	w := &walker.Walker{
		BaseDir:       c.baseDir,
		InterfaceName: c.interfaceName,
		WriterProvider: &walker.FileWriterProvider{
			Path: path.Join(c.baseDir, "z_noop_"+underscoreCaseName(c.interfaceName)+".go"),
		},
		Ext: ".go",
	}
	err := filepath.Walk(c.baseDir, w.Walk)
	if err != nil {
		logger.Error("noopgen", "Error while walking path %v, err: %v.", c.baseDir, err)
		os.Exit(1)
	}
}

func parseConfigFromArgs(args []string) *config {
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)

	c := &config{}

	flagSet.StringVar(&c.baseDir, "dir", ".", "directory to search for interfaces")
	flagSet.StringVar(&c.interfaceName, "i", "", "interface name")

	_ = flagSet.Parse(args[1:])
	logger.Info("noopgen", "Config: %+v", c)
	return c
}

func underscoreCaseName(caseName string) string {
	rxp1 := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s1 := rxp1.ReplaceAllString(caseName, "${1}_${2}")
	rxp2 := regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(rxp2.ReplaceAllString(s1, "${1}_${2}"))
}
