package users

import "database/sql"

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
	var u User
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
	var fus []FollowUser
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
	var fus []FollowUser
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
