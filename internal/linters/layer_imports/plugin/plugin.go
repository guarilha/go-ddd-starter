package main

import (
	"github.com/guarilha/go-ddd-starter/internal/linters/layer_imports"
	"golang.org/x/tools/go/analysis"
)

// New is the entry point for golangci-lint plugin
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		layer_imports.Analyzer,
	}, nil
}
