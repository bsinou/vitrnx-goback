package model

// This file centralises constanst that should be instance specific
// but that are still hard coded in a first approach to shorten implementation of v0

const ApiPrefix = "/api/"

// Known Roles
const RoleAdmin = "ADMIN"
const RoleUserAdmin = "USER_ADMIN"
const RoleModerator = "MODERATOR"
const RoleEditor = "EDITOR"
const RoleVolunteer = "VOLUNTEER"
const RoleGuest = "GUEST"
const RoleRegistered = "REGISTERED"
const RoleAnonymous = "ANONYMOUS"

// Known Policies
const PolicyCanRead = "CAN_READ"
const PolicyCanEdit = "CAN_EDIT"
const PolicyCanManage = "CAN_MANAGE"
