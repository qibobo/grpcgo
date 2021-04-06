package store

import (
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/qibobo/grpcgo/models"
)

type Store interface {
	Get(id string) interface{}
	Save(o interface{}) (string, error)
	List() []interface{}
}

type InMemoryStore struct {
	store map[string]interface{}
	lock  sync.RWMutex
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		store: make(map[string]interface{}, 1024),
		lock:  sync.RWMutex{},
	}
}

func (ims *InMemoryStore) Get(id string) interface{} {
	ims.lock.RLock()
	defer ims.lock.RUnlock()
	return ims.store[id]
}

func (ims *InMemoryStore) Save(o interface{}) (string, error) {
	log.Printf("saving obj %v\n", o)
	p, ok := o.(*models.Person)
	if !ok {
		return "", errors.New("object is not a person")
	}
	ims.lock.Lock()
	defer ims.lock.Unlock()
	p.Id = uuid.New().String()
	ims.store[p.Id] = p
	return p.Id, nil

}
func (ims *InMemoryStore) List() []interface{} {
	result := []interface{}{}
	for _, p := range ims.store {
		result = append(result, p)
	}
	return result
}
