package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/wallet-app/internal/tools"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context, f Filter) ([]Users, error)
	GetByID(ctx context.Context, id int) (Users, error)
	Upsert(ctx context.Context, model Users) (Users, error)
	GetDB() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll(ctx context.Context, f Filter) (res []Users, err error) {
	where := make(map[string]interface{})

	if f.Email != "" {
		where["email"] = f.Email
	}

	if f.Name != "" {
		where["name"] = f.Name
	}

	if f.Username != "" {
		where["username"] = f.Username
	}

	if f.CreatedStart != "" && f.CreatedEnd != "" {
		if tools.DateValidation(f.CreatedStart, f.CreatedEnd) {
			where["to_char(created_at, 'YYYY-MM-DD') >= ?"] = f.CreatedStart
			where["to_char(created_at, 'YYYY-MM-DD') <= ?"] = f.CreatedEnd
		} else {
			return res, fmt.Errorf("invalid date format")
		}
	}

	err = r.db.Preload("Topups").Where(where).Find(&res).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return res, err
	}

	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (res Users, err error) {
	result := r.db.Preload("Topups").Preload("Payments").Where("id = ?", id).Find(&res)
	return res, result.Error
}

func (r *repository) Upsert(ctx context.Context, model Users) (res Users, err error) {
	result := r.db.Save(&model)
	return model, result.Error
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}
