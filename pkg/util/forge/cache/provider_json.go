package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ProviderJSON is foobar
type ProviderJSON struct {
	Path          string
	HashFilenames bool
}

// Get is foobar
func (p ProviderJSON) Get(key string) (Item, error) {
	filename := p.getFilename(key)
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Item{}, ErrorCacheKeyNotFound
	}

	item := Item{}
	err = json.Unmarshal(jsonBytes, &item)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

// Set is foobar
func (p ProviderJSON) Set(item Item) (Item, error) {
	filename := p.getFilename(item.Key)

	jsonBytes, _ := json.MarshalIndent(item, "", "    ")

	if err := ioutil.WriteFile(filename, jsonBytes, 0644); err != nil {
		return Item{}, err
	}

	return item, nil
}

// Delete is foobar
func (p ProviderJSON) Delete(key string) error {
	filename := p.getFilename(key)

	has, _ := p.Has(key)

	if !has {
		return nil
	}

	return os.Remove(filename)
}

// Has is foobar
func (p ProviderJSON) Has(key string) (bool, error) {
	filename := p.getFilename(key)

	if _, err := os.Stat(filename); err != nil {
		return false, nil
	}

	return true, nil
}

func (p ProviderJSON) getFilename(key string) string {
	if p.Path == "" {
		p.Path = "."
	}

	p.Path = strings.TrimSuffix(p.Path, "/")

	if p.HashFilenames {
		h := sha256.New()
		if _, err := h.Write([]byte(key)); err == nil {
			key = hex.EncodeToString(h.Sum(nil))
		}
	}

	return fmt.Sprintf("%s/%s.json", p.Path, key)
}
