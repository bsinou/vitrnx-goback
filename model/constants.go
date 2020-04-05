package model

/* CONTEXT KEYS */

// various
const (
	KeyDb         = "db"
	StoreType     = "gorm"
	KeyID         = "id"
	KeyPostRecord = "postRecord"
	KeyUserRecord = "userRecord"
)

// Posts
const (
	KeyPost = "post"
	KeySlug = "slug"
	KeyTags = "tags"
	KeyTag  = "tag"
)

/* CONTEXT KEYS */

// Users
const (
	KeyIdToken = "id_token"

	KeyAuthorization   = "Authorization"
	KeyUser            = "user"
	KeyUserID          = "userId"
	KeyUserDisplayName = "userDisplayName"
	KeyEmail           = "email"
	KeyEmailVerified   = "emailVerified"
	KeyClaims          = "claims"
	KeyUserRoles       = "userRoles"
	KeyEditedUser      = "user"
	KeyEditedUserRoles = "editedUserRoles"
)
