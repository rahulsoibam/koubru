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

func (a *App) dbAuthenticatedGetUser(userID int64, quserID int64) (*User, error) {
	u := User{}
	err := a.DB.QueryRow(`
	SELECT 
		user_id, 
		username, 
		full_name, 
		email_verified, 
		photo_url, 
		bio,
		CASE WHEN EXISTS (SELECT 1 FROM UserMap AS map WHERE map.user_id = $1 AND map.follower_id=$2) THEN 1 ELSE 0 END AS is_following 
		FROM KUser 
		WHERE user_id=$1
	`, quserID, userID).Scan(&u.ID, &u.Username, &u.FullName, &u.EmailVerfied, &u.PhotoURL, &u.Bio, &u.IsFollowing)
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

func (a *App) dbGetUser(userID int64) (*User, error) {
	u := User{}
	err := a.DB.QueryRow("SELECT user_id, username, full_name, photo_url, bio FROM KUser WHERE user_id=$1", userID).Scan(&u.ID, &u.Username, &u.FullName, &u.PhotoURL, &u.Bio)
	if err != nil {
		return nil, err
	}

	u.IsFollowing = false

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
func (a *App) dbAuthenticatedGetFollowing(userID int64) (*[]FollowUser, error) {
	fus := []FollowUser{}
	rows, err := a.DB.Query(`
	SELECT 
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
		fu.IsFollowing = true
		if err := rows.Scan(&fu.Username, &fu.FullName, &fu.PhotoURL, &fu.FollowedOn); err != nil {
			return nil, err
		}
		fus = append(fus, fu)
	}
	return &fus, nil
}

func (a *App) dbGetFollowingByID(userID int64) (*[]FollowUser, error) {
	fus := []FollowUser{}
	rows, err := a.DB.Query(`
	SELECT 
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
		if err := rows.Scan(&fu.Username, &fu.FullName, &fu.PhotoURL, &fu.FollowedOn); err != nil {
			return nil, err
		}
		fus = append(fus, fu)
	}
	return &fus, nil
}

func (a *App) dbAuthenticatedGetFollowersSelf(userID int64) (*[]FollowUser, error) {
	fs := []FollowUser{}
	rows, err := a.DB.Query(`
	select u.username, u.full_name, u.photo_url, map.followed_on, case when following.user_id is null then 0 else 1 end as is_following
	from kuser u inner join usermap map on u.user_id = map.follower_id left join usermap following on following.user_id=map.follower_id AND following.follower_id=map.user_id
	where map.user_id=$1
	order by following.followed_on desc nulls last;
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fs, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var f FollowUser
		if err := rows.Scan(&f.Username, &f.FullName, &f.PhotoURL, &f.FollowedOn, &f.IsFollowing); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	return &fs, nil
}

func (a *App) dbAuthenticatedGetFollowers(userID int64, quserID int64) (*[]FollowUser, error) {
	fs := []FollowUser{}
	rows, err := a.DB.Query(`
	select u.username, u.full_name, u.photo_url, map.followed_on, case when following.user_id = $1 then 1 else 0 end as is_following
    from kuser u inner join usermap map on u.user_id = map.follower_id full join usermap following on following.user_id=map.follower_id AND following.follower_id=map.user_id
    where map.user_id = $2 
    order by following.followed_on desc nulls last;
	`, userID, quserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fs, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		f := FollowUser{}
		if err := rows.Scan(&f.Username, &f.FullName, &f.PhotoURL, &f.FollowedOn, &f.IsFollowing); err != nil {
			return &fs, nil
		}
		return nil, err
	}
	return &fs, nil
}

func (a *App) dbGetFollowers(quserID int64) (*[]FollowUser, error) {
	fs := []FollowUser{}
	rows, err := a.DB.Query(`
	select u.username, u.full_name, u.photo_url, map.followed_on
	from usermap map inner join kuser u on u.user_id = map.follower_id
	where map.user_id = $1
	`, quserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fs, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		f := FollowUser{}
		if err := rows.Scan(&f.Username, &f.FullName, &f.PhotoURL, &f.FollowedOn); err != nil {
			return &fs, nil
		}
		return nil, err
	}
	return &fs, nil
}

func (a *App) dbAuthenticatedListTopics(userID int64, quserID int64, limit int, offset int, orderBy string, order string) (*[]Topic, error) {
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
	WHERE t.created_by = $2
	GROUP BY t.topic_id, u.user_id
	ORDER BY t.`+orderBy+` `+order+`
	LIMIT $3 OFFSET $4
	`, userID, quserID, limit, offset)
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

func (a *App) dbListTopics(quserID int64, limit int, offset int, orderBy string, order string) (*[]Topic, error) {
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
	WHERE t.created_by = $1
	GROUP BY t.topic_id, u.user_id
	ORDER BY t.`+orderBy+` `+order+`
	LIMIT $2 OFFSET $3
	`, quserID, limit, offset)
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
