package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

var globalUserStore UserStore

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type FileUserStore struct {
	filename string
	Users    map[string]User
}

func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user
	store.Users[user.Email] = user
	store.Users[user.Username] = user
	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store FileUserStore) Find(id string) (*User, error) {
	return store.getUserByKey(id)
}

func (store FileUserStore) FindByUsername(username string) (*User, error) {
	return store.getUserByKey(username)
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
	return store.getUserByKey(email)
}

func (store FileUserStore) getUserByKey(key string) (*User, error) {
	if key == "" {
		return nil, nil
	}

	user, ok := store.Users[key]
	if ok {
		return &user, nil
	}
	return nil, nil
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		Users:    map[string]User{},
		filename: filename,
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}

	for _, user := range store.Users {
		store.Users[strings.ToLower(user.Email)] = user
		store.Users[strings.ToLower(user.Username)] = user
	}

	return store, nil
}
