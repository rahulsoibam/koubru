package opinions

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru/middleware"

	"github.com/rahulsoibam/koubru/types"
)

func (a *App) AuthListQuery(ctx context.Context, userID int64) ([]types.Opinion, error) {
	os := []types.Opinion{}
	limit := ctx.Value(middleware.PaginationKeys("per_page")).(int)
	offset := ctx.Value(middleware.PaginationKeys("db_offset")).(int)

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
        LEFT JOIN Opinion_Vote ov on ov.opinion_id=o.opinion_id and voter_id=$1
        LEFT JOIN Topic_Category tc on tc.topic_id = t.topic_id
        LEFT JOIN Category c on c.category_id=tc.category_id
	GROUP BY o.opinion_id, u.user_id, t.topic_id, views, ov.vote
	ORDER BY o.created_on DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := a.DB.Query(sqlQuery, userID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return os, nil
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
func (a *App) ListQuery(ctx context.Context) ([]types.Opinion, error) {
	os := []types.Opinion{}
	limit := ctx.Value(middleware.PaginationKeys("per_page")).(int)
	offset := ctx.Value(middleware.PaginationKeys("db_offset")).(int)

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
        0 as is_following_topic,
        o.is_anonymous,
        0 as is_following,
        o.thumbnails,
        case when o.is_anonymous then '_blank' else o.hls end as hls,
        case when o.is_anonymous then '_blank' else o.dash end as dash,
        case when o.is_anonymous then o.aac else '_blank' end as aac,
        'none' as vote,
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
	GROUP BY o.opinion_id, u.user_id, t.topic_id, views
	ORDER BY o.created_on DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := a.DB.Query(sqlQuery, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return os, nil
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

func (a *App) AuthGetQuery(userID int64, opinionID int64) (types.Opinion, error) {
	o := types.Opinion{}
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
        LEFT JOIN Opinion_Vote ov on ov.opinion_id=o.opinion_id and voter_id=$1
        LEFT JOIN Topic_Category tc on tc.topic_id = t.topic_id
		LEFT JOIN Category c on c.category_id=tc.category_id
	WHERE o.opinion_id=$2
	`

	err := a.DB.QueryRow(sqlQuery, userID, opinionID).Scan(&o.ID, &o.CreatedBy.Username, &o.CreatedBy.FullName, &o.CreatedBy.Picture, &o.CreatedBy.IsSelf, &o.Topic.ID, &o.Topic.Title, &o.Topic.Details, (*[]byte)(&o.Topic.Categories), &o.Topic.IsFollowing, &o.IsAnonymous, &o.IsFollowing, pq.Array(&o.Thumbnails), &o.Sources.Hls, &o.Sources.Dash, &o.Sources.Aac, &o.Vote, &o.Reaction, &o.CreatedOn, &o.Counts.Views, &o.Counts.Upvotes, &o.Counts.Downvotes, &o.Counts.Followers, &o.Counts.Replies)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (a *App) GetQuery(opinionID int64) (types.Opinion, error) {
	o := types.Opinion{}
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
        0 as is_following_topic,
        o.is_anonymous,
        0 as is_following,
        o.thumbnails,
        case when o.is_anonymous then '_blank' else o.hls end as hls,
        case when o.is_anonymous then '_blank' else o.dash end as dash,
        case when o.is_anonymous then o.aac else '_blank' end as aac,
        'none' as vote,
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
    WHERE o.opinion_id=$1
    GROUP BY o.opinion_id, u.user_id, t.topic_id, views;
	`

	err := a.DB.QueryRow(sqlQuery, opinionID).Scan(&o.ID, &o.CreatedBy.Username, &o.CreatedBy.FullName, &o.CreatedBy.Picture, &o.CreatedBy.IsSelf, &o.Topic.ID, &o.Topic.Title, &o.Topic.Details, (*[]byte)(&o.Topic.Categories), &o.Topic.IsFollowing, &o.IsAnonymous, &o.IsFollowing, pq.Array(&o.Thumbnails), &o.Sources.Hls, &o.Sources.Dash, &o.Sources.Aac, &o.Vote, &o.Reaction, &o.CreatedOn, &o.Counts.Views, &o.Counts.Upvotes, &o.Counts.Downvotes, &o.Counts.Followers, &o.Counts.Replies)

	if err != nil {
		return o, err
	}
	return o, nil
}

func (a *App) AuthCreateQuery(userID int64, no types.NewOpinion) (types.Opinion, error) {
	o := types.Opinion{}
	tx, err := a.DB.Begin()
	if err != nil {
		return o, err
	}

	var opinionID int64
	err = tx.QueryRow("INSERT INTO Opinion (topic_id, creator_id, reaction, is_anonymous, dash) VALUES ($1, $2, $3, $4, $5) RETURNING opinion_id", no.TopicID, userID, no.Reaction, no.IsAnonymous, no.Mp4).Scan(&opinionID)
	if err != nil {
		tx.Rollback()
		return o, err
	}

	_, err = tx.Exec("INSERT INTO opinion_follower (opinion_id, follower_id) VALUES ($1, $2)", opinionID, userID)
	if err != nil {
		tx.Rollback()
		return o, err
	}
	err = tx.Commit()
	if err != nil {
		return o, err
	}

	o, err = a.AuthGetQuery(userID, opinionID)
	if err != nil {
		return o, err
	}
	return o, nil
}

func (a *App) AuthRepliesQuery(userID int64, opinionID int64) ([]types.Opinion, error) {
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
        LEFT JOIN Opinion_Vote ov on ov.opinion_id=o.opinion_id and voter_id=$1
        LEFT JOIN Topic_Category tc on tc.topic_id = t.topic_id
        LEFT JOIN Category c on c.category_id=tc.category_id
    WHERE o.parent_id=$2
	GROUP BY o.opinion_id, u.user_id, t.topic_id, views, ov.vote
	ORDER BY o.created_on DESC
    `

	rows, err := a.DB.Query(sqlQuery, userID, opinionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return os, nil
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

func (a *App) RepliesQuery(opinionID int64) ([]types.Opinion, error) {
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
        0 as is_following_topic,
        o.is_anonymous,
        0 as is_following,
        o.thumbnails,
        case when o.is_anonymous then '_blank' else o.hls end as hls,
        case when o.is_anonymous then '_blank' else o.dash end as dash,
        case when o.is_anonymous then o.aac else '_blank' end as aac,
        'none' as vote,
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
    WHERE o.parent_id=$1
	GROUP BY o.opinion_id, u.user_id, t.topic_id, views
	ORDER BY o.created_on DESC
	`

	rows, err := a.DB.Query(sqlQuery, opinionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return os, nil
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
