package topics

import (
	"database/sql"

	"github.com/lib/pq"
)

func (a *App) dbAuthenticatedListTopics(userID int64, limit int, offset int, orderBy string, order string) (*[]Topic, error) {
	var rows *sql.Rows

	topics := []Topic{}

	rows, err := a.DB.Query(`
	SELECT 
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.photo_url,
		array_agg(c.category_id),
		array_agg(c.name),
		CASE WHEN EXISTS (SELECT 1 FROM topic_follower AS tf WHERE tf.topic_id = t.topic_id AND tf.followed_by=$1) THEN 1 ELSE 0 END AS is_following,
		COUNT(DISTINCT tf.topic_id)
	FROM
		Topic t LEFT JOIN Topic_Category tc ON t.topic_id = tc.topic_id LEFT JOIN Category c ON c.category_id = tc.category_id LEFT JOIN KUser as u ON u.user_id = t.created_by LEFT JOIN Topic_Follower tf ON tf.topic_id = t.topic_id
	GROUP BY t.topic_id, u.user_id
	ORDER BY t.`+orderBy+` `+order+`
	LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err == sql.ErrNoRows {
		return &topics, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		t := Topic{}
		t.Categories = []Category{}
		var cids []sql.NullInt64
		var cnames []sql.NullString
		err := rows.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, pq.Array(&cids), pq.Array(&cnames), &t.IsFollowing, &t.Counts.Followers)
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

func (a *App) dbListTopics(limit int, offset int, orderBy string, order string) (*[]Topic, error) {
	var rows *sql.Rows

	topics := []Topic{}

	rows, err := a.DB.Query(`
	SELECT 
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.photo_url,
		array_agg(c.category_id),
		array_agg(c.name),
		COUNT(DISTINCT tf.topic_id)
	FROM
		Topic t LEFT JOIN Topic_Category tc ON t.topic_id = tc.topic_id LEFT JOIN Category c ON c.category_id = tc.category_id LEFT JOIN KUser as u ON u.user_id = t.created_by LEFT JOIN Topic_Follower tf ON tf.topic_id = t.topic_id
	GROUP BY t.topic_id, u.user_id
	ORDER BY t.`+orderBy+` `+order+`
	LIMIT $1 OFFSET $2
	`, limit, offset)
	if err == sql.ErrNoRows {
		return &topics, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		t := Topic{}
		t.Categories = []Category{}
		var cids []sql.NullInt64
		var cnames []sql.NullString
		err := rows.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, pq.Array(&cids), pq.Array(&cnames), &t.Counts.Followers)
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

func (a *App) dbAuthenticatedGetTopicByID(userID int64, topicID int64) (*Topic, error) {
	t := Topic{}
	t.Categories = []Category{}
	row := a.DB.QueryRow(`
	SELECT
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.photo_url,
		array_agg(c.category_id),
		array_agg(c.name),
		CASE WHEN EXISTS (SELECT 1 FROM topic_follower AS tf WHERE tf.topic_id = t.topic_id AND tf.followed_by=$1) THEN 1 ELSE 0 END AS is_following
	FROM
		Topic t INNER JOIN Topic_Category tc USING (topic_id) INNER JOIN Category c USING(category_id) INNER JOIN KUser as u ON t.created_by=u.user_id
	WHERE t.topic_id=$1
	GROUP BY t.topic_id, u.user_id
	`, topicID)

	var cids []sql.NullInt64
	var cnames []sql.NullString

	err := row.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, pq.Array(&cids), pq.Array(&cnames), &t.IsFollowing)
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

	err = a.DB.QueryRow("SELECT COUNT(*) FROM Topic_Follower WHERE topic_id=$1", topicID).Scan(&t.Counts.Followers)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
func (a *App) dbGetTopicByID(topicID int64) (*Topic, error) {
	t := Topic{}
	t.Categories = []Category{}
	row := a.DB.QueryRow(`
	SELECT
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.photo_url,
		array_agg(c.category_id),
		array_agg(c.name),
	FROM
		Topic t INNER JOIN Topic_Category tc USING (topic_id) INNER JOIN Category c USING(category_id) INNER JOIN KUser as u ON t.created_by=u.user_id
	WHERE t.topic_id=$1
	GROUP BY t.topic_id, u.user_id
	`, topicID)

	var cids []sql.NullInt64
	var cnames []sql.NullString

	err := row.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, pq.Array(&cids), pq.Array(&cnames))
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

	err = a.DB.QueryRow("SELECT COUNT(*) FROM Topic_Follower WHERE topic_id=$1", topicID).Scan(&t.Counts.Followers)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
func (a *App) dbCreateTopic(nt *NewTopic, userID int64) (*Topic, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return nil, err
	}
	var newTopicID int64
	err = tx.QueryRow("INSERT INTO Topic (title, created_by) VALUES ($1, $2) RETURNING topic_id", nt.Title, userID).Scan(&newTopicID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range nt.Categories {
		_, err = tx.Exec("INSERT INTO Topic_Category (topic_id, category_id) VALUES ($1, $2)", nt.Categories[i].ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	_, err = tx.Exec("INSERT INTO Topic_Follower (topic_id, followed_by) VALUES ($1, $2)", newTopicID, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	topic, err := a.dbAuthenticatedGetTopicByID(userID, newTopicID)
	if err != nil {
		return nil, err
	}
	return topic, nil
}
