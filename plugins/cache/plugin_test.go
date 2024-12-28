package cache

import (
	"testing"

	kv "github.com/philippgille/gokv/redis"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSetGet(t *testing.T) {
	key := "test-setget-key"
	val := "test-setget-val"
	log := zap.NewNop().Sugar()
	opts := &kv.Options{
		Address: "localhost:6379",
	}
	c, err := New(log, opts)
	require.NoError(t, err)
	defer c.Close()
	defer c.Delete(key)

	err = c.Set(key, val)
	require.NoError(t, err)

	var v string
	ok, err := c.Get(key, &v)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, val, v)
}

func TestFetch(t *testing.T) {
	key := "test-fetch-key"
	val := 1
	log := zap.NewNop().Sugar()
	opts := &kv.Options{
		Address: "localhost:6379",
	}
	c, err := New(log, opts)
	require.NoError(t, err)
	defer c.Close()
	defer c.Delete(key)

	i := 0
	f := func() (interface{}, error) {
		i++
		return i, nil
	}

	var v int
	ok, err := c.Fetch(key, &v, f)
	require.False(t, ok)
	require.NoError(t, err)
	require.Equal(t, val, v)

	ok, err = c.Fetch(key, &v, f)
	require.True(t, ok)
	require.NoError(t, err)
	require.Equal(t, val, v)
}
