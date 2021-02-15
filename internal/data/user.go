package user

import (
	"context"

	"gorm.io/gorm"

	"moocss.com/tiga/internal/biz"
	"moocss.com/tiga/pkg/log"
)

// userRepository struct
type userRepo struct {
	db  *gorm.DB
	log *log.Helper
}

// NewUserRepo returns an instance of `userRepo`.
func NewUserRepo(db *gorm.DB, logger log.Logger) biz.UserRepo {
	return &userRepo{
		db:  db,
		log: log.NewHelper("user_repository", logger),
	}
}

// implement biz.UserRepo
func (r *userRepo) Exist(ctx context.Context, user *biz.User) (bool, error) {
	r.log.Info("Received UserRepository.Exist")

	return true, nil
}
func (r *userRepo) List(ctx context.Context, limit, page int, sort string, user *biz.User) (total int, users []*biz.User, err error) {
	r.log.Info("Received UserRepository.List")
	return 0, nil, nil
}

func (r *userRepo) Get(ctx context.Context, id int) (*biz.User, error) {
	r.log.Info("Received UserRepository.Get")
	return nil, nil
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Info("Received UserRepository.Create")
	return nil, nil
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Info("Received UserRepository.Update")
	return nil, nil
}

func (r *userRepo) DeleteFull(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Info("Received UserRepository.DeleteFull")
	return nil, nil
}

func (r *userRepo) Delete(ctx context.Context, id int) (*biz.User, error) {
	r.log.Info("Received UserRepository.Delete")
	return nil, nil
}

func (r *userRepo) Count(ctx context.Context) (int, error) {
	r.log.Info("Received UserRepository.Count")
	return 0, nil
}
