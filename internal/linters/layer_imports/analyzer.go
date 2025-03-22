package layer_imports

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name: "layer_imports",
	Doc:  "Checks that layer imports follow the dependency rule",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

// determineLayer identifies which architectural layer a file belongs to
func determineLayer(filename string) string {
	switch {
	case strings.Contains(filename, "/app/"):
		return "app"
	case strings.Contains(filename, "/domain/"):
		return "domain"
	case strings.Contains(filename, "/gateways/"):
		return "gateways"
	case strings.Contains(filename, "/internal/"):
		return "internal"
	default:
		return "unknown"
	}
}

// isExternalPackage checks if an import path is from an external package
func isExternalPackage(importPath string) bool {
	return !strings.Contains(importPath, "/app/") &&
		!strings.Contains(importPath, "/domain/") &&
		!strings.Contains(importPath, "/gateways/") &&
		!strings.Contains(importPath, "/internal/")
}

// isExemptDomainFile checks if the file is one of the special domain files that are exempt from import restrictions
func isExemptDomainFile(filename string) bool {
	// Convert to forward slashes for consistency across platforms
	filename = filepath.ToSlash(filename)

	// Check if it's the root domains.go file
	if strings.HasSuffix(filename, "/domain/domains.go") {
		return true
	}

	// Check if it's a domain.go file inside a domain package
	parts := strings.Split(filename, "/")
	if len(parts) >= 4 && // Minimum path: <root>/domain/<name>/domain.go
		parts[len(parts)-2] != "domain" && // Not the root domain folder
		parts[len(parts)-1] == "domain.go" { // Is domain.go file
		return true
	}

	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		currentLayer := determineLayer(filename)

		// Skip unknown layers and test files
		if currentLayer == "unknown" || strings.HasSuffix(filename, "_test.go") {
			continue
		}

		// Check if this is an exempt domain file
		isExempt := isExemptDomainFile(filename)

		for _, imp := range file.Imports {
			// Get the import path without quotes
			importPath := strings.Trim(imp.Path.Value, "\"")

			// Skip external packages
			if isExternalPackage(importPath) {
				continue
			}

			switch currentLayer {
			case "app":
				// app/ can import from domain/, gateways/, internal/
				// All internal imports are allowed for app layer
				continue

			case "domain":
				// If it's an exempt file, allow all imports
				if isExempt {
					continue
				}

				// domain/ can only import from internal/ and other domain/ packages
				if strings.Contains(importPath, "/app/") || strings.Contains(importPath, "/gateways/") {
					pass.Report(analysis.Diagnostic{
						Pos: imp.Pos(),
						Message: fmt.Sprintf(
							"domain layer cannot import from app or gateways layers: %s",
							importPath,
						),
						Category: "layer-imports",
					})
				}

			case "gateways":
				// gateways/ can only import from internal/
				if strings.Contains(importPath, "/app/") || strings.Contains(importPath, "/domain/") {
					pass.Report(analysis.Diagnostic{
						Pos: imp.Pos(),
						Message: fmt.Sprintf(
							"gateways layer cannot import from app or domain layers: %s",
							importPath,
						),
						Category: "layer-imports",
					})
				}

			case "internal":
				// internal/ should not import from other layers
				if strings.Contains(importPath, "/app/") ||
					strings.Contains(importPath, "/domain/") ||
					strings.Contains(importPath, "/gateways/") {
					pass.Report(analysis.Diagnostic{
						Pos: imp.Pos(),
						Message: fmt.Sprintf(
							"internal layer cannot import from app, domain, or gateways layers: %s",
							importPath,
						),
						Category: "layer-imports",
					})
				}
			}
		}
	}

	return nil, nil
}
