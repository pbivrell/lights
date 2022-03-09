package storage

import (
	"errors"
	"time"
)

type Light struct {
	ID      string    `dynamodbav:"id"`
	Mac     string    `dynamodbav:"mac"`
	IP      string    `dynamodbav:"ip"`
	Key     string    `dynamodbav:"key"`
	Version string    `dynamodbav:"version"`
	Status  bool      `dynamodbav:"status"`
	Updated time.Time `dynamodbav:"updated"`
	Count   int       `dynamodbav:"count"`
	Running string    `dynamodbav:"running"`
}

type User struct {
	ID       string       `json:"id" dynamodbav:"id"`
	Hash     string       `json:"hash" dynamodbav:"hash"`
	Email    string       `json:"email" dynamodbav:"email"`
	Lights   []LightAlias `json:"lights" dynamodbav:"lights"`
	Groups   []Group      `json:"groups" dynamodbav:"groups"`
	Password string       `json:"password"`
}

type LightAlias struct {
	ID    string `json:"id" dynamodbav:"id"`
	Alias string `json:"alias" dynamodbav:"alias"`
}

type Group struct {
	Name   string   `dynamodbav:"name"`
	Lights []string `dynamodbav:"lights"`
}

type Session struct {
	ID      string    `json:"id" dynamodbav:"id"`
	UserID  string    `json:"uid" dynamodbav:"uid"`
	Created time.Time `json:"created" dynamodbav:"created"`
}

var ErrorNotFound = errors.New("no item with id")

type Storage interface {
	WriteSession(*Session) error
	DeleteSession(key string) error
	ReadSession(key string) (*Session, error)
	WriteUser(u *User) error
	ReadUser(key string) (*User, error)
	DeleteUser(key string) error
	WriteLight(l *Light) error
	ReadLight(key string) (*Light, error)
}
