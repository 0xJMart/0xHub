package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

var (
	_ ports.TagRepository = (*Store)(nil)
)

// ListTags returns all tags sorted by display name.
func (s *Store) ListTags(ctx context.Context) ([]domain.Tag, error) {
	const query = `
SELECT
    id,
    name,
    display_name,
    color,
    created_at,
    updated_at
FROM tags
ORDER BY display_name ASC
`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres: list tags query: %w", err)
	}
	defer rows.Close()

	type tagRow struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		Display   string    `db:"display_name"`
		Color     *string   `db:"color"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	tagRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[tagRow])
	if err != nil {
		return nil, fmt.Errorf("postgres: list tags scan: %w", err)
	}

	tags := make([]domain.Tag, 0, len(tagRows))
	for _, row := range tagRows {
		tags = append(tags, domain.Tag{
			ID:          row.ID,
			Name:        row.Name,
			DisplayName: row.Display,
			Color:       stringOrDefault(row.Color, ""),
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return tags, nil
}
