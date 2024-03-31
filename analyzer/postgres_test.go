package analyzer_test

import (
	"testing"

	"github.com/nerocrux/migration-ddl-checker/analyzer"
	"github.com/nerocrux/migration-ddl-checker/ddl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var query_Postgres_CREATE = []string{
	`CREATE TABLE users (id INT)`,
	`CREATE INDEX idx_users ON users (id) WHERE deleted_at IS NULL`,
	`CREATE UNIQUE INDEX idx_users ON users (id) WHERE deleted_at IS NULL`,
}

var query_Postgres_DROP = []string{
	`DROP TABLE users`,
	`DROP INDEX idx_users`,
}

var query_Postgres_ALTER_ADD = []string{
	`ALTER TABLE users ADD COLUMN name VARCHAR(255)`,
}

var query_Postgres_ALTER_DROP = []string{
	`ALTER TABLE users DROP COLUMN name`,
}

func TestAnalyzePostgres_CREATE(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		want     bool
	}{
		{
			name:     "create queries",
			contents: query_Postgres_CREATE,
			want:     true,
		},
		{
			name:     "drop queries",
			contents: query_Postgres_DROP,
			want:     false,
		},
		{
			name:     "alter add queries",
			contents: query_Postgres_ALTER_ADD,
			want:     false,
		},
		{
			name:     "alter drop queries",
			contents: query_Postgres_ALTER_DROP,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := analyzer.NewPostgresqlAnalyzer(ddl.FromConfig([]string{"CREATE"}))
			for _, q := range tt.contents {
				got, err := a.Analyze(q)
				require.NoError(t, err, q)
				assert.Equal(t, tt.want, got, q)
			}
		})
	}
}
