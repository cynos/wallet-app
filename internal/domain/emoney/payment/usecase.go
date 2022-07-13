package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/wallet-app/internal/tools"
)

type UseCase interface {
	GetPayment(ctx context.Context, userID uint, trxid string) (Payment, error)
	CreatePayment(ctx context.Context, userID uint, productID string, amount int, remark string) (Payment, error)
	ConfirmPayment(ctx context.Context, userID uint, trxid string) (Payment, error)
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (uc *usecase) GetPayment(ctx context.Context, userID uint, trxid string) (res Payment, err error) {
	if trxid == "" {
		return res, fmt.Errorf("invalid parameter")
	}

	res, err = uc.repo.GetByTrx(ctx, trxid)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc *usecase) CreatePayment(ctx context.Context, userID uint, productID string, amount int, remark string) (res Payment, err error) {
	if userID == 0 {
		return res, fmt.Errorf("cannot create payment, invalid usersid productid")
	}

	if productID == "" {
		return res, fmt.Errorf("cannot create payment, invalid parameter productid")
	}

	if amount == 0 {
		return res, fmt.Errorf("cannot create payment, invalid parameter amount")
	}

	if remark == "" {
		return res, fmt.Errorf("cannot create payment, invalid parameter remark")
	}

	model := Payment{
		UsersID:       userID,
		TransactionID: tools.RandString(20),
		ProductID:     productID,
		Amount:        amount,
		Remark:        remark,
	}

	res, err = uc.repo.Upsert(ctx, model)
	if err != nil {
		return res, fmt.Errorf("cannot create payment, detail : %v", err.Error())
	}

	return res, nil
}

func (uc *usecase) ConfirmPayment(ctx context.Context, userID uint, trxid string) (res Payment, err error) {
	if trxid == "" {
		return res, fmt.Errorf("cannot confirm payment, invalid parameter")
	}

	res, err = uc.repo.GetByTrx(ctx, trxid)
	if err != nil {
		return res, fmt.Errorf("cannot confirm payment, detail : %v", err.Error())
	}

	if res.Confirm == 1 {
		return res, fmt.Errorf("cannot confirm payment, detail : %v", "transaction has been confirmed")
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("InternalConnection", "1").
		SetHeader("UserID", fmt.Sprint(userID)).
		Get("http://emoney-service:8080/users/balance")
	if err != nil {
		return res, fmt.Errorf("cannot get balance, detail : %v", err.Error())
	}

	p := struct {
		Message string `json:"message"`
		Result  bool   `json:"result"`
		Data    struct {
			Balance int `json:"balance"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(resp.Body(), &p)
	if err != nil {
		return res, fmt.Errorf("cannot get balance, detail : %v", err.Error())
	}

	if !p.Result {
		return res, fmt.Errorf("cannot get balance, detail : %v", p.Message)
	}

	if p.Data.Balance < res.Amount {
		return res, fmt.Errorf("cannot confirm payment, detail : %v", "insufficent balance")
	}

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"users_id": %d, "amount": %d}`, userID, -res.Amount)).
		Put("http://emoney-service:8080/users/balance")
	if err != nil {
		return res, fmt.Errorf("cannot confirm payment, failed deduct balance, detail : %v", err.Error())
	}

	if resp.StatusCode() != 200 {
		return res, fmt.Errorf("cannot confirm payment, failed deduct balance, detail : %v", resp.String())
	}

	model := res
	model.Confirm = 1
	res, err = uc.repo.Upsert(ctx, model)
	if err != nil {
		return res, fmt.Errorf("cannot confirm payment, detail : %v", err.Error())
	}

	return res, nil
}
