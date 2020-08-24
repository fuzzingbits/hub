package mockableprovider

import (
	"sync"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// Provider is a mockable provider
type Provider struct {
	Mutex        *sync.Mutex
	CreateError  error
	UpdateError  error
	DeleteError  error
	GetByIDError error
	GetAllError  error
	store        map[string]interface{}
}

// NewProvider builds and returns a valid Provider
func NewProvider() *Provider {
	return &Provider{
		Mutex: &sync.Mutex{},
		store: make(map[string]interface{}),
	}
}

// Delete an item from the store
func (p *Provider) Delete(id string) error {
	if p.DeleteError != nil {
		return p.DeleteError
	}

	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	delete(p.store, id)

	return nil
}

// Update an item in the store by id
func (p *Provider) Update(id string, v interface{}) error {
	if p.UpdateError != nil {
		return p.UpdateError
	}

	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.store[id] = v

	return nil
}

// Create and item in the store
func (p *Provider) Create(id string, v interface{}) error {
	if p.CreateError != nil {
		return p.CreateError
	}

	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.store[id] = v

	return nil
}

// GetByID gets an item by id
func (p *Provider) GetByID(id string) (interface{}, error) {
	if p.GetByIDError != nil {
		return nil, p.GetByIDError
	}

	item, found := p.store[id]
	if !found {
		return nil, entity.ErrRecordNotFound
	}

	return item, nil
}

// GetAll items from the store
func (p *Provider) GetAll() ([]interface{}, error) {
	if p.GetAllError != nil {
		return nil, p.GetAllError
	}

	result := []interface{}{}
	for _, item := range p.store {
		result = append(result, item)
	}

	return result, nil
}
