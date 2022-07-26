package model

// User represents the credentials of an actor trying to log in.
//
// swagger:model
type User struct {
	// the username
	// required: true
	// example: admin
	Username string `json:"username" bson:"username"`

	// the password
	// required: true
	// example: password
	Password string `json:"password" bson:"password"`
}
