package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		c.Clear()

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("overflow", func(t *testing.T) {
		capacity := 3
		c := NewCache(capacity)

		for i := 0; i < capacity+1; i++ {
			_ = c.Set(Key(strconv.Itoa(i)), i)
		}

		val, ok := c.Get("1")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("2")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("3")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("0")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("least used out", func(t *testing.T) {
		capacity := 3
		c := NewCache(capacity)

		for i := 0; i < capacity; i++ {
			_ = c.Set(Key(strconv.Itoa(i)), i)
		} // [0, 1, 2]

		_ = c.Set("2", 22) // [2, 0, 1]
		_, _ = c.Get("0")  // [0, 2, 1]
		_ = c.Set("1", 11) // [1, 0, 2]
		_ = c.Set("0", 10) // [0, 1, 2]

		_ = c.Set("3", 3) // [3, 0, 1] 2 is out

		val, ok := c.Get("2")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("zero capacity", func(t *testing.T) {
		capacity := 0
		tKey := Key("1")
		c := NewCache(capacity)

		_ = c.Set(tKey, 1)

		val, ok := c.Get(tKey)
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
