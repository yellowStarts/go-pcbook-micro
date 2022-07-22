package service

import "sync"

type UserStore interface {
	// 保存用户
	Save(user *User) error
	// 通过姓名查找用户
	Find(username string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.Username] != nil {
		return ErrAlreadyExits
	}

	store.users[user.Username] = user.Clone()
	return nil
}

func (store *InMemoryUserStore) Find(username string) (*User, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	user := store.users[username]
	if user == nil {
		return nil, nil
	}
	return user.Clone(), nil
}
