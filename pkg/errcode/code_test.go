package errcode

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsNew(t *testing.T) {
	code := 200
	message := "This is a message"
	data := map[string]string{
		"id":   "1",
		"name": "test-name",
	}
	result := New("app.error", int32(code), Msg(message), Data(data))
	assert.Equal(t, int32(code), result.Code)
	assert.Equal(t, message, result.Message)
	assert.Equal(t, data, result.Data)
	assert.Nil(t, result.Extra)

	t.Logf("输出错误: %v", result.Error())
}

func TestFromError(t *testing.T) {
	err := NotFound("app.user.test", "%s", "example")
	merr := FromError(err)
	t.Logf("输出merr: %v", merr.Error())
	if merr.Id != "app.user.test" || merr.Code != 404 {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
	err = errors.New(err.Error())
	merr = FromError(err)
	t.Logf("输出merr: %v", merr.Error())
	if merr.Id != "app.user.test" || merr.Code != 404 {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
}

func TestEqual(t *testing.T) {
	err1 := NotFound("myid1", "msg1")
	err2 := NotFound("myid2", "msg2")

	if !Equal(err1, err2) {
		t.Fatal("errors must be equal")
	}

	err3 := errors.New("my test err")
	if Equal(err1, err3) {
		t.Fatal("errors must be not equal")
	}

}

func TestErrors(t *testing.T) {
	testData := []*Error{
		{
			Id:      "test",
			Code:    500,
			Message: "Internal server error",
		},
		{
			Id:      "test2",
			Code:    401,
			Message: "用户没有访问权限",
		},
	}

	for _, e := range testData {
		ne := New(e.Id, e.Code, Msg(e.Message))

		if e.Error() != ne.Error() {
			t.Fatalf("Expected %s got %s", e.Error(), ne.Error())
		}

		pe := Parse(ne.Error())

		if pe == nil {
			t.Fatalf("Expected error got nil %v", pe)
		}

		if pe.Id != e.Id {
			t.Fatalf("Expected %s got %s", e.Id, pe.Id)
		}

		if pe.Message != e.Message {
			t.Fatalf("Expected %s got %s", e.Message, pe.Message)
		}

		if pe.Code != e.Code {
			t.Fatalf("Expected %d got %d", e.Code, pe.Code)
		}
	}
}
