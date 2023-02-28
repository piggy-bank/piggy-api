package entities

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Piggy struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	Goal        int64      `json:"goal"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	UserAddress int        `json:"expiry_month"`
	Donations   []Donation `json:"donation"`
}
