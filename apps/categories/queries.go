package categories

import (
	"context"
	"database/sql"
)

func (a *App) ListQuery(ctx context.Context) ([]Category, error) {
	cs := []Category{}
	var rows *sql.Rows
	var err error
	query = query + "%"
	rows, err = a.DB.Query("SELECT category_id, name FROM Category WHERE name LIKE $1 ORDER BY created_on LIMIT $2 OFFSET $3", query, limit, offset)
	if err == sql.ErrNoRows {
		return &categories, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return &categories, nil
}
