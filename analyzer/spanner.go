package analyzer

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

// https://pkg.go.dev/github.com/cloudspannerecosystem/memefish@v0.0.0-20231128072053-0a1141e8eb65/ast

type SpannerAnalyzer struct{}

func NewSpannerAnalyzer() *SpannerAnalyzer {
	return &SpannerAnalyzer{}
}

func (a *SpannerAnalyzer) Analyze(contents string) (bool, error) {
	p := &memefish.Parser{
		Lexer: &memefish.Lexer{
			File: &token.File{
				Buffer: heredoc.Doc(contents),
			},
		},
	}

	ddls, err := p.ParseDDLs()
	if err != nil {
		return true, err
	}

	for _, ddl := range ddls {
		if a.isSingleStmtHazardous(ddl) {
			return true, nil
		}
	}
	return false, nil
}

func (a *SpannerAnalyzer) isSingleStmtHazardous(ddl ast.DDL) bool {
	switch ddl.(type) {
	case *ast.CreateChangeStream, *ast.CreateDatabase, *ast.CreateIndex, *ast.CreateRole, *ast.CreateSequence, *ast.CreateTable, *ast.CreateView:
		return false
	case *ast.AlterTable:
		return a.isAlterHazardous(ddl)
	default:
	}
	return true
}

func (a *SpannerAnalyzer) isAlterHazardous(ddl ast.DDL) bool {
	alterTable, _ := ddl.(*ast.AlterTable)
	switch alterTable.TableAlteration.(type) {
	// ALTER TABLE ... ADD COLUMN ... is safe
	case *ast.AddColumn:
		return false
	default:
	}
	return true
}
