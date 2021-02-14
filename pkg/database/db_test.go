package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/conf/file"
)

type TestUser struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)" auth:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" auth:"email"`
}

type TestUserPromoted struct {
	TestUser
}

type TestUserOverride struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100);column:password_override" auth:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" auth:"email"`
}

type DatabaseTestSuite struct {
	suite.Suite
}

func (suite *DatabaseTestSuite) TestGetConnection() {
	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))

	err := cfg.Load()
	suite.Nil(err)

	v, err := conf.Sub("database")
	suite.Nil(err)
	dsn := v.GetString("dsn")

	fmt.Printf("mode: %s \n", dsn)
	mode := conf.Get("app.mode")
	fmt.Printf("mode: %s \n", mode)
	mode2 := conf.File("tiga").Get("app.mode")
	fmt.Printf("mode2: %s \n", mode2)
	host := conf.File("config").Get("features.nsq.host")
	fmt.Printf("host: %s \n", host)

	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)
	suite.NotNil(GetDatabase())
	suite.Equal(GetDatabase(), orm)
	suite.Nil(db.Close())
}

func (suite *DatabaseTestSuite) TestLogLevel() {
	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))
	err := cfg.Load()
	suite.Nil(err)

	dsn := conf.Get("database.dsn")
	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)

	suite.Equal(logger.Default.LogMode(logger.Info), orm.Logger)
	suite.Nil(db.Close())
}

func (suite *DatabaseTestSuite) TestSetup() {
	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))
	err := cfg.Load()
	suite.Nil(err)

	dsn := conf.Get("database.dsn")
	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)

	suite.NotNil(GetDatabase())
	suite.Equal(GetDatabase(), orm)

	ClearRegisteredModels()
	RegisterModel(&TestUser{})

	err2 := db.Migrate()
	suite.Nil(err2)

	user := &TestUser{
		Name:     "Admin",
		Password: "$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		Email:    "test@example.com",
	}

	orm.Create(user)

	suite.Nil(db.Close())
}

func (suite *DatabaseTestSuite) TestModelAndMigrate() {
	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))
	err := cfg.Load()
	suite.Nil(err)

	dsn := conf.Get("database.dsn")
	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)

	ClearRegisteredModels()
	RegisterModel(&TestUser{})
	suite.Len(models, 1)

	registeredModels := GetRegisteredModels()
	suite.Len(registeredModels, 1)
	suite.Same(models[0], registeredModels[0])

	err2 := db.Migrate()
	suite.Nil(err2)
	ClearRegisteredModels()
	suite.Equal(0, len(models))

	defer orm.Migrator().DropTable(&TestUser{})

	rows, err := orm.Raw("SHOW TABLES;").Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		if name == "test_user" {
			found = true
			break
		}
	}

	suite.True(found)
}

func (suite *DatabaseTestSuite) TestInitializers() {
	AddInitializer(func(db *gorm.DB) {
		// 跳过默认事务
		db.Config.SkipDefaultTransaction = true
	})

	suite.Len(initializers, 1)

	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))
	err := cfg.Load()
	suite.Nil(err)

	dsn := conf.Get("database.dsn")
	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)

	suite.True(orm.Config.SkipDefaultTransaction)

	AddInitializer(func(db *gorm.DB) {
		// GORM 会自动 ping 数据库以检查数据库的可用性，若要禁用该特性，可将其设置为 true
		// db.Config.DisableAutomaticPing = true
		db.Statement.Settings.Store("gorm:table_options", "ENGINE=InnoDB")
	})
	suite.Len(initializers, 2)

	val, ok := orm.Get("gorm:table_options")
	suite.True(ok)
	suite.Equal("ENGINE=InnoDB", val)

	suite.Nil(db.Close())

	ClearInitializers()
	suite.Empty(initializers)
}

func (suite *DatabaseTestSuite) TestFindColumns() {
	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))
	err := cfg.Load()
	suite.Nil(err)

	dsn := conf.Get("database.dsn")
	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)

	suite.NotNil(GetDatabase())
	suite.Equal(GetDatabase(), orm)

	user := &TestUser{}
	fields := FindColumns(user, "username", "password")
	suite.Len(fields, 2)
	// fmt.Printf("输出fields[0].Name: %#v\n", fields[0].Name)
	suite.Equal("email", fields[0].Name)
	suite.Equal("password", fields[1].Name)

	fields = FindColumns(user, "username", "notatag", "password")
	suite.Len(fields, 3)
	suite.Equal("email", fields[0].Name)
	suite.Nil(fields[1])
	suite.Equal("password", fields[2].Name)

	userOverride := &TestUserOverride{}
	fields = FindColumns(userOverride, "password")
	suite.Len(fields, 1)
	suite.Equal("password_override", fields[0].Name)

	suite.Nil(db.Close())
}

func (suite *DatabaseTestSuite) TestFindColumnsPromoted() {
	cfg := conf.New(conf.WithSource(
		file.NewFile("../../config/"),
	))
	err := cfg.Load()
	suite.Nil(err)

	dsn := conf.Get("database.dsn")
	db := NewDatabase(DSN(dsn))
	orm, err := db.Init()
	suite.Nil(err)
	suite.NotNil(orm)

	suite.NotNil(GetDatabase())
	suite.Equal(GetDatabase(), orm)

	user := &TestUserPromoted{}
	fields := FindColumns(user, "username", "password")
	suite.Len(fields, 2)
	suite.Equal("email", fields[0].Name)
	suite.Equal("password", fields[1].Name)

	fields = FindColumns(user, "username", "notatag", "password")
	suite.Len(fields, 3)
	suite.Equal("email", fields[0].Name)
	suite.Nil(fields[1])
	suite.Equal("password", fields[2].Name)

	suite.Nil(db.Close())
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
