package analyzer

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// https://pkg.go.dev/github.com/xwb1989/sqlparser

type MysqlAnalyzer struct{}

func NewMysqlAnalyzer() *MysqlAnalyzer {
	return &MysqlAnalyzer{}
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

func (a *MysqlAnalyzer) isSingleStmtHazardous(sql string) (bool, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return true, err
	}

	if ddl, ok := stmt.(*sqlparser.DDL); ok {
		switch ddl.Action {
		case sqlparser.CreateStr, sqlparser.CreateVindexStr, sqlparser.AddColVindexStr:
			return false, nil
		case sqlparser.AlterStr:
			return a.isAlterHazardous(sql), nil
		case sqlparser.DropStr, sqlparser.RenameStr, sqlparser.TruncateStr, sqlparser.DropColVindexStr:
			return true, nil
		default:
			return true, fmt.Errorf("unknown DDL action: %s in query: %s", ddl.Action, sql)
		}
	}

	// query except DDL (i.e. select, insert, delete, etc...) should be fine
	return false, nil
}

func (a *MysqlAnalyzer) isAlterHazardous(sql string) bool {
	var safeWords = []string{
		"ADD",
		"ADD COLUMN",
	}
	for _, word := range safeWords {
		if strings.Contains(strings.ToUpper(sql), word) {
			return false
		}
	}
	return true
}
