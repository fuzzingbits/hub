package cache

import "sync"

// ProviderMemory is foobar
type ProviderMemory struct {
	store map[string]Item
	mutex *sync.Mutex
}

// Get is foobar
func (p *ProviderMemory) Get(key string) (Item, error) {
	p.setup()
	has, _ := p.Has(key)

	if !has {
		return Item{}, ErrorCacheKeyNotFound
	}

	return p.store[key], nil
}

// Set is foobar
func (p *ProviderMemory) Set(item Item) (Item, error) {
	p.setup()
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.store[item.Key] = item

	return item, nil
}

// Delete is foobar
func (p *ProviderMemory) Delete(key string) error {
	p.setup()
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.store, key)

	return nil
}

// Has is foobar
func (p *ProviderMemory) Has(key string) (bool, error) {
	p.setup()
	_, ok := p.store[key]

	return ok, nil
}

// Has is foobar
func (p *ProviderMemory) setup() {
	if p.store == nil {
		p.store = make(map[string]Item)
	}

	if p.mutex == nil {
		p.mutex = &sync.Mutex{}
	}
}
