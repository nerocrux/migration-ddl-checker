package analyzer

import (
	"slices"

	"github.com/nerocrux/migration-ddl-checker/ddl"
	pg_query "github.com/pganalyze/pg_query_go/v5"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type PostgresqlAnalyzer struct {
	HazardousDDLs []ddl.DDL
}

func NewPostgresqlAnalyzer(ddl []ddl.DDL) *PostgresqlAnalyzer {
	return &PostgresqlAnalyzer{
		HazardousDDLs: ddl,
	}
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

func (a *PostgresqlAnalyzer) IsHazardousDDL(d ddl.DDL) bool {
	return slices.Contains(a.HazardousDDLs, d)
}

func (a *PostgresqlAnalyzer) isSingleStmtHazardous(stmt *pg_query.RawStmt) bool {
	switch stmt.Stmt.Node.(type) {
	case *pg_query.Node_CreateStmt, *pg_query.Node_CreateTableAsStmt:
		return a.IsHazardousDDL(ddl.CreateDDL)
	// "ALTER TABLE ... RENAME ..." goes to Node_RenameStmt,
	// "DROP INDEX ..." goes to Node_DropStmt,
	// so it seems that only "CREATE INDEX ..." goes here, so it should be safe.
	case *pg_query.Node_IndexStmt:
		return a.IsHazardousDDL(ddl.CreateDDL)
	case *pg_query.Node_AlterTableStmt:
		return a.isAlterHazardous(stmt.Stmt.GetAlterTableStmt().GetObjtype().Enum())
	case *pg_query.Node_DropStmt, *pg_query.Node_DeleteStmt:
		return a.IsHazardousDDL(ddl.DropDDL)
	default:
	}
	return false
}

func (a *PostgresqlAnalyzer) isAlterHazardous(typ protoreflect.Enum) bool {
	switch typ {
	case pg_query.AlterTableType_AT_AddColumn, pg_query.AlterTableType_AT_AddIndex:
		return a.IsHazardousDDL(ddl.CreateDDL)
	case pg_query.AlterTableType_AT_DropColumn:
		return a.IsHazardousDDL(ddl.DropDDL)
	default:
	}
	return false
}
