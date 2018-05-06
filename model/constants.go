package model

/* CONTEXT KEYS */

// various
const KeyDb = "db"
const KeyUserRecord = "userRecord"
const KeyMgoID = "id"

// Users
const KeyUser = "user"
const KeyUserID = "userId"
const KeyUserName = "userName"
const KeyEmailVerified = "emailVerified"
const KeyAuth = "authorization"
const KeyClaims = "claims"

// Posts
const KeyPath = "path"
const KeyTags = "tags"

// Comments
const KeyParentID = "parentId"

/* FIREBASE specific keys */

const FbKeyUserID = "user_id"
const FbKeyEmail = "email"
const FbKeyEmailVerified = "email_verified"

/* LOCAL CONSTANTS */

const ApiPrefix = "/api/"

// Known Policies
const PolicyCanRead = "CAN_READ"
const PolicyCanEdit = "CAN_EDIT"
const PolicyCanManage = "CAN_MANAGE"

// Store types
const StoreTypeMgo = "mongodb"
const StoreTypeGorm = "gorm"
