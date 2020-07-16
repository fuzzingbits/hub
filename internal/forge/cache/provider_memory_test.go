package cache

import "testing"

func TestCacheMemory(t *testing.T) {
	cache := Cache{
		Providers: []Provider{
			&ProviderMemory{},
		},
	}

	cacheTestHelper(t, cache)
}
