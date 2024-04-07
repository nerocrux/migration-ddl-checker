package analyzer_test

import (
	"testing"

	"github.com/nerocrux/migration-ddl-checker/analyzer"
	"github.com/nerocrux/migration-ddl-checker/ddl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var query_Mysql_CREATE = []string{
	`CREATE TABLE users (id INT)`,
	`CREATE INDEX idx_users ON users (id)`,
	`CREATE UNIQUE INDEX idx_users ON users (id)`,
	`ALTER TABLE users ADD COLUMN name VARCHAR(255)`,
	`ALTER TABLE users ADD INDEX idbyname (name)`,
}

var query_Mysql_DROP = []string{
	`DROP TABLE users`,
	`ALTER TABLE users DROP COLUMN name`,
	`DROP INDEX idx_users ON users`,
}

func TestAnalyzeMysql_CREATE(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		want     bool
	}{
		{
			name:     "create queries",
			contents: query_Mysql_CREATE,
			want:     true,
		},
		{
			name:     "drop queries",
			contents: query_Mysql_DROP,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := analyzer.NewMysqlAnalyzer(ddl.FromConfig([]string{"CREATE"}))
			for _, q := range tt.contents {
				got, err := a.Analyze(q)
				require.NoError(t, err, q)
				assert.Equal(t, tt.want, got, q)
			}
		})
	}
}

func TestAnalyzeMysql_DROP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		want     bool
	}{
		{
			name:     "create queries",
			contents: query_Mysql_CREATE,
			want:     false,
		},
		{
			name:     "drop queries",
			contents: query_Mysql_DROP,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := analyzer.NewMysqlAnalyzer(ddl.FromConfig([]string{"DROP"}))
			for _, q := range tt.contents {
				got, err := a.Analyze(q)
				require.NoError(t, err, q)
				assert.Equal(t, tt.want, got, q)
			}
		})
	}
}

func TestAnalyzeMysql_ALTER(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		want     bool
	}{
		{
			name:     "create queries",
			contents: query_Mysql_CREATE,
			want:     false,
		},
		{
			name:     "drop queries",
			contents: query_Mysql_DROP,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := analyzer.NewMysqlAnalyzer(ddl.FromConfig([]string{"ALTER"}))
			for _, q := range tt.contents {
				got, err := a.Analyze(q)
				require.NoError(t, err, q)
				assert.Equal(t, tt.want, got, q)
			}
		})
	}
}

func TestAnalyzeMysql_MULTI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		want     bool
	}{
		{
			name:     "create queries",
			contents: query_Mysql_CREATE,
			want:     true,
		},
		{
			name:     "drop queries",
			contents: query_Mysql_DROP,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := analyzer.NewMysqlAnalyzer(ddl.FromConfig([]string{"CREATE", "DROP", "ALTER"}))
			for _, q := range tt.contents {
				got, err := a.Analyze(q)
				require.NoError(t, err, q)
				assert.Equal(t, tt.want, got, q)
			}
		})
	}
}
