package cache

import (
	"errors"
	"time"
)

// ErrorCacheKeyNotFound is foobar
var ErrorCacheKeyNotFound = errors.New("Cache Key Not Found")

// Item is foobar
type Item struct {
	Key       string
	Value     []byte
	ExpiresAt *time.Time
}

// Provider is foobar
type Provider interface {
	Get(key string) (Item, error)
	Set(item Item) (Item, error)
	Delete(key string) error
	Has(key string) (bool, error)
}

// Cache is foobar
type Cache struct {
	Providers []Provider
}

// Get the first CacheItem found by a provider
func (c Cache) Get(key string) (Item, error) {
	for _, provider := range c.Providers {
		item, err := provider.Get(key)
		if err != nil {
			// Skip if the error is that the key does not exist
			if err == ErrorCacheKeyNotFound {
				continue
			}

			return Item{}, err
		}

		return item, nil
	}

	return Item{}, ErrorCacheKeyNotFound
}

// Set the CacheItem in all of the providers
func (c Cache) Set(item Item) (Item, error) {
	var cacheItem Item
	var err error

	for _, provider := range c.Providers {
		cacheItem, err = provider.Set(item)
		if err != nil {
			return Item{}, err
		}
	}

	return cacheItem, nil
}

// Delete the CacheItem in all of the providers
func (c Cache) Delete(key string) error {
	for _, provider := range c.Providers {
		err := provider.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// Has will return true if any of the providers have the key
func (c Cache) Has(key string) (bool, error) {
	for _, provider := range c.Providers {
		has, err := provider.Has(key)
		if err != nil {
			return false, err
		}

		if !has {
			continue
		}

		return true, nil
	}

	return false, nil
}
