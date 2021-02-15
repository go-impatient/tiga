package data

import (
	"context"

	"moocss.com/tiga/internal/biz"
	"moocss.com/tiga/pkg/log"
)

type postRepo struct {
	data *Data
	log  *log.Helper
}

// NewPostRepo .
func NewPostRepo(data *Data, logger log.Logger) biz.PostRepo {
	return &postRepo{
		data: data,
		log:  log.NewHelper("post_repo", logger),
	}
}

func (r *postRepo) Create(ctx context.Context, post *biz.Post) error {
	return nil
}

func (r *postRepo) Update(ctx context.Context, post *biz.Post) error {
	return nil
}