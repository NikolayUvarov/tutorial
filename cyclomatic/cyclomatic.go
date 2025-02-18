package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// Main function
func main() {
	// Get directory path from command-line arguments
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	// Process all Go files in the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Ошибка доступа:", err)
			return nil
		}

		// Process only Go files
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			complexity, err := calculateCyclomaticComplexity(path)
			if err == nil {
				fmt.Printf("Файл: %s, Цикломатическая сложность: %d\n", path, complexity)
			} else {
				fmt.Printf("Ошибка обработки %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Ошибка обхода директорий:", err)
	}
}

// Function to calculate cyclomatic complexity of a Go file
func calculateCyclomaticComplexity(filename string) (int, error) {
	// Read the file
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if err != nil {
		return 0, err
	}

	// Base complexity is 1
	complexity := 1

	// Walk through the AST and count complexity-increasing nodes
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.BinaryExpr:
			complexity++
		}
		return true
	})

	return complexity, nil
}

