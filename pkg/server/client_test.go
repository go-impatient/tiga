package server

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test http.Get
func TestGet(t *testing.T) {
	server_addr := "localhost:4000"
	c := NewClient()
	res, err := c.Get("http://" + server_addr)
	require.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "HandleFunc called.", string(all))
}

// Test client.Do request.
func TestDo(t *testing.T) {
	server_addr := "localhost:4000"
	c := NewClient()

	req, err := NewRequest("Get", "http://"+server_addr, nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "HandleFunc called.", string(all))
}

// Test client.Do request.
func TestDoPost(t *testing.T) {
	server_addr := "localhost:4000"
	c := NewClient()
	data := strings.NewReader(`{"data":"This is a post request."}`)
	req, err := NewRequest("Post", "http://"+server_addr+"/posttest", data)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Post test called.", string(all))
}
