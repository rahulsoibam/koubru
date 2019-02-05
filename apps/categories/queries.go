package categories

import (
	"database/sql"

	"github.com/rahulsoibam/koubru/types"
)

// DONE
func (a *App) ListQuery(q string, limit int, offset int) ([]types.Category_, error) {
	var err error
	cs := []types.Category_{}
	q = q + "%"

	// Query to list all categories by follower count
	sqlQuery := `
	SELECT	
	    c.category_id,
		c.name,
		0
	FROM category c FULL JOIN category_follower cf ON c.category_id=cf.category_id
	WHERE name LIKE $1
	GROUP BY c.category_id
	ORDER BY (select count(cf.follower_id)) DESC;
	LIMIT $3 OFFSET $4
	`
	rows, err := a.DB.Query(sqlQuery, q, limit, offset)
	if err == sql.ErrNoRows {
		return cs, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var c types.Category_
		err := rows.Scan(&c.ID, &c.Name, &c.IsFollowing)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// DONE
func (a *App) AuthListQuery(userID int64, q string, limit int, offset int) ([]types.Category_, error) {
	cs := []types.Category_{}

	// Query to list all categories and put following ones at the top
	sqlQuery := `
	SELECT
		c.category_id, 
		c.name, 
		CASE WHEN (cf.follower_id IS NULL) THEN  0 ELSE 1 END AS is_following
	FROM category c LEFT JOIN category_follower cf ON c.category_id=cf.category_id AND cf.follower_id=$1
	WHERE c.name LIKE $2
	ORDER BY cf.followed_on DESC NULLS LAST
	LIMIT $3 OFFSET $4
	`

	var err error
	query := q + "%"

	rows, err := a.DB.Query(sqlQuery, userID, query, limit, offset)
	if err == sql.ErrNoRows {
		return cs, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var c types.Category_
		err := rows.Scan(&c.ID, &c.Name, &c.IsFollowing)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// DONE
func (a *App) AuthCreateQuery(userID int64, c types.NewCategory) (types.Category, error) {
	cres := types.Category{}
	tx, err := a.DB.Begin()
	if err != nil {
		return cres, err
	}
	var categoryID int64
	err = tx.QueryRow("INSERT INTO category (name, created_by) VALUES ($1, $2) RETURNING category_id", c.Name, userID).Scan(&categoryID)
	if err != nil {
		tx.Rollback()
		return cres, err
	}

	_, err = tx.Exec("INSERT INTO category_follower (category_id, follower_id) VALUES ($1, $2)", categoryID, userID)
	if err != nil {
		tx.Rollback()
		return cres, err
	}

	// TODO GET CATEGORY PAGE
	cres, err = a.AuthGetQuery(userID, categoryID)
	if err != nil {
		tx.Rollback()
		return cres, err
	}

	err = tx.Commit()
	if err != nil {
		return cres, nil
	}
	return cres, nil
}

// DONE
func (a *App) AuthGetQuery(userID int64, categoryID int64) (types.Category, error) {
	c := types.Category{}
	sqlQuery := `
	SELECT
        c.category_id,
        c.name,
        c.created_on,
        u.username,
        u.picture,
		CASE WHEN EXISTS (SELECT 1 from category_follower cf WHERE cf.category_id=c.category_id AND cf.follower_id=$1) THEN 1 ELSE 0 END AS is_following,
		(select count(*) from topic_category where category_id=c.category_id) as topics_count,
		(select count(*) from category_follower where category_id=c.category_id) as followers_count
    FROM
        Category c INNER JOIN KUser u ON c.creator_id=u.user_id
    WHERE
        c.category_id=$2
	`

	err := a.DB.QueryRow(sqlQuery, userID, categoryID).Scan(&c.ID, &c.Name, &c.CreatedOn, &c.CreatedBy.Username, &c.CreatedBy.FullName, &c.CreatedBy.Picture, &c.IsFollowing, &c.Counts.Topics, &c.Counts.Followers)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, nil
		}
		return c, err
	}

	return c, nil
}

// DONE
func (a *App) GetQuery(categoryID int64) (types.Category, error) {
	c := types.Category{}
	sqlQuery := `
	SELECT
        c.category_id,
        c.name,
        c.created_on,
        u.username,
		u.picture,
		0 as is_following,
		(select count(*) from topic_category where category_id=c.category_id) as topics_count,
		(select count(*) from category_follower where category_id=c.category_id) as followers_count
    FROM
        Category c INNER JOIN KUser u ON c.creator_id=u.user_id
    WHERE
        c.category_id=$1
	`

	err := a.DB.QueryRow(sqlQuery, categoryID).Scan(&c.ID, &c.Name, &c.CreatedOn, &c.CreatedBy.Username, &c.CreatedBy.FullName, &c.CreatedBy.Picture, &c.IsFollowing)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, nil
		}
		return c, err
	}

	return c, nil
}

// TODO ADD PAGINATION
func (a *App) AuthFollowersQuery(userID int64, categoryID int64) ([]types.User_, error) {
	fs := []types.User_{}
	sqlQuery := `
	SELECT
        u.username,
        u.full_name,
        u.picture,
        CASE WHEN cf.follower_id=$1 THEN 1 ELSE 0 END AS is_self,
        CASE WHEN EXISTS (SELECT 1 FROM user_follower uf where uf.user_id=u.user_id AND uf.follower_id=$1) THEN 1 ELSE 0 END AS is_following
    FROM
        KUser u INNER JOIN Category_Follower cf ON u.user_id = cf.follower_id
    WHERE cf.category_id=$2;
    ORDER BY is_self desc, is_following desc
	`

	rows, err := a.DB.Query(sqlQuery, userID, categoryID)
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
			return fs, err
		}
		fs = append(fs, f)
	}
	err = rows.Err()
	if err != nil {
		return fs, err
	}

	// No errors found, return nil as errors
	return fs, nil

}

// TODO ADD PAGINATION
func (a *App) FollowersQuery(categoryID int64) ([]types.User_, error) {
	fs := []types.User_{}
	// List followers of a category ordered by their ids
	sqlQuery := `
	SELECT
        u.username,
        u.full_name,
        u.picture
    FROM
        KUser u INNER JOIN Category_Follower cf on u.user_id = cf.follower_id left join user_follower uf on u.user_id=uf.user_id
    WHERE cf.category_id=$1
    GROUP BY u.user_id
    ORDER BY (SELECT count(uf.follower_id)) DESC
	`

	rows, err := a.DB.Query(sqlQuery, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fs, nil
		}
		return fs, err
	}

	defer rows.Close()
	for rows.Next() {
		f := types.User_{}
		err := rows.Scan(&f.Username, &f.FullName, &f.Picture) // Need only scan these three fields, others will be initialized to false by default
		if err != nil {
			return fs, err
		}
		fs = append(fs, f)
	}
	err = rows.Err()
	if err != nil {
		return fs, err
	}
	// No errors found, return nil as error
	return fs, nil
}

// TODO ADD PAGINATION
func (a *App) AuthTopicsQuery(userID int64, categoryID int64) ([]types.Topic_, error) {
	ts := []types.Topic_{}
	// List all topics of a category sorted by follower count then in chronologial order
	sqlQuery := `
	SELECT
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.picture,
		coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) FILTER (WHERE c.category_id IS NOT NULL OR c.name IS NOT NULL), '[]'::json),
		CASE WHEN EXISTS (SELECT 1 FROM topic_follower tofo WHERE tofo.topic_id=t.topic_id AND tofo.follower_id=$1) THEN 1 ELSE 0 END AS is_following
	FROM topic t inner join kuser u on t.creator_id=u.user_id inner join topic_category tc on t.topic_id = tc.topic_id and tc.category_id = $2
	inner join topic_category tc2 on tc2.topic_id=t.topic_id
	inner join category c on c.category_id=tc2.category_id
	group by t.topic_id, u.user_id
	ORDER BY (SELECT COUNT(tf.follower_id) FROM topic_follower tf WHERE tf.topic_id=t.topic_id GROUP BY tf.topic_id) DESC, t.created_on DESC
	`
	rows, err := a.DB.Query(sqlQuery, userID, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ts, nil
		}
		return ts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := types.Topic_{}
		// IsFollowing will default to false
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

// TODO ADD PAGINATION
func (a *App) TopicsQuery(categoryID int64) ([]types.Topic_, error) {
	ts := []types.Topic_{}
	// List all topics of a category sorted by follower count then in chronologial order
	sqlQuery := `
	SELECT
		t.topic_id,
		t.title,
		t.details,
		t.created_on,
		u.username,
		u.full_name,
		u.picture,
		coalesce(json_agg(json_build_object('id',c.category_id,'name',c.name)) FILTER (WHERE c.category_id IS NOT NULL OR c.name IS NOT NULL), '[]'::json)
	FROM topic t inner join kuser u on t.creator_id=u.user_id inner join topic_category tc on t.topic_id = tc.topic_id and tc.category_id = $1
	inner join topic_category tc2 on tc2.topic_id=t.topic_id
	inner join category c on c.category_id=tc2.category_id
	group by t.topic_id, u.user_id
	ORDER BY (SELECT COUNT(tf.follower_id) FROM topic_follower tf WHERE tf.topic_id=t.topic_id GROUP BY tf.topic_id) DESC, t.created_on DESC
	`
	rows, err := a.DB.Query(sqlQuery, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ts, nil
		}
		return ts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := types.Topic_{}
		// IsFollowing will default to false
		err := rows.Scan(&t.ID, &t.Title, &t.Details, &t.CreatedOn, &t.CreatedBy.Username, &t.CreatedBy.FullName, &t.CreatedBy.Picture, (*[]byte)(&t.Categories))
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}

	return ts, nil
}
