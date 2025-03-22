package domain_func_signature

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name: "domain_func_signature",
	Doc:  "Checks that domain functions follow signature conventions",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func isNativeType(typ ast.Expr) bool {
	switch t := typ.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string", "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "bool", "byte", "rune":
			return true
		}
	}
	return false
}

func isContextType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "context" && sel.Sel.Name == "Context"
		}
	}
	return false
}

// getTotalParams counts the actual number of parameters, taking into account
// that a single field can have multiple names
func getTotalParams(params *ast.FieldList) int {
	if params == nil {
		return 0
	}
	total := 0
	for _, field := range params.List {
		if len(field.Names) == 0 {
			total++ // Anonymous parameter
		} else {
			total += len(field.Names) // Multiple names share the same type
		}
	}
	return total
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		if !strings.Contains(filename, "/domain/") {
			continue
		}

		// Skip domains.go and repository folders
		if strings.HasSuffix(filename, "/domains.go") || strings.Contains(filename, "/repository/") {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if funcDecl.Name.Name == "New" {
				return true
			}

			// Count actual parameters (excluding receiver if it exists)
			totalParams := getTotalParams(funcDecl.Type.Params)

			// Check if we have too many parameters
			if totalParams > 2 {
				pass.Report(analysis.Diagnostic{
					Pos:      funcDecl.Type.Params.Pos(),
					Message:  fmt.Sprintf("domain function %s must have at most 2 parameters, got %d", funcDecl.Name.Name, totalParams),
					Category: "domain-function-signature",
				})
				return true
			}

			// If we have parameters, validate them
			if totalParams > 0 {
				params := funcDecl.Type.Params.List

				// First parameter must be context.Context
				firstField := params[0]
				if len(firstField.Names) > 1 || !isContextType(firstField.Type) {
					pass.Report(analysis.Diagnostic{
						Pos:      firstField.Pos(),
						Message:  fmt.Sprintf("first parameter of domain function %s must be context.Context", funcDecl.Name.Name),
						Category: "domain-function-signature",
					})
				}

				// If we have a second parameter group, validate its type
				if len(params) > 1 {
					secondField := params[1]
					isValidType := isNativeType(secondField.Type)

					if !isValidType {
						// Check if it's a struct type
						if _, ok := secondField.Type.(*ast.StructType); ok {
							isValidType = true
						}
					}

					if !isValidType {
						// Check if it's a named type ending with Params
						if ident, ok := secondField.Type.(*ast.Ident); ok {
							if strings.HasSuffix(ident.Name, "Params") {
								isValidType = true
							}
						}
					}

					if !isValidType {
						// Check if it's an imported type (selector expression)
						if _, ok := secondField.Type.(*ast.SelectorExpr); ok {
							isValidType = true
						}
					}

					if !isValidType {
						pass.Report(analysis.Diagnostic{
							Pos:      secondField.Pos(),
							Message:  fmt.Sprintf("second parameter of domain function %s must be either a struct or a native Go type", funcDecl.Name.Name),
							Category: "domain-function-signature",
						})
					}
				}
			}

			// Count return values
			resultCount := 0
			if funcDecl.Type.Results != nil {
				resultCount = funcDecl.Type.Results.NumFields()
			}

			if resultCount < 1 || resultCount > 2 {
				pass.Report(analysis.Diagnostic{
					Pos:      funcDecl.Type.Results.Pos(),
					Message:  fmt.Sprintf("domain function %s must have 1 or 2 return values, got %d", funcDecl.Name.Name, resultCount),
					Category: "domain-function-signature",
				})
			}

			return true
		})
	}

	return nil, nil
}
