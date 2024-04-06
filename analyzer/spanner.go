package analyzer

import (
	"slices"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/nerocrux/migration-ddl-checker/ddl"
)

// https://pkg.go.dev/github.com/cloudspannerecosystem/memefish@v0.0.0-20231128072053-0a1141e8eb65/ast

type SpannerAnalyzer struct {
	HazardousDDLs []ddl.DDL
}

func NewSpannerAnalyzer(ddl []ddl.DDL) *SpannerAnalyzer {
	return &SpannerAnalyzer{
		HazardousDDLs: ddl,
	}
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

func (a *SpannerAnalyzer) IsHazardousDDL(d ddl.DDL) bool {
	return slices.Contains(a.HazardousDDLs, d)
}

func (a *SpannerAnalyzer) isSingleStmtHazardous(d ast.DDL) bool {
	switch d.(type) {
	case *ast.CreateChangeStream, *ast.CreateDatabase, *ast.CreateIndex, *ast.CreateRole, *ast.CreateSequence, *ast.CreateTable, *ast.CreateView:
		return a.IsHazardousDDL(ddl.CreateDDL)
	case *ast.DropIndex, *ast.DropTable:
		return a.IsHazardousDDL(ddl.DropDDL)
	case *ast.AlterTable:
		return a.isAlterHazardous(d)
	default:
	}
	return true
}

func (a *SpannerAnalyzer) isAlterHazardous(d ast.DDL) bool {
	alterTable, _ := d.(*ast.AlterTable)
	switch alterTable.TableAlteration.(type) {
	// ALTER TABLE ... ADD COLUMN ...
	case *ast.AddColumn:
		return a.IsHazardousDDL(ddl.CreateDDL)
	case *ast.DropColumn:
		return a.IsHazardousDDL(ddl.DropDDL)
	case *ast.AlterColumn:
		return a.IsHazardousDDL(ddl.AlterDDL)
	default:
	}
	return true
}
