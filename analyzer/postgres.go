package analyzer

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type PostgresqlAnalyzer struct{}

func NewPostgresqlAnalyzer() *PostgresqlAnalyzer {
	return &PostgresqlAnalyzer{}
}

func (a *PostgresqlAnalyzer) Analyze(contents string) (bool, error) {
	tree, err := pg_query.Parse(string(contents))
	if err != nil {
		return true, err
	}
	for _, stmt := range tree.Stmts {
		if a.isSingleStmtHazardous(stmt) {
			return true, nil
		}
	}
	return false, nil
}

func (a *PostgresqlAnalyzer) isSingleStmtHazardous(stmt *pg_query.RawStmt) bool {
	switch stmt.Stmt.Node.(type) {
	case *pg_query.Node_CreateStmt, *pg_query.Node_CreateTableAsStmt:
		return false
	case *pg_query.Node_AlterTableStmt:
		return a.isAlterHazardous(stmt.Stmt.GetAlterTableStmt().GetObjtype().Enum())
	case *pg_query.Node_IndexStmt:
		// "ALTER TABLE ... RENAME ..." goes to Node_RenameStmt,
		// "DROP INDEX ..." goes to Node_DropStmt,
		// so it seems that only "CREATE INDEX ..." goes here, so it should be saft.
		return false
	default:
	}
	return false
}

func (a *PostgresqlAnalyzer) isAlterHazardous(typ protoreflect.Enum) bool {
	switch typ {
	case pg_query.AlterTableType_AT_AddColumn, pg_query.AlterTableType_AT_AddIndex:
		return false
	default:
	}
	return true
}
