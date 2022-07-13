package payment

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (Payment, error)
	GetByTrx(ctx context.Context, trx string) (Payment, error)
	Upsert(ctx context.Context, model Payment) (Payment, error)
	GetDB() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByID(ctx context.Context, id int) (res Payment, err error) {
	result := r.db.Where("id = ?", id).Find(&res)
	return res, result.Error
}

func (r *repository) GetByTrx(ctx context.Context, trx string) (res Payment, err error) {
	result := r.db.Where("transaction_id = ?", trx).Find(&res)
	return res, result.Error
}

func (r *repository) Upsert(ctx context.Context, model Payment) (res Payment, err error) {
	result := r.db.Save(&model)
	return model, result.Error
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}
