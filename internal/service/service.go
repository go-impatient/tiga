package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService, NewPostService)

type Services struct {
	UserService *UserService
	PostService *PostService
}

