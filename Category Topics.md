koubru_prod=> select t.topic_id, t.title, array_agg(c.category_id), array_agg(c.name)
from topic t inner join topic_category tc on t.topic_id = tc.topic_id and tc.category_id = 1968424375409443863
inner join topic_category tc2 on tc2.topic_id=t.topic_id
inner join category c on c.category_id=tc2.category_id
group by t.topic_id;
      topic_id       |                       title                       |                           array_agg                           |         array_agg
---------------------+---------------------------------------------------+---------------------------------------------------------------+---------------------------
 1969001063495238692 | First example topic                               | {1968424375409443862,1968424375409443863}                     | {Politics,Science}
 1969286968462279728 | This is created from API with Categories          | {1968424375409443863}                                         | {Science}
 1969288191764595761 | This is created from API with Categories          | {1968424375409443863}                                         | {Science}
 1969295909837603891 | This is a topic created with categories from Paw  | {1968424375409443863,1968469911558685725,1968471428244177951} | {Science,google,internet}
 1969457782558032948 | This is a topic created with categories from Neng | {1968424375409443863,1968469911558685725,1968471428244177951} | {Science,google,internet}
 1969721059473097781 | This is a topic created with categories from Neng | {1968424375409443863,1968469911558685725,1968471428244177951} | {Science,google,internet}
 1969843750498731062 | Testing from Android device                       | {1968424375409443863,1968424375417832473,1968469911558685725} | {Science,India,google}
 1969924293324178487 | Testing topic post android                        | {1968424375409443862,1968424375409443863,1968424375417832472} | {Politics,Science,Famine}
 1970258447886713912 | what is your view towards koubru                  | {1968424375409443863,1968424375417832473,1968469911558685725} | {Science,India,google}
 1970534525792420921 | Testing new Android post                          | {1968424375409443862,1968424375409443863,1968424375417832472} | {Politics,Science,Famine}
(10 rows)

koubru_prod=>