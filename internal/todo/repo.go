package todo

import (
	"context"
	"database/sql"
	"time"
)

type Todo struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Repo struct {
	DB *sql.DB
}

func (r Repo) List(ctx context.Context, limit int) ([]Todo, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, title, created_at
		FROM todos
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Todo, 0, limit)
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}