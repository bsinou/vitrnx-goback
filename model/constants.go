package model

/* CONTEXT KEYS */

// various
const KeyDb = "db"
const KeyUserRecord = "userRecord"
const KeyMgoId = "id"

// Users
const KeyUser = "user"
const KeyUserId = "userId"
const KeyUserName = "userName"
const KeyEmailVerified = "emailVerified"
const KeyAuth = "authorization"
const KeyClaims = "claims"

// Posts
const KeyPath = "path"
const KeyTags = "tags"

/* FIREBASE specific keys */

const FbKeyUserId = "user_id"
const FbKeyEmail = "email"
const FbKeyEmailVerified = "email_verified"

/* LOCAL CONSTANTS */

// Known Policies
const PolicyCanRead = "CAN_READ"
const PolicyCanEdit = "CAN_EDIT"
const PolicyCanManage = "CAN_MANAGE"

// store type
const StoreTypeMgo = "mongodb"
const StoreTypeGorm = "gorm"
