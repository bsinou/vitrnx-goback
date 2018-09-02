// Package dummydata provides convenience data to test the backend
package dummydata

// TestPost1 is a simple post.
var TestPost1 = `{ 
	"isPublic" : true,
	"path" : "simple-test",
	"title" : "a title",
	"tags" : "news tags",
	"thumb" : "test.jpg",
	"hero" : "test.jpg",
	"desc" : "A short description",
	"body" : "A very long body"
	}`

// TestPost2 should be inserted correcty.
var TestPost2 = `{ 
	"isPublic" : true,
	"path" : "simple-test-2",
	"title" : "An other title",
	"tags" : "news tags second",
	"thumb" : "test2.jpg",
	"hero" : "test2.jpg",
	"desc" : "A short description",
	"body" : "A very long body"
}`

// TestPost3 should not be inserted path is sams as post1's.
var TestPost3 = `{ 
	"isPublic" : false,
	"path" : "simple-test",
	"title" : "Duplicate path",
	"tags" : "news tags",
	"thumb" : "test.jpg",
	"hero" : "test.jpg",
	"desc" : "A short description",
	"body" : "This post should not be imported"
}`
