package service

import (
	"context"

	"moocss.com/tiga/internal/biz"
	"moocss.com/tiga/internal/service/dto"
	"moocss.com/tiga/pkg/log"
)

type PostService struct {
	post *biz.PostUsecase

	log *log.Helper
}

func NewPostService(post *biz.PostUsecase, logger log.Logger) *PostService {
	return &PostService{
		post: post,
		log:  log.NewHelper("post", logger),
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *dto.PostRequest) (*dto.PostResponse, error) {
	err := s.post.CreatePost(ctx, &biz.Post{
		Title:   req.Title,
		Content: req.Content,
	})
	return &dto.PostResponse{}, err
}

func (s *PostService) UpdatePost(ctx context.Context, req *dto.PostRequest) (*dto.PostResponse, error) {
	err := s.post.CreatePost(ctx, &biz.Post{
		Title:   req.Title,
		Content: req.Content,
	})
	return &dto.PostResponse{}, err
}

func (s *PostService) DeletePost(ctx context.Context, req *dto.PostRequest) (*dto.PostResponse, error) {
	return &dto.PostResponse{}, nil
}

func (s *PostService) GetPost(ctx context.Context, req *dto.PostRequest) (*dto.PostResponse, error) {
	return &dto.PostResponse{}, nil
}

func (s *PostService) ListPost(ctx context.Context, req *dto.PostRequest) (*dto.PostResponse, error) {
	return &dto.PostResponse{}, nil
}
