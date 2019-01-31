package user

import (
	"database/sql"

	"github.com/lib/pq"
)

func (a *App) dbGetUserIDUsingUsername(username string) (int64, error) {
	var userID int64
	var err error
	err = a.DB.QueryRow("SELECT user_id from KUser WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *App) dbGetUserByID(userID int64) (*User, error) {
	u := User{}
	err := a.DB.QueryRow("SELECT user_id, username, full_name, email_verified, photo_url, bio FROM KUser WHERE user_id=$1", userID).Scan(&u.ID, &u.Username, &u.FullName, &u.EmailVerfied, &u.PhotoURL, &u.Bio)
	if err != nil {
		return nil, err
	}

	err = a.DB.QueryRow(`
	SELECT 
		count(*) FILTER (WHERE user_id=$1) as followers,
		count(*) FILTER (WHERE follower_id=$1) as following
	FROM UserMap
	`, userID).Scan(&u.Counts.Followers, &u.Counts.Following)
	if err != nil {
		return nil, err
	}
	// TODO Add topic and opinion count when their tables are created
	return &u, nil
}

func (a *App) dbGetFollowingByID(userID int64) (*[]FollowUser, error) {
	fus := []FollowUser{}
	rows, err := a.DB.Query(`
	SELECT 
		u.user_id, 
		u.username, 
		u.full_name, 
		u.photo_url, 
		map.followed_on 
	FROM
		KUser AS u INNER JOIN UserMap AS map USING (user_id)
	WHERE map.follower_id = $1
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fus, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var fu FollowUser
		if err := rows.Scan(&fu.ID, &fu.Username, &fu.FullName, &fu.PhotoURL, &fu.FollowedOn); err != nil {
			return nil, err
		}
		fus = append(fus, fu)
	}
	return &fus, nil
}

func (a *App) dbGetFollowersByID(userID int64) (*[]FollowUser, error) {
	fus := []FollowUser{}
	rows, err := a.DB.Query(`
	SELECT 
		u.user_id, u.username, u.full_name, u.photo_url, map.followed_on
	FROM
		KUser AS u INNER JOIN UserMap AS map ON u.user_id = map.follower_id
	WHERE map.user_id = $1
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fus, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var fu FollowUser
		if err := rows.Scan(&fu.ID, &fu.Username, &fu.FullName, &fu.PhotoURL, &fu.FollowedOn); err != nil {
			return nil, err
		}
		fus = append(fus, fu)
	}
	return &fus, nil
}

func (a *App) dbListTopicsByUserID(userID int64, limit int, offset int, orderBy string, order string) (*[]Topic, error) {
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
	WHERE u.user_id=$1
		GROUP BY t.topic_id, u.user_id
	ORDER BY `+orderBy+" "+order+`
	LIMIT $2 OFFSET $3
	`, userID, limit, offset)
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
