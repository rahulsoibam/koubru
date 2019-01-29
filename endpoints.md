# Koubru REST API

On selection of country on first entry into the app, the server will return a Bearer Token that corresponds to a guest session in the server. This session data will be used when the user eventually creates an account.

So every request here (with the exception of categories and countries related) needs to be sent with the header `Authorization: Bearer *Token*`



## Feed

#### Get feed

```http
GET /feed
```

## Auth

#### Login using username/phone/email and password

```http
POST /auth/login
```

#### Register by providing details

```http
POST /auth/register
```

#### Login using social media accounts

```http
POST /auth/social
```

#### Logout

```http
POST /auth/logout
```



## User

#### Get details of authenticated user

```http
GET /user
```

#### Update details of authenticated user

```http
PATCH /user
```

#### Delete or deactivate authenticated user

```http
DELETE /user
```

#### To list followers of authenticated user

```http
GET /user/followers
```

#### To list users who the authenticated user is following

```http
GET /user/following
```

#### To check if the a user is following the authenticated user - to show **Follows you** when viewing their profile

```http
GET /user/follower/:username
# response if following
Status: 204 No content
# response if not following
Status: 404 Not found
```

#### To check if the authenticated user is following another user - to show **Following** when viewing their profile

```http
GET /user/following/:username
# response if following
Status: 204 No content
# response if not following
Status: 404 Not found
```

#### List opinions of authenticated user, sorted in chronological order

```http
GET /user/opinions
```

*pagination example, chronological order*

```http
GET /user/opinions?page=2&per_page=100 
# default is 15, maybe change with user connection speed on calling side?
```

*ordering by time example (default)*

```http
GET /user/opinions?sort=created&order=desc
```

*ordering by popularity example`(add this or not?)`*

```http
GET /user/opinions?sort=popularity&order=desc
```

*combining*

```http
GET /user/opinions?sort=created&page=2&per_page=50
```

#### List topics by authenticated user sorted in chronological order

```http
GET /user/topics
```

*pagination example*

```http
GET /user/topics?page=2&per_page=50 
# default is 30, maybe change with user connection speed on calling side?
```

*ordering by time example (default)*

```http
GET /user/topics?sort=created&order=desc
```

*ordering by top example`(add this or not?)`*

```http
GET /user/topics?sort=top&order=desc
```

*combining*

```http
GET /user/topics?sort=popularity&page=2&per_page=50
```



## Users

#### To list all users by joining date `(add this or not?)`

```http
GET /users
```

#### To view the profile of a user

```http
GET /users/:username
```

#### To list opinions of user sorted in chronological order

```http
GET /users/:username/opinions
```

*pagination*

```http
GET /users/:username/opinions?page=2&per_page=30
# default is 15
```

*ordering by recommended*

```http
GET /users/:username/opinions?sort=recommended
```

*ordering by top*

```http
GET /users/:username/opinions?sort=top&duration=7
# duration in days
```

#### To list topics of user sorted in chronological order

```http
GET /users/:username/topics
```

*Same pagination and ordering as opinions*

#### To list followers of a user

```http
GET /user/:username/followers
```

#### To list the users a user is following

```http
GET /users/:username/following
```



## Explore

#### List contents of explore page

```http
GET /explore
```

#### Nearby

```http
GET /explore/location=nearby
```

#### By location

```http
GET /explore/location=IND
# 3 letter digit representing country
```



## Search

#### Search users

```http
GET /search/users?q=rahul+soibam&...
```

#### Search topics

```http
GET /search/topics?q=narendra+modi&...
```

#### Search categories

```http
GET /search/categories?q=politics&...
```



## Topics

#### List all topics `(add this or not?)`

```http
GET /topics
```

#### Create a topic

```http
POST /topics
```

#### Get details of a specific topic

```http
GET /topics/:id
```

#### Update details of a specific topic

```http
PATCH /topics/:id
```

#### Remove user from topic and make topic anonymous *(only way to delete a topic is by reporting it.)*

```http
DELETE /topics/:id
```

#### List followers of a topic

```http
GET /topics/:id/followers
```

#### Follow a topic

```http
PUT /topics/:id/follow
```

#### Unfollow a topic

```http
DELETE /topics/:id/follow
```

#### Report a topic

```http
POST /topic/:id/report
```



## Opinions

#### List all topics `(add this or not?)`

```http
GET /opinions
```

#### Create an opinion - including adding to topic and replying to opinion

```http
POST /opinions
```

#### Details of a specific opinion

```http
GET /opinions/:id
```

#### Delete an opinion

```http
DELETE /opinions/:id
```

#### Get followers of an opinion

```http
GET /opinion/:id/followers
```

####  Get replies of an opinion

```http
GET /opinion/:id/replies
```

#### Follow an opinion

```http
PUT /opinions/:id/follow
```

#### Unfollow an opinion

```http
DELETE /opinions/:id/unfollow
```

#### Report an opinion

```http
PUT /opinions/:id/report
```

#### Upvote an opinion

```http
PUT /opinions/:id/upvote
```

#### Un-upvote an opinion

```http
DELETE /opinions/:id/upvote
```

#### Downvote an opinion

```http
PUT /opinions/:id/downvote
```

#### Un-downvote an opinion

```http
DELETE /opinions/:id/downvote
```



## Categories

#### List all categories

```http
GET /categories
```

#### Create a category

```http
POST /categories
```

#### Get details of a specific category

```http
GET /categories/:id
```

#### Follow a topic

```http
PUT /categories/:id/follow
```

#### Unfollow a topic

```http
DELETE /categories/:id/follow
```

#### Bulk follow topics (usually first use of app)

```http
POST /categories/follow
```



## Countries

#### List countries

```http
GET /countries
```

#### Select country

```http
POST /countries
```

