package cache

import (
	"errors"
	"reflect"
	"testing"
)

type mockableCacheProvider struct {
	GetFunc    func(key string) (Item, error)
	SetFunc    func(item Item) (Item, error)
	DeleteFunc func(key string) error
	HasFunc    func(key string) (bool, error)
}

func (m *mockableCacheProvider) Get(key string) (Item, error) {
	return m.GetFunc(key)
}

func (m *mockableCacheProvider) Set(item Item) (Item, error) {
	return m.SetFunc(item)
}

func (m *mockableCacheProvider) Delete(key string) error {
	return m.DeleteFunc(key)
}

func (m *mockableCacheProvider) Has(key string) (bool, error) {
	return m.HasFunc(key)
}

func TestErrors(t *testing.T) {
	cacheProvider := &mockableCacheProvider{}
	cache := Cache{
		Providers: []Provider{
			cacheProvider,
		},
	}

	reset := func() {
		cacheProvider.SetFunc = nil
		cacheProvider.DeleteFunc = nil
		cacheProvider.HasFunc = nil
		cacheProvider.GetFunc = nil
	}

	{ // Test Failed Get
		reset()
		cacheProvider.GetFunc = func(key string) (Item, error) {
			return Item{}, errors.New("Get error")
		}

		if _, err := cache.Get(""); err == nil {
			t.Errorf("cache.Get should have returned an error")
		}
	}

	{ // Test Failed Set
		reset()
		cacheProvider.SetFunc = func(item Item) (Item, error) {
			return Item{}, errors.New("Set error")
		}

		if _, err := cache.Set(Item{}); err == nil {
			t.Errorf("cache.Set should have returned an error")
		}
	}

	{ // Test Failed Delete
		reset()
		cacheProvider.DeleteFunc = func(key string) error {
			return errors.New("Delete error")
		}

		if err := cache.Delete(""); err == nil {
			t.Errorf("cache.Delete should have returned an error")
		}
	}
	{ // Test Failed Has
		reset()
		cacheProvider.HasFunc = func(key string) (bool, error) {
			return false, errors.New("Has error")
		}

		if _, err := cache.Has(""); err == nil {
			t.Errorf("cache.Has should have returned an error")
		}
	}
}

func cacheTestHelper(t *testing.T, cache Cache) {
	testKey := "key1"
	fakeKey := "super-fake-key"
	testValue := []byte("a fake value")

	{ // Test Settings Value
		_, err := cache.Set(Item{
			Key:   testKey,
			Value: testValue,
		})
		if err != nil {
			t.Error(err)
		}
	}

	{ // Test Reading a value back
		_, err := cache.Get(fakeKey)
		if err != ErrorCacheKeyNotFound {
			t.Error("did not get the ErrorCacheKeyNotFound error")
		}
	}

	{ // Test trying to get a key that does not exist
		item, err := cache.Get(testKey)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(item.Value, testValue) {
			t.Error("value was wrong")
		}
	}

	{ // Test checking if a key exists
		has, err := cache.Has(testKey)
		if err != nil {
			t.Error(err)
		}

		if !has {
			t.Error("Has check should have returned true")
		}
	}

	{ // Test checking if a key exists when it does not
		has, err := cache.Has(fakeKey)
		if err != nil {
			t.Error(err)
		}

		if has {
			t.Error("Has check should have returned false")
		}
	}

	{ // Test deleting a key
		err := cache.Delete(testKey)
		if err != nil {
			t.Error(err)
		}
	}

	{ // Test deleting a key that does not exist
		err := cache.Delete(fakeKey)
		if err != nil {
			t.Error(err)
		}
	}
}
