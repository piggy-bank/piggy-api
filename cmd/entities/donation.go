package entities

import (
	"github.com/jinzhu/gorm"
)

type Donation struct {
	gorm.Model
	PiggyID                   string `json:"piggy_id"`
	Piggy                     Piggy  `json:"piggy"`
	SenderID                  string `json:"sender_id"`
	Amount                    int64  `json:"amount"`
	BrokePiggy                bool   `json:"broke"`
	PaymentRelatedTransaction string `json:"transaction_id"`
}
