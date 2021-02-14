package conf

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"moocss.com/tiga/pkg/log"
)

var (
	files map[string]*config

	_ Config = (*config)(nil)
)

// Config is a config interface.
type Config interface {
	Load() error
	Watch()
}

type config struct {
	opts *options

	viper *viper.Viper
	log   *log.Helper
}

// New new a config with options.
func New(opts ...Option) Config {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &config{
		opts: options,
		log:  log.NewHelper("config", options.logger),
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		fs, err := src.Load()
		if err != nil {
			return err
		}
		files = make(map[string]*config, len(fs))
		for _, f := range fs {
			if !strings.HasSuffix(f, ".yaml") {
				continue
			}
			c.log.Infof("配置文件: %s", f)

			v := viper.New()
			// Config's format: "json" | "toml" | "yaml" | "yml"
			v.SetConfigType("yaml")
			v.SetConfigFile(f)
			if err := v.ReadInConfig(); err != nil {
				c.log.Warnf("Using config file: %s [%s]\n", viper.ConfigFileUsed(), err)
				panic(err)
			}

			// 读取匹配的环境变量
			v.AutomaticEnv()

			name := strings.TrimSuffix(path.Base(f), ".yaml")
			files[name] = &config{
				viper: v,
			}
		}
	}

	return nil
}

func (c *config) Watch() {
	for _, v := range files {
		v.viper.WatchConfig()
		v.viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("Config file changed: %s \n", e.Name)
		})
	}
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 { return File("tiga").GetFloat64(key) }
func (c *config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// Get 获取字符串配置
func Get(key string) string { return File("tiga").Get(key) }
func (c *config) Get(key string) string {
	return c.viper.GetString(key)
}

// GetStrings 获取字符串列表
func GetStrings(key string) (s []string) { return File("tiga").GetStrings(key) }
func (c *config) GetStrings(key string) (s []string) {
	value := Get(key)
	if value == "" {
		return
	}

	for _, v := range strings.Split(value, ",") {
		s = append(s, v)
	}
	return
}

// GetInt32s 获取数字列表
// 1,2,3 => []int32{1,2,3}
func GetInt32s(key string) (s []int32, err error) { return File("tiga").GetInt32s(key) }
func (c *config) GetInt32s(key string) (s []int32, err error) {
	s64, err := GetInt64s(key)
	for _, v := range s64 {
		s = append(s, int32(v))
	}
	return
}

// GetInt64s 获取数字列表
func GetInt64s(key string) (s []int64, err error) { return File("tiga").GetInt64s(key) }
func (c *config) GetInt64s(key string) (s []int64, err error) {
	value := Get(key)
	if value == "" {
		return
	}

	var i int64
	for _, v := range strings.Split(value, ",") {
		i, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		s = append(s, i)
	}
	return
}

// GetInt 获取整数配置
func GetInt(key string) int { return File("tiga").GetInt(key) }
func (c *config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt32 获取 int32 配置
func GetInt32(key string) int32 { return File("tiga").GetInt32(key) }
func (c *config) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

// GetInt64 获取 int64 配置
func GetInt64(key string) int64 { return File("tiga").GetInt64(key) }
func (c *config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetDuration 获取时间配置
func GetDuration(key string) time.Duration { return File("tiga").GetDuration(key) }
func (c *config) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// GetTime 查询时间配置
// 默认时间格式为 "2006-01-02 15:04:05"，conf.GetTime("FOO_BEGIN")
// 如果需要指定时间格式，则可以多传一个参数，conf.GetString("FOO_BEGIN", "2006")
//
// 配置不存在或时间格式错误返回**空时间对象**
// 使用本地时区
func GetTime(key string, args ...string) time.Time { return File("tiga").GetTime(key, args...) }
func (c *config) GetTime(key string, args ...string) time.Time {
	fmt := "2006-01-02 15:04:05"
	if len(args) == 1 {
		fmt = args[0]
	}

	t, _ := time.ParseInLocation(fmt, c.viper.GetString(key), time.Local)
	return t
}

// GetBool 获取配置布尔配置
func GetBool(key string) bool { return File("tiga").GetBool(key) }
func (c *config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// Sub 返回新的Viper实例，代表该实例的子节点。
func Sub(key string) (*viper.Viper, error) { return File("tiga").Sub(key) }
func (c *config) Sub(key string) (*viper.Viper, error) {
	if app := c.viper.Sub(key); app != nil {
		return app, nil
	}
	return nil, fmt.Errorf("No found `%s` in the configuration", key)
}

// Set 设置配置，仅用于测试
func Set(key string, value string) { File("tiga").Set(key, value) }
func (c *config) Set(key string, value string) {
	c.viper.Set(key, value)
}

// File 根据文件名获取对应配置对象
// 目前仅支持 toml 文件，不用传扩展名
// 如果要读取 foo.toml 配置，可以 File("foo").Get("bar")
func File(name string) *config {
	res, _ := files[name]
	if res == nil {
		res = &config{viper: &viper.Viper{}}
	}
	return res
}
