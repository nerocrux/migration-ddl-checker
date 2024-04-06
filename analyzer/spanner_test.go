package analyzer_test

import (
	"testing"

	"github.com/nerocrux/migration-ddl-checker/analyzer"
	"github.com/nerocrux/migration-ddl-checker/ddl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var query_Spanner_CREATE = []string{
	`CREATE TABLE users (id INT64) PRIMARY KEY (id);`,
	`CREATE INDEX SingersByFirstName ON Singers(FirstName);`,
	`ALTER TABLE Songwriters ADD COLUMN Nickname STRING(MAX) NOT NULL;`,
}

var query_Spanner_DROP = []string{
	`DROP TABLE users`,
	`DROP INDEX idx_users`,
	`ALTER TABLE Songwriters DROP COLUMN Nickname;`,
}

var query_Spanner_ALTER = []string{
	`ALTER TABLE Songwriters ALTER COLUMN Nickname STRING(MAX) NOT NULL;`,
}

func TestAnalyzeSpanner_CREATE(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		want     bool
	}{
		{
			name:     "create queries",
			contents: query_Spanner_CREATE,
			want:     true,
		},
		{
			name:     "drop queries",
			contents: query_Spanner_DROP,
			want:     false,
		},
		{
			name:     "alter queries",
			contents: query_Spanner_ALTER,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := analyzer.NewSpannerAnalyzer(ddl.FromConfig([]string{"CREATE"}))
			for _, q := range tt.contents {
				got, err := a.Analyze(q)
				require.NoError(t, err, q)
				assert.Equal(t, tt.want, got, q)
			}
		})
	}
}
