package main

import (
	"errors"
	"sync"
	"unicode"

	sql "github.com/FloatTech/sqlite"
)

type user struct {
	Name     string `db:"name"`
	Password string `db:"pwd"`
}

type usersdb struct {
	sync.RWMutex
	sql.Sqlite
}

func NewUsersDB(file string) (*usersdb, error) {
	db := &usersdb{}
	db.DBPath = file
	err := db.Open()
	if err != nil {
		return nil, err
	}
	err = db.Create("default", &user{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *usersdb) Add(location, name, password string) error {
	db.Lock()
	defer db.Unlock()
	return db.Insert(location, &user{Name: name, Password: password})
}

func (db *usersdb) Password(location, name string) (pwd string, err error) {
	for _, c := range name {
		if !unicode.IsLetter(c) {
			err = errors.New("invaild user name")
			return
		}
	}
	u := &user{}
	db.RLock()
	defer db.RUnlock()
	err = db.Find(location, &u, "WHERE name='"+name+"'")
	pwd = u.Password
	return
}
