package entities

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID             string `json:"id"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Email          string     `gorm:"primary_key" json:"email"`
	StreetAddress  string     `json:"street_address"`
	FlowAddress    string     `json:"flow_address"`
	ExternalWallet bool       `json:"external_wallet"`
	Status         bool       `json:"status"`
}

func (u *User) Disable() {
	u.Status = false
}

func (u *User) Enable() {
	u.Status = true
}
