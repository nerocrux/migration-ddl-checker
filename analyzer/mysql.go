package analyzer

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/nerocrux/migration-ddl-checker/ddl"
	"github.com/xwb1989/sqlparser"
)

// https://pkg.go.dev/github.com/xwb1989/sqlparser

var createWords = []string{
	"CREATE INDEX",
	"CREATE UNIQUE INDEX",
	"ADD",
	"ADD COLUMN",
}

var dropWords = []string{
	"DROP INDEX",
	"DROP COLUMN",
}

type MysqlAnalyzer struct {
	HazardousDDLs []ddl.DDL
}

func NewMysqlAnalyzer(ddl []ddl.DDL) *MysqlAnalyzer {
	return &MysqlAnalyzer{
		HazardousDDLs: ddl,
	}
}

func (a *MysqlAnalyzer) Analyze(contents string) (bool, error) {
	for _, sql := range strings.Split(contents, ";") {
		isHazadous, err := a.isSingleStmtHazardous(sql)
		if err != nil {
			slog.Error(err.Error())
		}
		if isHazadous {
			return true, nil
		}
	}
	return false, nil
}

func (a *MysqlAnalyzer) IsHazardousDDL(d ddl.DDL) bool {
	return slices.Contains(a.HazardousDDLs, d)
}

func (a *MysqlAnalyzer) isSingleStmtHazardous(sql string) (bool, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return true, err
	}

	if d, ok := stmt.(*sqlparser.DDL); ok {
		switch d.Action {
		case sqlparser.CreateStr, sqlparser.CreateVindexStr, sqlparser.AddColVindexStr:
			return a.IsHazardousDDL(ddl.CreateDDL), nil
		case sqlparser.DropStr, sqlparser.RenameStr, sqlparser.TruncateStr, sqlparser.DropColVindexStr:
			return a.IsHazardousDDL(ddl.DropDDL), nil
		// CREATE INDEX or DROP INDEX are judged in alter statement
		case sqlparser.AlterStr:
			return a.isAlterHazardous(sql)
		default:
			return true, fmt.Errorf("unknown DDL action: %s in query: %s", d.Action, sql)
		}
	}

	// query except DDL (i.e. select, insert, delete, etc...) should be fine
	return false, nil
}

func (a *MysqlAnalyzer) isAlterHazardous(sql string) (bool, error) {
	for _, word := range createWords {
		if strings.Contains(strings.ToUpper(sql), word) {
			return a.IsHazardousDDL(ddl.CreateDDL), nil
		}
	}
	for _, word := range dropWords {
		if strings.Contains(strings.ToUpper(sql), word) {
			return a.IsHazardousDDL(ddl.DropDDL), nil
		}
	}
	return true, fmt.Errorf("unknown ALTER action in query: %s", sql)
}
