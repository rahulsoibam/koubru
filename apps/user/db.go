package user

import (
	"database/sql"
)

func dbGetUserByID(db *sql.DB, userID int64) (*User, error) {
	var u User
	err := db.QueryRow("SELECT user_id, username, full_name, email_verified, photo_url, bio FROM KUser WHERE user_id=$1", userID).Scan(&u.ID, &u.Username, &u.FullName, &u.EmailVerfied, &u.PhotoURL, &u.Bio)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`
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

func dbGetFollowingByID(db *sql.DB, userID int64) (*[]FollowUser, error) {
	var fus []FollowUser
	rows, err := db.Query(`
	SELECT 
		u.user_id, u.username, u.full_name, u.photo_url, map.followed_on
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

func dbGetFollowersByID(db *sql.DB, userID int64) (*[]FollowUser, error) {
	var fus []FollowUser
	rows, err := db.Query(`
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
