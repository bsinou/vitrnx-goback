// Package dummydata provides convenience data to test the backend
package dummydata

// TestUser1 is a simple user.
var TestUser1 = `{ 
	"userId":"test",
		"name":"Jane",
		"email":"jane@example.com",
		"address":"a perfect address",
		"roles":[
			{
				"roleId":"EDITOR",
				"label":"Editor"
				},
			{
				"roleId":"ADMIN",
				"label":"Admin"
				}
			]
		}
	`

var TestUser2 = `{
		"userId":"test2",
		"name":"John",
		"email":"john@example.com",
		"address":"a perfect address",
		"roles":[
			{
				"roleId":"EDITOR",
				"label":"Editor"
				},
			{
				"roleId":"ADMIN",
				"label":"Admin"
				}
			]
		}
	`
var fullJSONUserExample = `{
		"ID":0,
		"CreatedAt":"0001-01-01T00:00:00Z",
		"UpdatedAt":"0001-01-01T00:00:00Z",
		"DeletedAt":null,
		"userId":"test",
		"name":"John",
		"email":"john@example.com",
		"address":"",
		"roles":[
			{
				"ID":0,
				"CreatedAt":"0001-01-01T00:00:00Z",
				"UpdatedAt":"0001-01-01T00:00:00Z",
				"DeletedAt":null,
				"roleId":"EDITOR",
				"label":"Editor"
				},
			{
				"ID":0,
				"CreatedAt":"0001-01-01T00:00:00Z",
				"UpdatedAt":"0001-01-01T00:00:00Z",
				"DeletedAt":null,
				"roleId":"ADMIN",
				"label":"Admin"
				}
			]
		}
	`
