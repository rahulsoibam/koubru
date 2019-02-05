```sql
--- NON_AUTH CATEGORIES LIST BY FOLLOWER COUNT
select
   c.category_id,
   c.name,
   count(cf.user_id) as count_followers
from category c full join category_follower cf on c.category_id =cf.category_id
where name like '%'
group by c.category_id
order by count_followers desc;

--- Result
     category_id     |   name   | count_followers
---------------------+----------+-----------------
 1968424375417832472 | Famine   |               2
 1968424375409443863 | Science  |               1
 1968424375409443862 | Politics |               1
 1968469569395754011 | Facebook |               1
 1968424375417832473 | India    |               1
 1968469911558685725 | google   |               0
 1968471428244177951 | internet |               0
(7 rows)

--- Without followers count
select
    c.category_id,
    c.name
from category c full join category_follower cf on c.category_id =cf.category_id
where name like '%'
group by c.category_id
order by (select count(cf.user_id)) desc;

--- Result
 category_id     |   name
---------------------+----------
 1968424375417832472 | Famine
 1968424375409443863 | Science
 1968424375409443862 | Politics
 1968469569395754011 | Facebook
 1968424375417832473 | India
 1968469911558685725 | google
 1968471428244177951 | internet
(7 rows)

--- Details of one topic
select c.category_id, c.name, c.created_on, u.username, u.full_name, u.photo_url,
case when (cf.user_id=1967600534613394434) then 1 else 0 end as is_following,
count(cf.user_id) as follower_count,
count(tc.topic_id) as topic_count

from category c inner join kuser u on c.created_by = u.user_id and c.category_id = 1968424375417832472
right join category_follower cf on c.category_id = cf.category_id and cf.category_id = 1968424375417832472
right join topic_category tc on c.category_id = tc.category_id and tc.category_id =1968424375417832472
group by c.category_id, u.user_id, cf.user_id;

-[ RECORD 1 ]--+---------------------------------------------------------------------------------------------------
category_id    |
name           |
created_on     |
username       |
full_name      |
photo_url      |
is_following   | 0
follower_count | 0
topic_count    | 23
-[ RECORD 2 ]--+---------------------------------------------------------------------------------------------------
category_id    | 1968424375417832472
name           | Famine
created_on     | 2019-01-30 18:54:53.047921+00
username       | rahulsoibam
full_name      | Rahul Soibam
photo_url      | https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg
is_following   | 1
follower_count | 2
topic_count    | 2
-[ RECORD 3 ]--+---------------------------------------------------------------------------------------------------
category_id    | 1968424375417832472
name           | Famine
created_on     | 2019-01-30 18:54:53.047921+00
username       | rahulsoibam
full_name      | Rahul Soibam
photo_url      | https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg
is_following   | 0
follower_count | 2
topic_count    | 2

select c.category_id, c.name, c.created_on, u.username, u.full_name,
case when (cf.user_id=1967600534613394434) then 1 else 0 end as is_following



from category c inner join kuser u on c.created_by = u.user_id and c.category_id = 1968424375417832472
right join category_follower cf on c.category_id = cf.category_id and cf.category_id = 1968424375417832472
right join topic_category tc on c.category_id = tc.category_id and tc.category_id =1968424375417832472
where c.category_id=1968424375417832472;

     category_id     |  name  |          created_on           |  username   |  full_name   | is_following
---------------------+--------+-------------------------------+-------------+--------------+--------------
 1968424375417832472 | Famine | 2019-01-30 18:54:53.047921+00 | rahulsoibam | Rahul Soibam |            1
 1968424375417832472 | Famine | 2019-01-30 18:54:53.047921+00 | rahulsoibam | Rahul Soibam |            1
 1968424375417832472 | Famine | 2019-01-30 18:54:53.047921+00 | rahulsoibam | Rahul Soibam |            0
 1968424375417832472 | Famine | 2019-01-30 18:54:53.047921+00 | rahulsoibam | Rahul Soibam |            0
 
 
 
 
 
select c.category_id, u.user_id, cf.*, tc.*,
case when (cf.user_id=1967600534613394434) then 1 else 0 end as is_following
from category c inner join kuser u on c.created_by = u.user_id and c.category_id = 1968424375417832472
right join category_follower cf on c.category_id = cf.category_id and cf.category_id = 1968424375417832472
right join topic_category tc on c.category_id = tc.category_id and tc.category_id =1968424375417832472
where c.category_id=1968424375417832472;
-[ RECORD 1 ]+------------------------------
category_id  | 1968424375417832472
user_id      | 1967600534613394434
category_id  | 1968424375417832472
user_id      | 1967600534613394434
followed_on  | 2019-02-02 00:12:38.371699+00
topic_id     | 1969293378566751282
category_id  | 1968424375417832472
created_on   | 2019-01-31 23:41:26.296922+00
is_following | 1
-[ RECORD 2 ]+------------------------------
category_id  | 1968424375417832472
user_id      | 1967600534613394434
category_id  | 1968424375417832472
user_id      | 1967600534613394434
followed_on  | 2019-02-02 00:12:38.371699+00
topic_id     | 1969924293324178487
category_id  | 1968424375417832472
created_on   | 2019-02-01 20:34:57.195639+00
is_following | 1
-[ RECORD 3 ]+------------------------------
category_id  | 1968424375417832472
user_id      | 1967600534613394434
category_id  | 1968424375417832472
user_id      | 1967646297590596612
followed_on  | 2019-02-02 00:13:43.691412+00
topic_id     | 1969293378566751282
category_id  | 1968424375417832472
created_on   | 2019-01-31 23:41:26.296922+00
is_following | 0
-[ RECORD 4 ]+------------------------------
category_id  | 1968424375417832472
user_id      | 1967600534613394434
category_id  | 1968424375417832472
user_id      | 1967646297590596612
followed_on  | 2019-02-02 00:13:43.691412+00
topic_id     | 1969924293324178487
category_id  | 1968424375417832472
created_on   | 2019-02-01 20:34:57.195639+00
is_following | 0


koubru_prod=> select c.category_id, u.user_id,
--case when (cf.user_id=1967600534613394434) then 1 else 0 end as is_following,
count(distinct cf.*) as follower_count,
count(distinct tc.*) as topic_count



from category c inner join kuser u on c.created_by = u.user_id and c.category_id = 1968424375417832472
right join category_follower cf on c.category_id = cf.category_id and cf.category_id = 1968424375417832472
right join topic_category tc on c.category_id = tc.category_id and tc.category_id =1968424375417832472
where c.category_id=1968424375417832472
group by c.category_id, u.user_id;
     category_id     |       user_id       | follower_count | topic_count
---------------------+---------------------+----------------+-------------
 1968424375417832472 | 1967600534613394434 |              2 |           2
(1 row)




select c.category_id, u.user_id,
case when (cf.user_id=1967600534613394434) then 1 else 0 end as is_following,
count(distinct cf.*) as follower_count,
count(distinct tc.*) as topic_count



from category c inner join kuser u on c.created_by = u.user_id and c.category_id = 1968424375417832472
right join category_follower cf on c.category_id = cf.category_id and cf.category_id = 1968424375417832472
right join topic_category tc on c.category_id = tc.category_id and tc.category_id =1968424375417832472
where c.category_id=1968424375417832472
group by c.category_id, u.user_id;
```

