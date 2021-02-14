package conf

import (
	"moocss.com/tiga/pkg/conf/file"
	"testing"
)

func TestConfig(t *testing.T) {
	c := New(WithSource(
		file.NewFile("../../config/"),
		// file.NewFile("../../config/tiga.yaml"),
		// file.NewFile("../../"),
	))

	testConfig(t, c)
}

func testConfig(t *testing.T, c Config) {
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	v, _ := Sub("database")
	dsn := v.GetString("dsn")
	t.Logf("mode: %s", dsn)
	mode := Get("app.mode")
	t.Logf("mode: %s", mode)
	mode2 := File("tiga").Get("app.mode")
	t.Logf("mode2: %s", mode2)
	host := File("config").Get("features.nsq.host")
	t.Logf("host: %s", host)
}
