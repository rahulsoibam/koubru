package errs

import "errors"

var (
	InternalServerError            = errors.New("Opps, something went wrong. Please try again.")
	TopicNotFound                  = errors.New("Topic not found.")
	CategoryNotFound               = errors.New("Category not found.")
	CategoryFollowAlreadyFollowing = errors.New("You have already followed this topic.")
	TopicFollowAlreadyFollowing    = errors.New("You have already followed this category.")
	CategoryUnfollowNotFollowing   = errors.New("You do not follow this category.")
	TopicUnfollowNotFollowing      = errors.New("You do not follow this topic.")
	BadRequest                     = errors.New("Bad request.")
	Unauthorized                   = errors.New("You are not authorized to perform this action.")
	InvalidHash                    = errors.New("The encoded hash is not in the correct format.")
	Argon2IncompatibleVersion      = errors.New("Incompatible version of argon2.")
	NoHeader                       = errors.New("Authorization header required.")
	NoBearer                       = errors.New("Authorization requires Bearer scheme.")
	UnintendedExecution            = errors.New("This should not be executing.")
	UserNotFound                   = errors.New("User not found.")
	NoPasswordSet                  = errors.New("A password is not set for this account. Login using social account and set a password for future login.")
	PasswordNotMatch               = errors.New("Password does not match.")
	UserFollowAlreadyFollowing     = errors.New("You are already following this user.")
	UserUnfollowNotFollowing       = errors.New("You do not follow this user.")
)
