package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheAdd(t *testing.T) {
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "http://example.com",
			val: []byte("Example 1"),
		},
		{
			key: "http://example.com/2",
			val: []byte("Example 2"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
			cache := NewCache(5 * time.Second)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Fatal("expected to find key")
			}
			if string(val) != string(c.val) {
				t.Fatal("expected value to be the same")
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(5 * time.Millisecond)
	cache.Add("http://example.com/1", []byte("data 1"))

	go func() {
		time.Sleep(8 * time.Millisecond)
		cache.Add("http://example.com/2", []byte("data 2"))
	}()
	time.Sleep(10 * time.Millisecond)

	_, ok := cache.Get("http://example.com/1")
	if ok {
		t.Errorf("expected key http://example.com/1 to have been deleted")
	}

	_, ok = cache.Get("http://example.com/2")
	if !ok {
		t.Errorf("expected key http://example.com/2 to remain")
	}
}
