package data

import (
	"context"

	"github.com/pkg/errors"

	"moocss.com/tiga/internal/biz"
	"moocss.com/tiga/pkg/log"
)

// userRepository struct
type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo returns an instance of `userRepo`.
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper("user_repository", logger),
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
	_, ok := r.data.DB.Get(string(id))
	if !ok {
		return nil, errors.New("查询成功")
	}

	return nil, nil
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Info("Received UserRepository.Create")

	po := user.UserToPO()
	if err := r.data.DB.Create(po).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Info("Received UserRepository.Update")
	// po := user.UserToPO()
	// if err := r.data.DB.Model(po).Where("id = ? AND is_deleted = ?", user.ID, 0).Updates(po).Error; err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func (r *userRepo) DeleteFull(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Info("Received UserRepository.DeleteFull")
	return nil, nil
}

func (r *userRepo) Delete(ctx context.Context, id int) (*biz.User, error) {
	r.log.Info("Received UserRepository.Delete")
	// if err := r.data.DB.Where("id = ? AND is_deleted = ?", id, 0).Delete(schema.User{}).Error; err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func (r *userRepo) Count(ctx context.Context) (int, error) {
	r.log.Info("Received UserRepository.Count")
	return 0, nil
}
