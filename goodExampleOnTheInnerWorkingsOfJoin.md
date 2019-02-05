koubru_prod=> SELECT
t.topic_id,
t.title,
t.details,
t.created_on,
u.username,
u.full_name,
json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null)
FROM topic t inner join kuser u on t.created_by=u.user_id left join topic_category tc on t.topic_id = tc.topic_id and tc.topic_id=1969001063495238692
left join category c on c.category_id=tc.category_id
group by t.topic_id, u.user_id;
topic_id | title | details | created_on | username | full_name | json_agg
---------------------+---------------------------------------------------+---------+-------------------------------+---------------------+----------------------+-------------------------------------------------------------------------------------------------------
1969001063495238692 | First example topic | | 2019-01-31 14:00:39.62561+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443863, "name" : "Science"}, {"id" : 1968424375409443862, "name" : "Politics"}]
1969001138715886629 | Second example topic | | 2019-01-31 14:00:48.591991+00 | rahulsoibam | Rahul Soibam |
1969001195338990630 | Third example topic | | 2019-01-31 14:00:55.342028+00 | rahulsoibam | Rahul Soibam |
1969001300381140007 | Fourth example topic | | 2019-01-31 14:01:07.864258+00 | rahulsoibam | Rahul Soibam |
1969256566880207917 | This is created from API | | 2019-01-31 22:28:18.000012+00 | rahulsoibam | Rahul Soibam |
1969257777029186606 | This is created from API | | 2019-01-31 22:30:42.261383+00 | rahulsoibam | Rahul Soibam |
1969286968462279728 | This is created from API with Categories | | 2019-01-31 23:28:42.151447+00 | rahulsoibam | Rahul Soibam |
1969288191764595761 | This is created from API with Categories | | 2019-01-31 23:31:07.980426+00 | rahulsoibam | Rahul Soibam |
1969293378566751282 | This is from postman | | 2019-01-31 23:41:26.296922+00 | rahulsoibam | Rahul Soibam |
1969295909837603891 | This is a topic created with categories from Paw | | 2019-01-31 23:46:28.048499+00 | rahulsoibam | Rahul Soibam |
1969457782558032948 | This is a topic created with categories from Neng | | 2019-02-01 05:08:04.778396+00 | rahulsoibam | Rahul Soibam |
1969721059473097781 | This is a topic created with categories from Neng | | 2019-02-01 13:51:09.832169+00 | rahulsoibam | Rahul Soibam |
1969843750498731062 | Testing from Android device | | 2019-02-01 17:54:55.744406+00 | nengkhoibachungkham | nengkhoiba chungkham |
1969924293324178487 | Testing topic post android | | 2019-02-01 20:34:57.195639+00 | nengkhoibachungkham | nengkhoiba chungkham |
1970258447886713912 | what is your view towards koubru | | 2019-02-02 07:38:51.52363+00 | rahulsoibam | Rahul Soibam |
1970534525792420921 | Testing new Android post | | 2019-02-02 16:47:22.576135+00 | nengkhoibachungkham | nengkhoiba chungkham |
(16 rows)

---

koubru_prod=> SELECT
t.topic_id,
t.title,
t.details,
t.created_on,
u.username,
u.full_name,
json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null)
FROM topic t inner join kuser u on t.created_by=u.user_id left join topic_category tc on tc.topic_id = 1969001063495238692
left join category c on c.category_id=tc.category_id
group by t.topic_id, u.user_id;
topic_id | title | details | created_on | username | full_name | json_agg
---------------------+---------------------------------------------------+---------+-------------------------------+---------------------+----------------------+-------------------------------------------------------------------------------------------------------
1969001063495238692 | First example topic | | 2019-01-31 14:00:39.62561+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969001138715886629 | Second example topic | | 2019-01-31 14:00:48.591991+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969001195338990630 | Third example topic | | 2019-01-31 14:00:55.342028+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969001300381140007 | Fourth example topic | | 2019-01-31 14:01:07.864258+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969256566880207917 | This is created from API | | 2019-01-31 22:28:18.000012+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969257777029186606 | This is created from API | | 2019-01-31 22:30:42.261383+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969286968462279728 | This is created from API with Categories | | 2019-01-31 23:28:42.151447+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969288191764595761 | This is created from API with Categories | | 2019-01-31 23:31:07.980426+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969293378566751282 | This is from postman | | 2019-01-31 23:41:26.296922+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969295909837603891 | This is a topic created with categories from Paw | | 2019-01-31 23:46:28.048499+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969457782558032948 | This is a topic created with categories from Neng | | 2019-02-01 05:08:04.778396+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969721059473097781 | This is a topic created with categories from Neng | | 2019-02-01 13:51:09.832169+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969843750498731062 | Testing from Android device | | 2019-02-01 17:54:55.744406+00 | nengkhoibachungkham | nengkhoiba chungkham | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969924293324178487 | Testing topic post android | | 2019-02-01 20:34:57.195639+00 | nengkhoibachungkham | nengkhoiba chungkham | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1970258447886713912 | what is your view towards koubru | | 2019-02-02 07:38:51.52363+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
1970534525792420921 | Testing new Android post | | 2019-02-02 16:47:22.576135+00 | nengkhoibachungkham | nengkhoiba chungkham | [{"id" : 1968424375409443862, "name" : "Politics"}, {"id" : 1968424375409443863, "name" : "Science"}]
(16 rows)

---

koubru_prod=> SELECT
t.topic_id,
t.title,
t.details,
t.created_on,
u.username,
u.full_name,
json_agg(json_build_object('id',c.category_id,'name',c.name)) filter (where c.category_id is not null or c.name is not null)
FROM topic t inner join kuser u on t.created_by=u.user_id left join topic_category tc on t.topic_id = tc.topic_id
left join category c on c.category_id=tc.category_id
group by t.topic_id, u.user_id;
topic_id | title | details | created_on | username | full_name | json_agg
---------------------+---------------------------------------------------+---------+-------------------------------+---------------------+----------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------
1969001063495238692 | First example topic | | 2019-01-31 14:00:39.62561+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443863, "name" : "Science"}, {"id" : 1968424375409443862, "name" : "Politics"}]
1969001138715886629 | Second example topic | | 2019-01-31 14:00:48.591991+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375417832473, "name" : "India"}]
1969001195338990630 | Third example topic | | 2019-01-31 14:00:55.342028+00 | rahulsoibam | Rahul Soibam |
1969001300381140007 | Fourth example topic | | 2019-01-31 14:01:07.864258+00 | rahulsoibam | Rahul Soibam |
1969256566880207917 | This is created from API | | 2019-01-31 22:28:18.000012+00 | rahulsoibam | Rahul Soibam |
1969257777029186606 | This is created from API | | 2019-01-31 22:30:42.261383+00 | rahulsoibam | Rahul Soibam |
1969286968462279728 | This is created from API with Categories | | 2019-01-31 23:28:42.151447+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443863, "name" : "Science"}]
1969288191764595761 | This is created from API with Categories | | 2019-01-31 23:31:07.980426+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443863, "name" : "Science"}]
1969293378566751282 | This is from postman | | 2019-01-31 23:41:26.296922+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375417832473, "name" : "India"}, {"id" : 1968424375417832472, "name" : "Famine"}]
1969295909837603891 | This is a topic created with categories from Paw | | 2019-01-31 23:46:28.048499+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968424375409443863, "name" : "Science"}, {"id" : 1968469911558685725, "name" : "google"}, {"id" : 1968471428244177951, "name" : "internet"}]
1969457782558032948 | This is a topic created with categories from Neng | | 2019-02-01 05:08:04.778396+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968469911558685725, "name" : "google"}, {"id" : 1968471428244177951, "name" : "internet"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969721059473097781 | This is a topic created with categories from Neng | | 2019-02-01 13:51:09.832169+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968469911558685725, "name" : "google"}, {"id" : 1968471428244177951, "name" : "internet"}, {"id" : 1968424375409443863, "name" : "Science"}]
1969843750498731062 | Testing from Android device | | 2019-02-01 17:54:55.744406+00 | nengkhoibachungkham | nengkhoiba chungkham | [{"id" : 1968424375417832473, "name" : "India"}, {"id" : 1968424375409443863, "name" : "Science"}, {"id" : 1968469911558685725, "name" : "google"}]
1969924293324178487 | Testing topic post android | | 2019-02-01 20:34:57.195639+00 | nengkhoibachungkham | nengkhoiba chungkham | [{"id" : 1968424375409443863, "name" : "Science"}, {"id" : 1968424375417832472, "name" : "Famine"}, {"id" : 1968424375409443862, "name" : "Politics"}]
1970258447886713912 | what is your view towards koubru | | 2019-02-02 07:38:51.52363+00 | rahulsoibam | Rahul Soibam | [{"id" : 1968469911558685725, "name" : "google"}, {"id" : 1968424375417832473, "name" : "India"}, {"id" : 1968424375409443863, "name" : "Science"}]
1970534525792420921 | Testing new Android post | | 2019-02-02 16:47:22.576135+00 | nengkhoibachungkham | nengkhoiba chungkham | [{"id" : 1968424375409443863, "name" : "Science"}, {"id" : 1968424375417832472, "name" : "Famine"}, {"id" : 1968424375409443862, "name" : "Politics"}]
(16 rows)

---
