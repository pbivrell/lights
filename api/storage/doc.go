package storage

import (
	"errors"
	"time"
)

type Hub struct {
	Key    string   `json:"hubKey" dynamodbav:"key"`
	ID     string   `json:"hubID" dynamodbav:"id"`
	Lights []string `json:"lights" dynamodbav:"lights"`
}

type LightNullable struct {
	Alias     *string    `dynamodbav:"alias" csv:"-" json:",omitempty"`
	ID        *string    `dynamodbav:"id" csv:"id" json:",omitempty"`
	Status    *bool      `dynamodbav:"status" csv:"status" json:",omitempty"`
	Updated   *time.Time `dynamodbav:"updated" csv:"-" json:",omitempty"`
	PatternID *string    `dynamodbav:"running" csv:"patternID" json:",omitempty"`
	Count     *uint8     `dynamodbav:"count" csv:"count" json:",omitempty"`
	Changed   *time.Time `dynamodbav:"changed" csv:"-" json:",omitempty"`
}

type Light struct {
	Alias     string    `dynamodbav:"alias" csv:"-"`
	ID        string    `dynamodbav:"id" csv:"id"`
	Status    bool      `dynamodbav:"status" csv:"status"`
	Updated   time.Time `dynamodbav:"updated" csv:"-"`
	PatternID string    `dynamodbav:"running" csv:"patternID"`
	Count     uint8     `dynamodbav:"count" csv:"count"`
	Changed   time.Time `dynamodbav:"changed" csv:"-"`
}

type Pattern struct {
	Alias   string `json:"alias"`
	ID      string `json:"id"`
	Data    []byte `json:"-"`
	Creator string `json:"creator"`
}

type User struct {
	ID       string     `json:"id" dynamodbav:"id"`
	Hash     string     `json:"hash" dynamodbav:"hash"`
	Email    string     `json:"email" dynamodbav:"email"`
	Hubs     []HubAlias `json:"hubs" dynamodbav:"hubs"`
	Groups   []Group    `json:"groups" dynamodbav:"groups"`
	Password string     `json:"password"`
}

type HubAlias struct {
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

type Storage[T any] interface {
	Write(string, *T) error
	Read(string) (*T, error)
	List() ([]*T, error)
	Delete(string) error
}

type AppStorage struct {
	User    Storage[User]
	Session Storage[Session]
	Hub     Storage[Hub]
	Light   Storage[Light]
	Pattern Storage[Pattern]
}
