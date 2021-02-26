//+build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"moocss.com/tiga/app"

	"moocss.com/tiga/internal/biz"
	"moocss.com/tiga/internal/data"
	"moocss.com/tiga/internal/service"
	"moocss.com/tiga/pkg/log"
)

//go:generate wire
// InitApp init application dependency injection.
func InitApp(logger log.Logger) (*app.App, error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		wire.Struct(new(service.Services), "*"),
		newApp,
	))
}
