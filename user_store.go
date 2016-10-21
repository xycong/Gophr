package main

import (

        "encoding/json"
    "fmt"
        "io/ioutil"
    "os"
    "strings"

)

var globalUserStore UserStore

func init() {
    store, err := NewFileUserStore("./data/users.json")
    if err != nil {
        panic(fmt.Errorf("Error creating user store: %s", err))
    }
    globalUserStore = store
}

type UserStore interface {
    Find(string) (*User, Error)
    FindByEmail(string) (*User, Error)
    FindByUsername(string) (*User, Error)
    Save(User) error
}

type FileUserStore struct {
    filename string
    Users map[string]User
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
    store := &FileUserStore{
        Users: map[string]User{},
        filename: filename,
    }
    
    contents, err:= ioutil.ReadFile(filename)

    if err != nil {
        if os.IsNotExist(err) {
            return store, nil
        }
        return nil, err
    }
    err = json.Unmarshal(contents, store) 
    if err != nil {
        return nil, err
    }
    return store, nil
}

func (store FileUserStore) Save(user User) error {
    store.Users[user.ID] = user

    contents,err := json.MarshalIndent(store, "", "  ")
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(store.filename, contents, 0660)
    if err != nil {
        return err
    }
    return nil
}

func (store FileUserStore) Find(id string) (*User, error) {
    user, ok := store.Users[id]
    if ok {
        return &user, nil
    }
    return nil, nil
}

func (store FileUserStore) FindByUsername(username string) (*User, error) {
    if username == "" {
        return nil, nil
    }

    for _, user := range store.Users {
        if strings.ToLower(username) == strings.ToLower(user.Username) {
            return &user, nil
        }
    }
    return nil, nil
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
    if email == "" {
        return nil, nil
    }

    for _, user:= range store.Users {
        if strings.ToLower(email) == strings.ToLower(user.Email) {
            return &user, nil
        }
    }
    return nil, nil
}