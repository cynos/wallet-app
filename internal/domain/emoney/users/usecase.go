package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/wallet-app/internal/tools"
)

type UseCase interface {
	GetUserAccount(ctx context.Context, id int) (Users, error)
	GenerateUser(ctx context.Context, dto DTOUsers) (Users, error)
	AuthenticateUser(ctx context.Context, dto DTOUsers) (Users, error)
	Update(ctx context.Context, model Users, id int) (Users, error)
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo}
}

func (uc *usecase) GetUserAccount(ctx context.Context, id int) (res Users, err error) {
	res, err = uc.repo.GetByID(ctx, id)
	return
}

func (uc *usecase) GenerateUser(ctx context.Context, dto DTOUsers) (res Users, err error) {
	dto.Username = strings.TrimSpace(dto.Username)
	dto.Username = strings.ToLower(dto.Username)

	_, err = uc.repo.GetAll(ctx, Filter{Username: dto.Username})
	if err != nil {
		return res, fmt.Errorf("username already exist | detail : %v", err)
	}

	hash, err := tools.HashPassword(dto.Password)
	if err != nil {
		return res, fmt.Errorf("cannot hash password | detail : %v", err)
	}

	res, err = uc.repo.Upsert(ctx, Users{
		Name:         dto.Name,
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: hash,
	})
	if err != nil {
		return res, fmt.Errorf("cannot generate users | detail : %s", err.Error())
	}

	return res, nil
}

func (uc *usecase) AuthenticateUser(ctx context.Context, dto DTOUsers) (res Users, err error) {
	dto.Username = strings.TrimSpace(dto.Username)
	dto.Username = strings.ToLower(dto.Username)

	result, err := uc.repo.GetAll(ctx, Filter{Username: dto.Username})
	if err != nil || len(result) < 1 {
		return res, fmt.Errorf("username not found | detail : %v", err)
	}

	match := tools.CheckPasswordHash(dto.Password, result[0].PasswordHash)
	if !match {
		return res, fmt.Errorf("cannot chech password hash | detail : %v", err)
	}

	return result[0], nil
}

func (uc *usecase) Update(ctx context.Context, model Users, id int) (res Users, err error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return res, fmt.Errorf("record not found")
	}

	if model.Name != "" {
		user.Name = model.Name
	}

	if model.Email != "" {
		user.Email = model.Email
	}

	if model.Username != "" {
		user.Username = model.Username
	}

	if model.PasswordHash != "" {
		user.PasswordHash = model.PasswordHash
	}

	if model.LastLogin.Unix() != 0 {
		user.LastLogin = model.LastLogin
	}

	if model.Token != "" {
		user.Token = model.Token
	}

	if model.Balance != 0 {
		user.Balance = model.Balance
	}

	res, err = uc.repo.Upsert(ctx, user)
	return
}
