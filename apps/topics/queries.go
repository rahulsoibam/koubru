package topics

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rahulsoibam/koubru-prod-api/types"
)

func (a *App) AuthCreateQuery(userID int64, t types.NewTopic) (types.Topic, error) {
	tres := types.Topic{}
	tx, err := a.DB.Begin()
	if err != nil {
		return tres, err
	}
	var topicID int64
	err = tx.QueryRow("INSERT INTO Topic (title, details, created_by) VALUES ($1, $2, $3) RETURNING topic_id", t.Title, t.Details, userID).Scan(&topicID)
	if err != nil {
		tx.Rollback()
		return tres, err
	}

	for i := 0; i < 3; i++ {
		fmt.Println(t.Categories[i])
		_, err = tx.Exec("INSERT INTO Topic_Category (topic_id, category_id) VALUES ($1, $2)", topicID, t.Categories[i])
		if err != nil {
			tx.Rollback()
			return tres, err
		}
	}
	_, err = tx.Exec("INSERT INTO Topic_Follower (topic_id, followed_by) VALUES ($1, $2)", topicID, userID)
	if err != nil {
		tx.Rollback()
		return tres, err
	}
	err = tx.Commit()
	if err != nil {
		return tres, err
	}

	tres, err = a.AuthGetQuery(userID, topicID)
	if err != nil {
		return tres, err
	}

	// TODO return topic page
	return tres, nil
}

func (a *App) AuthGetQuery(userID int64, topicID int64) (types.Topic, error) {
	t := types.Topic{}
	sqlQuery := `
	SELECT
        t.topic_id,
        t.title,
        t.details,
        t.created_on,
        u.username,
        u.full_name,
        u.photo_url,
        coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null), '[]'::json),
        case when exists (select 1 from topic_follower where topic_id=t.topic_id and followed_by=$1) then 1 else 0 end as is_following
    FROM topic t inner join kuser u on t.created_by=u.user_id left join topic_category tc on t.topic_id = tc.topic_id
    left join category c on c.category_id=tc.category_id
    where t.topic_id=$2
    group by t.topic_id, u.user_id
	`

	err := a.DB.QueryRow(sqlQuery, userID, topicID).Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, (*[]byte)(&t.Categories), &t.IsFollowing)
	if err != nil {
		// check for sql.ErrNoRows and return 404 if that is the case
		return t, err
	}

	err = a.DB.QueryRow(`
	SELECT COUNT(*)
	FROM topic_follower tf
	WHERE tf.topic_id=$1
	`, topicID).Scan(&t.Counts.Followers)
	if err != nil {
		return t, err
	}

	err = a.DB.QueryRow(`
	SELECT COUNT(*)
	FROM Opinion o
	WHERE o.topic_id=$1
	`, topicID).Scan(&t.Counts.Opinions)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (a *App) GetQuery(topicID int64) (types.Topic, error) {
	t := types.Topic{}
	sqlQuery := `
	SELECT
        t.topic_id,
        t.title,
        t.details,
        t.created_on,
        u.username,
        u.full_name,
        u.photo_url,
		coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null), '[]'::json),
		0 as is_following
    FROM topic t inner join kuser u on t.created_by=u.user_id left join topic_category tc on t.topic_id = tc.topic_id
    left join category c on c.category_id=tc.category_id
    where t.topic_id=$1
    group by t.topic_id, u.user_id
	`

	err := a.DB.QueryRow(sqlQuery, topicID).Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, (*[]byte)(&t.Categories), &t.IsFollowing)
	if t.Categories == nil {
		t.Categories = json.RawMessage("[]")
	}
	if err != nil {
		// check for sql.ErrNoRows and return 404 if that is the case
		return t, err
	}

	err = a.DB.QueryRow(`
	SELECT COUNT(*)
	FROM topic_follower tf
	WHERE tf.topic_id=$1
	`, topicID).Scan(&t.Counts.Followers)
	if err != nil {
		return t, err
	}

	err = a.DB.QueryRow(`
	SELECT COUNT(*)
	FROM Opinion o
	WHERE o.topic_id=$1
	`, topicID).Scan(&t.Counts.Opinions)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (a *App) AuthFollowersQuery(userID int64, topicID int64) ([]types.User_, error) {
	fs := []types.User_{}

	sqlQuery := `
	SELECT
		u.username,
		u.full_name,
		u.photo_url,
		CASE WHEN tf.followed_by=$1 THEN 1 ELSE 0 END AS is_self,
		CASE WHEN EXISTS (SELECT 1 FROM Usermap map WHERE map.user_id=u.user_id AND map.follower_id=$1) THEN 1 ELSE 0 END AS is_following
	FROM
		KUser u INNER JOIN Topic_Follower tf ON u.user_id=tf.followed_by
	WHERE tf.topic_id=$2
	ORDER BY is_self desc, is_following desc
	`

	rows, err := a.DB.Query(sqlQuery, userID, topicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fs, nil
		}
		return fs, err
	}

	defer rows.Close()
	for rows.Next() {
		f := types.User_{}
		err := rows.Scan(&f.Username, &f.FullName, &f.Picture, &f.IsSelf, &f.IsFollowing)
		if err != nil {
			return fs, nil
		}
		fs = append(fs, f)
	}

	err = rows.Err()
	if err != nil {
		return fs, err
	}

	return fs, nil
}

func (a *App) FollowersQuery(topicID int64) ([]types.User_, error) {
	fs := []types.User_{}

	sqlQuery := `
	SELECT
		u.username,
		u.full_name,
		u.photo_url,
		0 as is_self,
		0 as is_following
	FROM
		KUser u INNER JOIN Topic_Follower tf ON u.user_id=tf.followed_by left join usermap map on u.user_id=map.user_id
	WHERE tf.topic_id=$1
	GROUP BY u.user_id
	ORDER BY (SELECT COUNT(map.follower_id)) DESC
	`

	rows, err := a.DB.Query(sqlQuery, topicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fs, nil
		}
		return fs, err
	}

	defer rows.Close()
	for rows.Next() {
		f := types.User_{}
		err := rows.Scan(&f.Username, &f.FullName, &f.Picture, &f.IsSelf, &f.IsFollowing)
		if err != nil {
			return fs, nil
		}
		fs = append(fs, f)
	}

	err = rows.Err()
	if err != nil {
		return fs, err
	}

	return fs, nil
}

func (a *App) AuthOpinionsQuery(userID int64, topicID int64) ([]types.Opinion_, error) {

}

func (a *App) OpinionsQuery(topicID int64) ([]types.Opinion_, error) {

}
