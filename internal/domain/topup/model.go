package topup

import "time"

type Topup struct {
	ID            uint      `gorm:"primaryKey" json:"id,omitempty"`
	Amount        int64     `gorm:"index" json:"amount,omitempty"`
	TransactionID string    `gorm:"index" json:"transaction_id,omitempty"`
	Remark        string    `gorm:"index" json:"remark,omitempty"`
	UsersID       uint      `json:"users_id,omitempty"`
	CreatedAt     time.Time `gorm:"default:current_timestamp;index" json:"created_at,omitempty"`
	UpdatedAt     time.Time `gorm:"default:current_timestamp" json:"updated_at,omitempty"`
}

type Filter struct {
	TransactionID string `form:"transaction_id"`
	Remark        string `form:"remark"`
	CreatedStart  string `form:"created_start"`
	CreatedEnd    string `form:"created_end"`
}
