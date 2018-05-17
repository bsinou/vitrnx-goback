package model

/* CONTEXT KEYS */

// various
// const KeyDb = "db"
const KeyUserDb = "userDb"
const KeyDataDb = "dataDb"
const KeyUserRecord = "userRecord"
const KeyMgoID = "id"

// Users
const KeyUser = "user"
const KeyUserID = "userId"
const KeyUserDisplayName = "userDisplayName"
const KeyEmail = "email"
const KeyEmailVerified = "emailVerified"
const KeyAuth = "authorization"
const KeyClaims = "claims"
const KeyUserRoles = "userRoles"

const KeyEditedUser = "user"
const KeyEditedUserRoles = "editedUserRoles"

// User Meta
const KeyPresence = "presence"

// Posts
const KeyPost = "post"
const KeyPath = "path"
const KeyTags = "tags"

// Comments
const KeyComment = "comment"
const KeyParentID = "parentId"

// Tasks
const KeyTask = "task"
const KeyCategory = "category"
const KeyCategoryID = "categoryId"

/* FIREBASE specific keys */

const FbKeyUserID = "user_id"
const FbKeyEmail = "email"
const FbKeyEmailVerified = "email_verified"

/* LOCAL CONSTANTS */

// Store types
const StoreTypeMgo = "mongodb"
const StoreTypeGorm = "gorm"
