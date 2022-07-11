package topup

import (
	"context"
	"errors"
	"fmt"

	"github.com/wallet-app/internal/tools"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context, f Filter) ([]Topup, error)
	GetByID(ctx context.Context, id int) (Topup, error)
	Upsert(ctx context.Context, model Topup) (Topup, error)
	DeleteByID(ctx context.Context, id int) error
	GetDB() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll(ctx context.Context, f Filter) (res []Topup, err error) {
	where := make(map[string]interface{})

	if f.TransactionID != "" {
		where["transaction_id"] = f.TransactionID
	}

	if f.Remark != "" {
		where["remark"] = f.Remark
	}

	if f.CreatedStart != "" && f.CreatedEnd != "" {
		if tools.DateValidation(f.CreatedStart, f.CreatedEnd) {
			where["to_char(created_at, 'YYYY-MM-DD') >= ?"] = f.CreatedStart
			where["to_char(created_at, 'YYYY-MM-DD') <= ?"] = f.CreatedEnd
		} else {
			return res, fmt.Errorf("invalid date format")
		}
	}

	err = r.db.Where(where).Find(&res).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return res, err
	}

	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (res Topup, err error) {
	result := r.db.First(&res, id)
	return res, result.Error
}

func (r *repository) Upsert(ctx context.Context, model Topup) (res Topup, err error) {
	result := r.db.Save(&model)
	return model, result.Error
}

func (r *repository) DeleteByID(ctx context.Context, id int) error {
	result := r.db.Unscoped().Delete(&Topup{ID: uint(id)})
	return result.Error
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}
