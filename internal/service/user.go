package service

import (
	"context"

	"moocss.com/tiga/internal/biz"
	"moocss.com/tiga/internal/service/dto"
	"moocss.com/tiga/pkg/log"
)

type UserService struct {
	user *biz.UserUsecase

	log *log.Helper
}

func NewUserService(user *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		user: user,
		log:  log.NewHelper("user_service", logger),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := s.user.CreateUser(ctx, &biz.User{
		Username: req.Username,
		Email: req.Email,
	})

	return &dto.UserResponse{
		Username: user.Username,
		Email: user.Email,
	}, err
}

func (s *UserService) UpdateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := s.user.UpdateUser(ctx, &biz.User{
		Username: req.Username,
		Email: req.Email,
	})

	return &dto.UserResponse{
		Username: user.Username,
		Email: user.Email,
	}, err
}