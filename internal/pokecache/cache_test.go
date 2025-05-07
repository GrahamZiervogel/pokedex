package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata1"),
		},
		{
			key: "https://example.com/path",
			val: []byte("testdata2"),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache.Add(tc.key, tc.val)
			val, ok := cache.Get(tc.key)
			if !ok {
				t.Errorf("expected to find key %s", tc.key)
				return
			}
			if string(val) != string(tc.val) {
				t.Errorf("expected value %s, got %s", string(tc.val), string(val))
			}
		})
	}
}

func TestCacheReap(t *testing.T) {
	const reapInterval = 50 * time.Millisecond
	const testDuration = reapInterval + (20 * time.Millisecond)
	cache := NewCache(reapInterval)

	keyToExpire := "key1"
	valToExpire := []byte("data1")

	cache.Add(keyToExpire, valToExpire)

	_, ok := cache.Get(keyToExpire)
	if !ok {
		t.Fatalf("expected key %s to be in cache initially", keyToExpire)
	}

	time.Sleep(testDuration)

	_, ok = cache.Get(keyToExpire)
	if ok {
		t.Errorf("expected key %s to be reaped from cache, but it was found", keyToExpire)
	}
}

func TestCacheReapNotTooSoon(t *testing.T) {
	const reapInterval = 100 * time.Millisecond
	const waitTime = reapInterval / 2
	cache := NewCache(reapInterval)

	keyToKeep := "keyToKeep"
	valToKeep := []byte("dataToKeep")

	cache.Add(keyToKeep, valToKeep)

	time.Sleep(waitTime)

	_, ok := cache.Get(keyToKeep)
	if !ok {
		t.Errorf("expected key %s to still be in cache, but it was reaped too soon", keyToKeep)
	}
}
