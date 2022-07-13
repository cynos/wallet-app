package payment

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	TransactionID string `gorm:"varchar(100);index"`
	ProductID     string `gorm:"varchar(100);index"`
	UsersID       uint
	Amount        int `gorm:"default:0"`
	Confirm       int
	Remark        string
}
