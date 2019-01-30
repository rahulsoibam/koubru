package categories

import "database/sql"

func (a *App) dbListTopics(query string, limit int, offset int) (*[]Category, error) {
	categories := []Category{}
	var rows *sql.Rows
	var err error
	if query == "" {
		rows, err = a.DB.Query("SELECT category_id, name FROM Category ORDER BY created_on LIMIT = $1 OFFSET = $2", limit, offset)
	} else {
		query = query + "%"
		rows, err = a.DB.Query("SELECT category_id, name FROM Category WHERE name LIKE $1 ORDER BY created ON LIMIT=$2 OFFSET=$3", query, limit, offset)
	}
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
