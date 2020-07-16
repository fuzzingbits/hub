package cache

import (
	"os"
	"testing"
)

func TestCacheJSON(t *testing.T) {
	// Define Variables
	testPath := "provider_json_git_excluded"

	{ // Setup
		if err := os.RemoveAll(testPath); err != nil {
			t.Error(err)
		}

		if err := os.MkdirAll(testPath, 0700); err != nil {
			t.Error(err)
		}
	}

	cache := Cache{
		Providers: []Provider{
			ProviderJSON{
				Path:          testPath,
				HashFilenames: true,
			},
		},
	}

	cacheTestHelper(t, cache)
}

func TestCacheJSONSet(t *testing.T) {
	cache := Cache{
		Providers: []Provider{
			ProviderJSON{
				Path:          "not/real/path",
				HashFilenames: true,
			},
		},
	}

	if _, err := cache.Set(Item{Key: "foobar"}); err == nil {
		t.Errorf("cache.Set should have returned an error")
	}
}
func TestCacheJSONGet(t *testing.T) {
	cache := Cache{
		Providers: []Provider{
			ProviderJSON{
				Path: "provider_json_test_files",
			},
		},
	}

	if _, err := cache.Get("malformed"); err == nil {
		t.Errorf("cache.Get should have returned an error")
	}
}

func TestCacheJSONGetNoPath(t *testing.T) {
	cache := Cache{
		Providers: []Provider{
			ProviderJSON{
				Path: "",
			},
		},
	}

	if _, err := cache.Get("malformed"); err == nil {
		t.Errorf("cache.Get should have returned an error")
	}
}
