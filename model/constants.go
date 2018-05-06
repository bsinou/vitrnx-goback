package model

/* CONTEXT KEYS */

// various
const KeyDb = "db"
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
const KeyRoles = "roles"

// Posts
const KeyPost = "post"
const KeyPath = "path"
const KeyTags = "tags"

// Comments
const KeyParentID = "parentId"

/* FIREBASE specific keys */

const FbKeyUserID = "user_id"
const FbKeyEmail = "email"
const FbKeyEmailVerified = "email_verified"

/* LOCAL CONSTANTS */

// Store types
const StoreTypeMgo = "mongodb"
const StoreTypeGorm = "gorm"
