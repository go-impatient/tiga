package data

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/database"
	"moocss.com/tiga/pkg/log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewPostRepo)

// Data .
type Data struct {
	*gorm.DB
}

// NewData .
func NewData(logger log.Logger) (*Data, error) {
	log := log.NewHelper("data", logger)

	dsn := conf.Get("database.dsn")
	dialect := conf.Get("database.dialect")
	db := database.NewDatabase(
		database.Dialect(dialect),
		database.DSN(dsn),
	)

	instance, err := db.Init()

	defer db.Close()

	if err != nil {
		log.Errorf("failed opening connection to %s: %v", dialect, err)
		return nil, err
	}

	return &Data{DB: instance}, nil
}
