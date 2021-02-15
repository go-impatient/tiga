/*
# stats
**stats base on golang prometheus client**

* essy api
* support gin middleware, fix dymamic params in url

![stats.png](stats.png)

### methods

**http**

* SetHttpReqStats
* SetHttpReqStatsWrap

**func**

* NewFuncDurationStats
* SetFuncDurationStatsWrap
* SetFuncDurationStats

**database**

* NewDatabaseDurationStats
* SetDatabaseDurationStats
* SetDatabaseDurationStatsWrap

[`more...`](stats.go)

### gin usage:

**metrics**

* metric:http_request_duration_seconds
* metric:http_request_total
* metric:http_response_bytes
* metric:http_request_bytes

**how to use in gin web frame ?**

```
router := gin.Default()

// add gin metric middleware
router.Use(stats.GinMetricMiddleware())

router.GET("/metrics", gin.WrapH(promhttp.Handler()))
router.Run(":80")
```

### example usage:

func NewDatabaseDurationStats(dao string, filter string, args []interface{}) func()

```go
func() {
    var (
        dao = "user"
        filter = "queryUserByIds"
    )

    done := NewDatabaseDurationStats(dao, filter, []string{"admin"})

    time.Sleep(1e6)
    time.Sleep(1e6)
    time.Sleep(1e6)

    done()
}
```

func SetDatabaseDurationStats(dao string, filter string, args []interface{}, start time.Time)

```go
func() {
    var (
        dao = "user"
        filter = "queryUserByIds"
        start = time.Now()
    )

    time.Sleep(1e6)
    time.Sleep(1e6)
    time.Sleep(1e6)

    SetDatabaseDurationStats(dao, filter, nil, start)
}
```

func SetDatabaseDurationStatsWrap(dao string, filter string, args []interface{}) func()

```go
func() {
    var (
        dao = "user"
        filter = "queryUserByIds"
    )

    defer SetDatabaseDurationStatsWrap(dao, filter, nil)()

    time.Sleep(1e6)
    time.Sleep(1e6)
    time.Sleep(1e6)
}
```
*/

package middleware

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestFormatArgs(t *testing.T) {
	args := []interface{}{"name", "addr", 88}
	v := formatArgs(args)
	assert.Equal(t, v, "name_addr_88")
}

func TestTimeSince(t *testing.T) {
	start := time.Now()
	ds := timeSince(start)
	assert.Equal(t, ds, 0.001) // min val
}

func TestSimple(t *testing.T) {
	defer SetFuncDurationStatsWrap("test_func_1", nil)()
	start := time.Now()
	done := NewFuncDurationStats("test_func_2", nil)
	time.Sleep(1 * time.Second)
	SetFuncDurationStats("test_func_3", nil, start)
	done()
}

// test gin
func TestGin(t *testing.T) {
	router := gin.Default()

	// add gin metric middleware
	router.Use(GinMetricMiddleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go router.Run(":8888")

	time.AfterFunc(10*time.Second, unregister)
	time.Sleep(20 * time.Second)
}
