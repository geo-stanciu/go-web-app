package models

import "time"

// UserModel - UserModel
type UserModel struct {
	UserID          int       `sql:"user_id"`
	Username        string    `sql:"username"`
	Name            string    `sql:"name"`
	Surname         string    `sql:"surname"`
	Email           string    `sql:"email"`
	PasswordExpires bool      `sql:"password_expires" json:"password_expires"`
	CreationTime    time.Time `sql:"creation_time" json:"creation_time"`
	LastUpdate      time.Time `sql:"last_update" json:"last_update"`
	Activated       bool      `sql:"activated" json:"activated"`
	LockedOut       bool      `sql:"locked_out" json:"locked_out"`
	Valid           bool      `sql:"valid" json:"valid"`
	Password        string    `json:"-"`
}
