package topics

import (
	"database/sql"

	"github.com/lib/pq"
)

func (a *App) dbListTopics(limit int, offset int, orderBy string, order string) (*[]Topic, error) {
	var rows *sql.Rows
	topics := []Topic{}

	rows, err := a.DB.Query(`
	SELECT 
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.user_id,
		u.username,
		u.full_name,
		array_agg(c.category_id),
		array_agg(c.name)
	FROM
		Topic t LEFT JOIN Topic_Category tc ON t.topic_id = tc.topic_id LEFT JOIN Category c ON c.category_id = tc.category_id LEFT JOIN KUser as u ON u.user_id = t.created_by
	GROUP BY t.topic_id, u.user_id
	ORDER BY t.`+orderBy+" "+order+`
	LIMIT $1 OFFSET $2
	`, limit, offset)
	if err == sql.ErrNoRows {
		return &topics, nil
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		t := Topic{}
		var cids []sql.NullInt64
		var cnames []sql.NullString
		err := rows.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.ID, &t.CreatedBy.Username, &t.CreatedBy.FullName, pq.Array(&cids), pq.Array(&cnames))
		if err != nil {
			return nil, err
		}
		for i := range cids {
			if cids[i].Valid && cnames[i].Valid {
				var c Category
				c.ID = cids[i].Int64
				c.Name = cnames[i].String
				t.Categories = append(t.Categories, c)
			}
		}
		topics = append(topics, t)
	}
	return &topics, nil
}

func dbGet(db *sql.DB, topicID int64) (interface{}, error) {
	return nil, nil
}
