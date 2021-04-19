package store

import (
	"sync"

	"github.com/qibobo/grpcgo/models"
)

type UserStore interface {
	Save(user *models.User) error
	Find(username string) *models.User
}

type InMemoryUserStore struct {
	userMap map[string]*models.User
	lock    sync.RWMutex
}

func NewInMemoryUserStore() UserStore {
	return &InMemoryUserStore{
		userMap: make(map[string]*models.User),
		lock:    sync.RWMutex{},
	}
}

func (us *InMemoryUserStore) Save(user *models.User) error {
	us.lock.Lock()
	defer us.lock.Unlock()
	us.userMap[user.UserName] = user
	return nil
}
func (us *InMemoryUserStore) Find(username string) *models.User {
	us.lock.RLock()
	defer us.lock.RUnlock()
	if user, ok := us.userMap[username]; ok && user != nil {
		return user
	}
	return nil
}
