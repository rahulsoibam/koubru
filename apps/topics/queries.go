package topics

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"github.com/rahulsoibam/koubru/types"
)

func (a *App) AuthCreateQuery(userID int64, t types.NewTopic) (types.Topic, error) {
	tres := types.Topic{}
	tx, err := a.DB.Begin()
	if err != nil {
		return tres, err
	}
	var topicID int64
	err = tx.QueryRow("INSERT INTO Topic (title, details, creator_by) VALUES ($1, $2, $3) RETURNING topic_id", t.Title, t.Details, userID).Scan(&topicID)
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
	_, err = tx.Exec("INSERT INTO Topic_Follower (topic_id, follower_id) VALUES ($1, $2)", topicID, userID)
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
        u.username,
        u.full_name,
        u.picture,
        case when t.creator_id=$1 then 1 else 0 end as is_self,
        coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null), '[]'::json),
		case when exists (select 1 from topic_follower where topic_id=t.topic_id and follower_id=$1) then 1 else 0 end as is_following,
		t.created_on,
        (select count(*) from topic_follower where topic_id=t.topic_id) as followers_count,
        (select count(*) from opinion where topic_id=t.topic_id) as opinions_count
    FROM
        topic t inner join kuser u on t.creator_id=u.user_id
        left join topic_category tc on t.topic_id = tc.topic_id
        left join category c on c.category_id=tc.category_id
    where t.topic_id=$2
    group by t.topic_id, u.user_id
	
	`

	err := a.DB.QueryRow(sqlQuery, userID, topicID).Scan(&t.ID, &t.Title, &t.Details, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, &t.CreatedBy.IsSelf, (*[]byte)(&t.Categories), &t.IsFollowing, &t.CreatedOn, &t.Counts.Followers, &t.Counts.Opinions)
	if err != nil {
		// check for sql.ErrNoRows and return 404 if that is the case
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
        u.username,
        u.full_name,
        u.picture,
        case when t.creator_id=$1 then 1 else 0 end as is_self,
        coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null), '[]'::json),
		0 as is_following,
		t.created_on,
        (select count(*) from topic_follower where topic_id=t.topic_id) as followers_count,
        (select count(*) from opinion where topic_id=t.topic_id) as opinions_count
    FROM
        topic t inner join kuser u on t.creator_id=u.user_id
        left join topic_category tc on t.topic_id = tc.topic_id
        left join category c on c.category_id=tc.category_id
    where t.topic_id=$1
    group by t.topic_id, u.user_id
	`

	err := a.DB.QueryRow(sqlQuery, topicID).Scan(&t.ID, &t.Title, &t.Details, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, &t.CreatedBy.IsSelf, (*[]byte)(&t.Categories), &t.IsFollowing, &t.CreatedOn, &t.Counts.Followers, &t.Counts.Opinions)
	if t.Categories == nil {
		t.Categories = json.RawMessage("[]")
	}
	if err != nil {
		// check for sql.ErrNoRows and return 404 if that is the case
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
        u.picture,
        CASE WHEN tf.follower_id=$1 THEN 1 ELSE 0 END AS is_self,
        CASE WHEN EXISTS (SELECT 1 FROM User_follower WHERE user_id=u.user_id AND follower_id=$1) THEN 1 ELSE 0 END AS is_following
    FROM
        KUser u INNER JOIN Topic_Follower tf ON u.user_id=tf.follower_id
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
        u.picture,
        0 AS is_self,
        0 AS is_following
    FROM
        KUser u INNER JOIN Topic_Follower tf ON u.user_id=tf.follower_id
    WHERE tf.topic_id=$1
    ORDER BY is_self desc, is_following desc
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

func (a *App) AuthOpinionsQuery(userID int64, topicID int64) ([]types.Opinion, error) {
	os := []types.Opinion{}

	sqlQuery := `
	SELECT
        o.opinion_id,
        u.username,
        u.full_name,
        u.picture,
        case when o.creator_id=$1 then 1 else 0 end as is_self,
        t.topic_id,
        t.title,
        t.details,
        coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null), '[]'::json) as categories,
        case when exists(select 1 from topic_follower where topic_id=o.topic_id and follower_id=$1) then 1 else 0 end as is_following_topic,
        o.is_anonymous,
        case when exists(select 1 from opinion_follower where opinion_id=o.opinion_id and follower_id=$1) then 1 else 0 end as is_following,
        o.thumbnails,
        case when o.is_anonymous then '_blank' else o.hls end as hls,
        case when o.is_anonymous then '_blank' else o.dash end as dash,
        case when o.is_anonymous then o.aac else '_blank' end as aac,
        case when ov.vote=true then 'upvote' when ov.vote=false then 'downvote' else 'none' end as vote,
        o.reaction,
        o.created_on,
        (select count(*) from opinion_view where opinion_id=o.opinion_id) as views,
        (select count(*) from opinion_vote where opinion_id=o.opinion_id and vote=true) as upvotes,
        (select count(*) from opinion_vote where opinion_id=o.opinion_id and vote=false) as downvotes,
        (select count(*) from opinion_follower where opinion_id=o.opinion_id) as followers,
        (select count(*) from opinion where parent_id=o.opinion_id) as replies
    FROM
        opinion o INNER JOIN topic t USING(topic_id)
        INNER JOIN kuser u on o.creator_id=u.user_id
        LEFT JOIN Opinion_Vote ov on ov.opinion_id=o.opinion_id and voter_id=1967600534613394434
        LEFT JOIN Topic_Category tc on tc.topic_id = t.topic_id
        LEFT JOIN Category c on c.category_id=tc.category_id
    WHERE o.topic_id=$2
    GROUP BY o.opinion_id, u.user_id, t.topic_id, views, ov.vote
	`

	rows, err := a.DB.Query(sqlQuery, userID, topicID)
	if err != nil {
		if err != sql.ErrNoRows {
			return os, err
		}
		return os, err
	}

	defer rows.Close()
	for rows.Next() {
		o := types.Opinion{}
		err := rows.Scan(&o.ID, &o.CreatedBy.Username, &o.CreatedBy.FullName, &o.CreatedBy.Picture, &o.CreatedBy.IsSelf, &o.Topic.ID, &o.Topic.Title, &o.Topic.Details, (*[]byte)(&o.Topic.Categories), &o.Topic.IsFollowing, &o.IsAnonymous, &o.IsFollowing, pq.Array(&o.Thumbnails), &o.Sources.Hls, &o.Sources.Dash, &o.Sources.Aac, &o.Vote, &o.Reaction, &o.CreatedOn, &o.Counts.Views, &o.Counts.Upvotes, &o.Counts.Downvotes, &o.Counts.Followers, &o.Counts.Replies)
		if err != nil {
			return os, err
		}
		os = append(os, o)
	}

	err = rows.Err()
	if err != nil {
		return os, err
	}

	return os, nil
}

func (a *App) OpinionsQuery(topicID int64) ([]types.Opinion, error) {
	os := []types.Opinion{}

	sqlQuery := `
	SELECT
        o.opinion_id,
        u.username,
        u.full_name,
        u.picture,
        0 as is_self,
        t.topic_id,
        t.title,
        t.details,
        coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null), '[]'::json) as categories,
        0 is_following_topic,
        o.is_anonymous,
        0 as is_following,
        o.thumbnails,
        case when o.is_anonymous then '_blank' else o.hls end as hls,
        case when o.is_anonymous then '_blank' else o.dash end as dash,
        case when o.is_anonymous then o.aac else '_blank' end as aac,
        'none',
        o.reaction,
        o.created_on,
        (select count(*) from opinion_view where opinion_id=o.opinion_id) as views,
        (select count(*) from opinion_vote where opinion_id=o.opinion_id and vote=true) as upvotes,
        (select count(*) from opinion_vote where opinion_id=o.opinion_id and vote=false) as downvotes,
        (select count(*) from opinion_follower where opinion_id=o.opinion_id) as followers,
        (select count(*) from opinion where parent_id=o.opinion_id) as replies
    FROM
        opinion o INNER JOIN topic t USING(topic_id)
        INNER JOIN kuser u on o.creator_id=u.user_id
        LEFT JOIN Topic_Category tc on tc.topic_id = t.topic_id
        LEFT JOIN Category c on c.category_id=tc.category_id
    WHERE o.topic_id=$1
    GROUP BY o.opinion_id, u.user_id, t.topic_id, views
	`

	rows, err := a.DB.Query(sqlQuery, topicID)
	if err != nil {
		if err != sql.ErrNoRows {
			return os, err
		}
		return os, err
	}

	defer rows.Close()
	for rows.Next() {
		o := types.Opinion{}
		err := rows.Scan(&o.ID, &o.CreatedBy.Username, &o.CreatedBy.FullName, &o.CreatedBy.Picture, &o.CreatedBy.IsSelf, &o.Topic.ID, &o.Topic.Title, &o.Topic.Details, (*[]byte)(&o.Topic.Categories), &o.Topic.IsFollowing, &o.IsAnonymous, &o.IsFollowing, pq.Array(&o.Thumbnails), &o.Sources.Hls, &o.Sources.Dash, &o.Sources.Aac, &o.Vote, &o.Reaction, &o.CreatedOn, &o.Counts.Views, &o.Counts.Upvotes, &o.Counts.Downvotes, &o.Counts.Followers, &o.Counts.Replies)
		if err != nil {
			return os, err
		}
		os = append(os, o)
	}

	err = rows.Err()
	if err != nil {
		return os, err
	}

	return os, nil
}
