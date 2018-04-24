package model

/* CONTEXT KEYS */

// various
const KeyDb = "db"
const KeyUserRecord = "userRecord"

// Users
const KeyUserId = "userId"
const KeyUserName = "userName"
const KeyEmailVerified = "emailVerified"
const KeyAuth = "authorization"
const KeyClaims = "claims"

// Posts
const KeyPath = "path"

/* FIREBASE specific keys */

const FbKeyUserId = "user_id"
const FbKeyEmail = "email"
const FbKeyEmailVerified = "email_verified"

/* LOCAL CONSTANTS */

// store type
const StoreTypeMgo = "mongodb"
const StoreTypeGorm = "gorm"
