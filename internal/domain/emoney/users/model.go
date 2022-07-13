package users

import (
	"time"

	"gorm.io/gorm"

	"github.com/wallet-app/internal/domain/emoney/payment"
	"github.com/wallet-app/internal/domain/topup"
)

type Users struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255)"`
	Email        string `gorm:"type:varchar(100);index"`
	Username     string `gorm:"type:varchar(50);index;unique"`
	PasswordHash string
	LastLogin    time.Time
	Token        string
	Balance      int64
	Topups       []topup.Topup
	Payments     []payment.Payment
}

type Filter struct {
	Name         string `form:"name"`
	Email        string `form:"email"`
	Username     string `form:"username"`
	CreatedStart string `form:"created_start"`
	CreatedEnd   string `form:"created_end"`
}
