package biz

import (
	"context"
	"time"
)

type Post struct {
	Id      int64
	Title   string
	Content string
	CTime   time.Time
}

type PostRepo interface {
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
}

type PostUsecase struct {
	repo PostRepo
}

func NewPostUsecase(repo PostRepo) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (uc *PostUsecase) CreatePost(ctx context.Context, post *Post) error {
	return uc.repo.Create(ctx, post)
}

func (uc *PostUsecase) UpdatePost(ctx context.Context, post *Post) error {
	return uc.repo.Update(ctx, post)
}
