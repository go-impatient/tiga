package dto

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-impatient/gaia/pkg/errcode"
)

func TestResponseNew(t *testing.T) {
	type Test struct {
		Id   int
		Name string
	}
	code := 200
	msg := "have some error"
	data := &Test{Id: 1, Name: "test-name"}
	err := errcode.New("gaia.api.error", int32(code), errcode.Msg(msg), errcode.Data(data))
	assert.NotNil(t, err)
	assert.Equal(t, "go-error: code = 200 ,message = have some error ,data = &{%!s(int=1) test-name}", err.Errorf())
	e:= errcode.FromError(err)
	assert.NotNil(t, e)
	assert.Equal(t, int32(code), e.Code)
	assert.Equal(t, msg, e.Message)
	assert.Equal(t, data, e.Data)
	assert.Nil(t, e.Extra)
}

func TestError(t *testing.T) {
	code1 := 500
	msg1 := "unknown error"
	err := errcode.New("gaia.api.error", int32(code1), errcode.Msg(msg1))
	e := errcode.FromError(err)
	assert.NotNil(t, e)
	assert.EqualValues(t, code1, e.Code)
	assert.EqualValues(t, msg1, e.Message)

	err = errcode.New("gaia.api.error", 500, errcode.Msg(msg1))
	e = errcode.FromError(err)
	assert.NotNil(t, e)
	assert.EqualValues(t, http.StatusInternalServerError, e.Code)
	assert.EqualValues(t, msg1, e.Message)
}

func TestErrorf(t *testing.T) {
	type Test struct {
		Id   int
		Name string
	}
	type ExtraTest struct {
		TotalCount int
	}
	code1 := 301
	msg1 := "unknown error"
	data1 := &Test{Id: 1, Name: "test-name"}
	extra1 := &ExtraTest{TotalCount: 500}
	err := errcode.New("gaia.api.error", int32(code1), errcode.Msg(msg1), errcode.Data(data1), errcode.Extra(extra1))
	e := errcode.FromError(err)
	assert.NotNil(t, e)
	assert.EqualValues(t, code1, e.Code)
	assert.EqualValues(t, msg1, e.Message)
	assert.EqualValues(t, data1, e.Data)
	assert.EqualValues(t, extra1, e.Extra)

	msg2 := "test error"
	err = errcode.New("gaia.api.error", 500, errcode.Msg("test error"))
	e = errcode.FromError(err)
	assert.NotNil(t, e)
	assert.EqualValues(t, http.StatusInternalServerError, e.Code)
	assert.EqualValues(t, msg2, e.Message)
}
