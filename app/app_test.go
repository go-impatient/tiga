package app

import (
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	app := New(
		Version("v1.0.0"),
	)
	time.AfterFunc(time.Second, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
