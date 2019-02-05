package user

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru/types"
)

func (a *App) AuthGetQuery(userID int64) (types.User, error) {
	u := types.User{}
	sqlQuery := `
	SELECT 
		u.username, 
		u.full_name, 
		u.picture, 
		u.bio,
		1 as is_self,
		0 as is_following,
		(select count(*) from user_follower where user_id=u.user_id) as followers_count,
		(select count(*) from user_follower where follower_id=u.user_id) as following_count,
		(select count(*) from topic where creator_id=u.user_id) as topic_count,
		(select count(*) from opinion where creator_id=u.user_id) as opinion_count
	FROM KUser u
	WHERE user_id=$1
	`
	err := a.DB.QueryRow(sqlQuery, userID).Scan(&u.Username, &u.FullName, &u.Picture, &u.Bio, &u.IsSelf, &u.IsFollowing, &u.Counts.Followers, &u.Counts.Following, &u.Counts.Topics, &u.Counts.Opinions)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (a *App) AuthFollowingQuery(userID int64) ([]types.Following, error) {
	fs := []types.Following{}
	sqlQuery := `
	SELECT
		u.username,
		u.full_name,
		u.picture,
		uf.followed_on,
		1 as is_following,
		0 as is_self
	FROM KUser u INNER JOIN user_follower uf on u.user_id=uf.user_id
	WHERE uf.follower_id=$1
	ORDER BY uf.followed_on DESC
	`

	rows, err := a.DB.Query(sqlQuery, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fs, nil
		}
		return fs, err
	}

	defer rows.Close()
	for rows.Next() {
		f := types.Following{}
		if err := rows.Scan(&f.Username, &f.FullName, &f.Picture, &f.FollowedOn, &f.IsFollowing, &f.IsSelf); err != nil {
			return fs, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}

func (a *App) AuthFollowersQuery(userID int64) ([]types.Follower, error) {
	fs := []types.Follower{}
	sqlQuery := `
	SELECT
        u.username,
        u.full_name,
		u.picture,
		uf.followed_on,
        case when exists(select 1 from user_follower where user_id=u.user_id and follower_id=uf.user_id) then 1 else 0 end as is_following,
		0 as is_self
    FROM
        Kuser u inner join user_follower uf on u.user_id=uf.follower_id
    WHERE uf.user_id=$1
    ORDER BY (select followed_on from user_follower where user_id=u.user_id and follower_id=uf.user_id) DESC, uf.followed_on DESC;
	`

	rows, err := a.DB.Query(sqlQuery, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fs, nil
		}
		return fs, err
	}

	defer rows.Close()
	for rows.Next() {
		f := types.Follower{}
		if err := rows.Scan(&f.Username, &f.FullName, &f.Picture, &f.FollowedOn, &f.IsFollowing, &f.IsSelf); err != nil {
			return fs, err
		}
		fs = append(fs, f)
	}
	err = rows.Err()
	if err != nil {
		return fs, err
	}

	return fs, nil
}

func (a *App) AuthTopicsQuery(userID int64) ([]types.Topic_, error) {
	ts := []types.Topic_{}
	sqlQuery := `
	SELECT
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.picture,
		coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) FILTER (WHERE c.category_id IS NOT NULL OR c.name IS NOT NULL), '[]'::json) as categories,
		CASE WHEN EXISTS (SELECT 1 FROM topic_follower WHERE topic_id=t.topic_id AND follower_id=$1) THEN 1 ELSE 0 END AS is_following
	FROM topic t inner join kuser u on t.creator_id=u.user_id left join topic_category tc on t.topic_id = tc.topic_id
	LEFT JOIN category c on c.category_id=tc.category_id
	WHERE t.creator_id=$1
	GROUP BY t.topic_id, u.user_id
	ORDER BY t.created_on desc
	`

	rows, err := a.DB.Query(sqlQuery, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ts, nil
		}
		return ts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := types.Topic_{}
		err := rows.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, (*[]byte)(&t.Categories), &t.IsFollowing)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}
	err = rows.Err()
	if err != nil {
		return ts, err
	}
	return ts, nil
}

func (a *App) AuthOpinionsQuery(userID int64) ([]types.Opinion, error) {
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
        case when exists(select 1 from topic_follower where topic_id=o.topic_id and follower_id=o.creator_id) then 1 else 0 end as is_following_topic,
        o.is_anonymous,
        case when exists(select 1 from opinion_follower where opinion_id=o.opinion_id and follower_id=o.creator_id) then 1 else 0 end as is_following,
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
        LEFT JOIN Opinion_Vote ov on ov.opinion_id=o.opinion_id and voter_id=o.creator_id
        LEFT JOIN Topic_Category tc on tc.topic_id = t.topic_id
        LEFT JOIN Category c on c.category_id=tc.category_id
    WHERE o.creator_id=$1
	GROUP BY o.opinion_id, u.user_id, t.topic_id, views, ov.vote
	ORDER BY o.created_on
	`

	rows, err := a.DB.Query(sqlQuery, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return os, err
		}
		return os, err
	}

	defer rows.Close()
	for rows.Next() {
		o := types.Opinion{}
		err := rows.Scan(&o.ID, &o.CreatedBy.Username, &o.CreatedBy.FullName, &o.CreatedBy.Picture, &o.CreatedBy.IsSelf, &o.Topic.ID, &o.Topic.Title, &o.Topic.Details, (*[]byte)(&o.Topic.Categories), &o.Topic.IsFollowing, &o.IsAnonymous, &o.IsFollowing, &o.IsFollowing, pq.Array(&o.Thumbnails), &o.Sources.Hls, &o.Sources.Dash, &o.Sources.Aac, &o.Vote, o.Reaction, &o.CreatedOn, &o.Counts.Views, &o.Counts.Upvotes, &o.Counts.Downvotes, &o.Counts.Followers, &o.Counts.Replies)
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
