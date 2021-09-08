package models

// The models package will store the database schema.
// We will use struct type to represent or map the database schema in golang.
// The User struct is a representation of users table which we created in ElephantSQL.

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int64  `json:"age"`
}
