package sqlxx_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestAssociations_GetAssociationQueries(t *testing.T) {
	article := Article{ID: 1}

	schema, err := sqlxx.GetSchema(article)
	assert.Nil(t, err)

	var assocs []sqlxx.Field
	for k, v := range schema.Associations {
		if len(k) <= 2 {
			assocs = append(assocs, v)
		}
	}

	queries, err := sqlxx.GetAssociationQueries(article, assocs)
	assert.Nil(t, err)

	results := []struct {
		field string
		query string
	}{
		{"Author", "SELECT articles.author_id, articles.created_at, articles.id, articles.is_published, articles.reviewer_id, articles.title, articles.updated_at FROM articles WHERE articles.author_id = ? LIMIT 1"},
		{"Reviewer", "SELECT articles.author_id, articles.created_at, articles.id, articles.is_published, articles.reviewer_id, articles.title, articles.updated_at FROM articles WHERE articles.reviewer_id = ? LIMIT 1"},
	}

	for _, tt := range results {
		for _, q := range queries {
			if q.Field.Name == tt.field {
				assert.Equal(t, tt.query, q.Query)
			}
		}
	}
}
