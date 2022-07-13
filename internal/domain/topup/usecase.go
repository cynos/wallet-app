package topup

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type UseCase interface {
	GetAll(context context.Context) ([]Topup, error)
	GetAllByFilter(context context.Context, f Filter) ([]Topup, error)
	GetByID(context context.Context, id int) (Topup, error)
	Add(context context.Context, model Topup) (Topup, error)
	Update(context context.Context, model Topup, id int) (Topup, error)
	UpdateBalance(context context.Context, amount, userid int) error
	Delete(context context.Context, id int) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (uc *usecase) GetAll(ctx context.Context) (res []Topup, err error) {
	res, err = uc.repo.GetAll(ctx, Filter{})
	return res, err
}

func (uc *usecase) GetAllByFilter(ctx context.Context, f Filter) (res []Topup, err error) {
	res, err = uc.repo.GetAll(ctx, f)
	return res, err
}

func (uc *usecase) GetByID(ctx context.Context, id int) (res Topup, err error) {
	res, err = uc.repo.GetByID(ctx, id)
	return res, err
}

func (uc *usecase) Add(ctx context.Context, model Topup) (res Topup, err error) {
	if model.TransactionID == "" {
		return res, fmt.Errorf("invalid parameters")
	}

	res, err = uc.repo.Upsert(ctx, model)
	return res, err
}

func (uc *usecase) Update(ctx context.Context, model Topup, id int) (res Topup, err error) {
	topup, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return res, fmt.Errorf("record not found")
	}

	if topup.Remark != "" {
		topup.Remark = model.Remark
	}

	if id == 0 {
		return res, fmt.Errorf("invalid parameters")
	}

	res, err = uc.repo.Upsert(ctx, topup)
	return
}

func (uc *usecase) UpdateBalance(ctx context.Context, amount, userid int) error {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"users_id": %d, "amount": %d}`, userid, amount)).
		Put("http://emoney-service:8080/users/balance")
	if err != nil {
		return fmt.Errorf("cannot update balance, detail : %v", err.Error())
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("cannot update balance, detail : %v", string(resp.Body()))
	}

	return nil
}

func (uc *usecase) Delete(ctx context.Context, id int) (err error) {
	err = uc.repo.DeleteByID(ctx, id)
	return err
}
