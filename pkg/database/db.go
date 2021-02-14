package database

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/slice"
)

var (
	defaultSQL *gorm.DB
	sqlMap     sync.Map

	mu           sync.Mutex
	models       []interface{}
	initializers []Initializer

	_ Database = (*database)(nil)
)

// 数据库驱动名称
type DriverName string

const (
	SQLITE   DriverName = "sqlite"
	MYSQL    DriverName = "mysql"
	POSTGRES DriverName = "postgres"
)

// AsDefault alias for "default"
const AsDefault = "default"

// Initializer is a function meant to modify a connection settings
// at the global scope when it's created.
//
// Use `db.InstantSet()` and not `db.Set()`, since the latter clones
// the gorm.DB instance instead of modifying it.
type Initializer func(*gorm.DB)

// Column matches a column name with a struct field.
type Column struct {
	Name  string
	Field *reflect.StructField
}

type Database interface {
	// 初始化数据库
	Init() (*gorm.DB, error)
	// 关闭数据库
	Close() error
	// 数据库迁移
	Migrate() error
	// 创建表
	CreateTables() error
	// 删除某些表
	DeleteTables(Models []interface{}) error
}

// database 基本数据结构
type database struct {
	opts *options
	orm  *gorm.DB
	log  *log.Helper
}

// NewDatabase new a database with options.
func NewDatabase(opts ...Option) Database {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &database{
		opts: options,
		log:  log.NewHelper("database", options.logger),
	}
}

func (db *database) Init() (*gorm.DB, error) {
	orm, err := db.connected()
	if err != nil {
		return nil, err
	}

	// Initializer functions are meant to modify a connection settings
	for _, initializer := range initializers {
		initializer(orm)
	}

	defaultSQL = orm
	sqlMap.Store("default", orm)

	db.orm = orm

	return orm, nil
}

// Close the database connections if they exist.
func (db *database) Close() error {
	mu.Lock()
	defer mu.Unlock()

	sqlDB, err := db.orm.DB()
	if err != nil {
		return errors.Wrap(err, "disconnect from database failed")
	}
	return sqlDB.Close()
}

// Connecte 连接数据库
func (db *database) connected() (*gorm.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	dialect, err := db.registerDialect(db.getDriverName(db.opts.dialect), db.opts.dsn)
	if err != nil {
		return nil, errors.Wrap(err, "database driver failed")
	}

	logLevel := logger.Silent
	if db.getMode() {
		logLevel = logger.Info
	}

	orm, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名, 例如: 't_user'
		},
	})

	// 尝试多连接几次, 是否能连接成功
	if err != nil || orm == nil {
		for i := 1; i <= 10; i++ {
			db, err := gorm.Open(dialect, &gorm.Config{
				Logger: logger.Default.LogMode(logLevel),
			})

			if db != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}
		db.log.Errorf("database connection failed: [%v]", err)
		return nil, errors.Wrap(err, "database connection failed")
	}

	// 数据库调优
	if sqlDB, err := orm.DB(); err == nil {
		// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		sqlDB.SetMaxIdleConns(db.getMaxIdleConnection())
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(db.getMaxOpenConnection())

		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
	}

	return orm, nil
}

// 注册方言
func (db *database) registerDialect(driver DriverName, dsn string) (gorm.Dialector, error) {
	switch driver {
	case SQLITE:
		return sqlite.Open(dsn), nil
	case MYSQL:
		return mysql.Open(dsn), nil
	case POSTGRES:
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("no database driver named `%s` found", driver)
	}
}

// 模式
func (db *database) getMode() bool {
	mode := conf.Get("app.mode")
	if mode != "" {
		if mode == "dev" {
			return true
		}
		return false
	}
	return true
}

func (db *database) getDriverName(driverName string) DriverName {
	switch strings.ToLower(driverName) {
	case "sqlite":
		return SQLITE
	case "mysql":
		return MYSQL
	case "postgres":
		return POSTGRES
	default:
		return SQLITE
	}
}

// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
func (db *database) getMaxOpenConnection() int {
	limit := db.opts.maxOpenConns

	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}

	if limit > 1024 {
		limit = 1024
	}

	return limit
}

// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
func (db *database) getMaxIdleConnection() int {
	limit := db.opts.maxIdleConns

	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}

	if limit > db.getMaxOpenConnection() {
		limit = db.getMaxOpenConnection()
	}

	return limit
}

// GetDatabase 全局的数据库服务
func GetDatabase(name ...string) *gorm.DB {
	if len(name) == 0 || name[0] == AsDefault {
		if defaultSQL == nil {
			fmt.Errorf("Invalid db `%s` \n", AsDefault)
		}
		return defaultSQL
	}

	v, ok := sqlMap.Load(name[0])
	if !ok {
		fmt.Errorf("Invalid db `%s` \n", name[0])
	}

	return v.(*gorm.DB)
}

// Migrate migrates all registered models.
func (db *database) Migrate() error {
	if err := db.orm.AutoMigrate(models...); err != nil {
		return errors.Wrap(err, "bcrypt migrate tables failed")
	}

	return nil
}

// creates necessary database tables
func (db *database) CreateTables() error {
	for _, model := range models {
		if !db.orm.Migrator().HasTable(model) {
			if err := db.orm.Migrator().CreateTable(model); err != nil {
				return errors.Wrap(err, "create table failed")
			}
		}
	}

	return nil
}

func (db *database) DeleteTables(Models []interface{}) error {
	if err := db.orm.Migrator().DropTable(Models...); err != nil {
		return errors.Wrap(err, "delete table failed")
	}
	return nil
}

// AddInitializer adds a database connection initializer function.
// Initializer functions are meant to modify a connection settings
// at the global scope when it's created.
//
// Initializer functions are called in order, meaning that functions
// added last can override settings defined by previous ones.
func AddInitializer(initializer Initializer) {
	initializers = append(initializers, initializer)
}

// ClearInitializers remove all database connection initializer functions.
func ClearInitializers() {
	initializers = []Initializer{}
}

// RegisterModel registers a model for auto-migration.
// When writing a model file, you should always register it in the init() function.
//  func init() {
//		database.RegisterModel(&MyModel{})
//  }
func RegisterModel(model interface{}) {
	models = append(models, model)
}

// GetRegisteredModels get the registered models.
// The returned slice is a copy of the original, so it
// cannot be modified.
func GetRegisteredModels() []interface{} {
	return append(make([]interface{}, 0, len(models)), models...)
}

// ClearRegisteredModels unregister all models.
func ClearRegisteredModels() {
	models = []interface{}{}
}

// FindColumns in given struct. A field matches if it has a "auth" tag with the given value.
// Returns a slice of found fields, ordered as the input "fields" slice.
// If the nth field is not found, the nth value of the returned slice will be nil.
//
// Promoted fields are matched as well.
//
// Given the following struct and "username", "notatag", "password":
//  type TestUser struct {
// 		gorm.Model
// 		Name     string `gorm:"type:varchar(100)"`
// 		Password string `gorm:"type:varchar(100)" auth:"password"`
// 		Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
//  }
//
// The result will be the "Email" field, "nil" and the "Password" field.
func FindColumns(strct interface{}, fields ...string) []*Column {
	length := len(fields)
	result := make([]*Column, length)

	value := reflect.ValueOf(strct)
	t := reflect.TypeOf(strct)
	if t.Kind() == reflect.Ptr {
		value = value.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := value.Field(i)
		fieldType := t.Field(i)
		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			// Check promoted fields recursively
			for i, v := range FindColumns(field.Interface(), fields...) {
				if v != nil {
					result[i] = v
				}
			}
			continue
		}

		tag := fieldType.Tag.Get("auth")
		if index := slice.IndexOf(fields, tag); index != -1 {
			result[index] = &Column{
				Name:  columnName(&fieldType),
				Field: &fieldType,
			}
		}
	}

	return result
}

func columnName(field *reflect.StructField) string {
	for _, t := range strings.Split(field.Tag.Get("gorm"), ";") { // Check for gorm column name override
		if strings.HasPrefix(t, "column") {
			v := strings.Split(t, ":")
			return strings.TrimSpace(v[1])
		}
	}

	return GetDatabase().NamingStrategy.ColumnName("", field.Name)
}
