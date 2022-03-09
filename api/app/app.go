package app

import (
	"errors"
	"time"

	"github.com/kjk/betterguid"
	"github.com/pbivrell/lights/api/handlers"
	"github.com/pbivrell/lights/api/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var ErrorInvalidLight = errors.New("invalid light")
var ErrorUnauthenticated = errors.New("unauthenticed")

type App struct {
	store     storage.Storage
	Responder handlers.Responder
	Dlogger   *logrus.Entry
}

func New(store storage.Storage, responder handlers.Responder, logger *logrus.Entry) *App {
	return &App{
		store:     store,
		Responder: responder,
		Dlogger:   logger,
	}
}

func (a *App) SetLights(light *storage.Light) error {

	l, _ := a.store.ReadLight(light.ID)

	light.Updated = time.Now()

	// If light was found ID will not be empty. If ID is not the uploading light ID
	// or the key does not match this request is invalid
	if l.ID != "" && (light.ID != l.ID || l.Key != light.Key) {
		return ErrorInvalidLight
	}

	return a.store.WriteLight(light)
}

func (a *App) GetLights(id string) (*storage.Light, error) {

	light, err := a.store.ReadLight(id)
	if err != nil {
		return light, err
	}

	light.Key = ""

	return light, nil
}

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (a *App) Login(c *Credentials) (string, error) {

	dlogger := a.Dlogger.WithFields(logrus.Fields{
		"app.Method": "login",
	})

	dlogger.Debug("start")
	defer dlogger.Debug("complete")

	user, err := a.store.ReadUser(c.User)
	if err != nil {
		return "", err
	}

	dlogger.Debug("read user from store")

	if !checkPasswordHash(c.Password, user.Hash) {
		return "", ErrorUnauthenticated
	}

	dlogger.Debug("matched password")

	session := betterguid.New()

	err = a.store.WriteSession(&storage.Session{
		ID:      session,
		UserID:  user.ID,
		Created: time.Now(),
	})
	if err != nil {
		return "", err
	}

	dlogger.Debug("wrote session")

	return session, nil
}

func (a *App) Logout(session string) error {
	return a.store.DeleteSession(session)
}

var ErrorInvalidUser = errors.New("invalid user")

func (a *App) Register(user *storage.User, session string) error {

	if user.ID == "" {
		return ErrorInvalidUser
	}

	_, err := a.store.ReadUser(user.ID)

	// We failed to look up existing user because of an storage issues
	// we can not create this user now
	if !errors.Is(err, storage.ErrorNotFound) {
		return err
	}

	// If we are attempt to register a new user and we succeeded reading
	// a user from the storage that means this user id is taken so we must error
	if session == "" && err == nil {
		return ErrorInvalidUser
	}

	// If someone is attempting to update an existing user we need to
	// verify the session token is valid
	if session != "" {
		sess, err := a.store.ReadSession(session)
		if err != nil {
			return err
		}

		if sess.UserID != user.ID {
			return ErrorUnauthenticated
		}

	}

	// Rehash password if one is provided
	if user.Password != "" {
		user.Hash, err = hashPassword(user.Password)
		if err != nil {
			return err
		}
	}

	user.Password = ""

	return a.store.WriteUser(user)
}

func (a *App) GetUser(session string, hidePrivate bool) (*storage.User, error) {

	sess, err := a.store.ReadSession(session)
	if err != nil {
		return nil, err
	}

	user, err := a.store.ReadUser(sess.UserID)

	if hidePrivate {
		user.Password = ""
		user.Hash = ""
	}
	return user, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
