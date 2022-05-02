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
	store     storage.AppStorage
	Responder handlers.Responder
	Dlogger   *logrus.Entry
}

func New(store storage.AppStorage, responder handlers.Responder, logger *logrus.Entry) *App {
	return &App{
		store:     store,
		Responder: responder,
		Dlogger:   logger,
	}
}

func (a *App) ListPatterns() ([]*storage.Pattern, error) {
	return []*storage.Pattern{}, nil
}

func (a *App) GetPatterns() ([]*storage.Pattern, error) {
	return a.store.Pattern.List()
}

func (a *App) GetPattern(patternID string) (*storage.Pattern, error) {
	return a.store.Pattern.Read(patternID)
}

func (a *App) SetPattern(pattern *storage.Pattern) error {

	oldPattern, err := a.store.Pattern.Read(pattern.ID)
	if err != nil {
		return err
	}

	if oldPattern.Alias != pattern.Alias && pattern.Alias != "" {
		pattern.Alias = oldPattern.Alias
	}

	return a.store.Pattern.Write(pattern.ID, pattern)
}

func (a *App) SetLight(light *storage.Light, update *storage.LightNullable) error {

	if update.Status != nil && light.Status != *update.Status {
		light.Status = *update.Status
		light.Changed = time.Now()

	}

	if update.PatternID != nil && light.PatternID != *update.PatternID {
		light.PatternID = *update.PatternID
		light.Changed = time.Now()
	}

	if update.Count != nil && light.Count != *update.Count {
		light.Count = *update.Count
		light.Changed = time.Now()
	}

	if update.Alias != nil && light.Alias != *update.Alias {
		light.Alias = *update.Alias
	}

	return a.store.Light.Write(light.ID, light)
}

var ErrorInvalidHub = errors.New("invalid hub")

func (a *App) GetLight(light string) (*storage.Light, error) {
	dlogger := a.Dlogger.WithFields(logrus.Fields{
		"app.Method": "getLight",
	})

	dlogger.Debug("start")
	defer dlogger.Debug("complete")

	l, err := a.store.Light.Read(light)
	dlogger.Debug("got light", light, l, err)
	return l, err

}

func (a *App) GetHub(hub string) (*storage.Hub, error) {
	dlogger := a.Dlogger.WithFields(logrus.Fields{
		"app.Method": "getHub",
	})

	dlogger.Debug("start")
	defer dlogger.Debug("complete")

	return a.store.Hub.Read(hub)
}

func (a *App) ReportHub(h *storage.Hub) ([]*storage.Light, error) {

	dlogger := a.Dlogger.WithFields(logrus.Fields{
		"app.Method": "reportHub",
	})

	dlogger.Debug("start")
	defer dlogger.Debug("complete")

	hub, err := a.store.Hub.Read(h.ID)
	if err != nil {
		return []*storage.Light{}, err
	}
	dlogger.Debug("read hub data")

	if hub.Key != h.Key {
		return []*storage.Light{}, ErrorInvalidHub
	}
	dlogger.Debug("valid hub key")

	updatedHubLights := make([]string, len(hub.Lights))

	existingLights := make(map[string]struct{})
	for i, light := range hub.Lights {
		existingLights[light] = struct{}{}
		updatedHubLights[i] = light
	}

	data := make([]*storage.Light, len(h.Lights))

	for i, light := range h.Lights {

		var updatedLight *storage.Light

		_, ok := existingLights[light]
		if !ok {
			updatedLight = &storage.Light{
				Alias:     "",
				ID:        light,
				Status:    false,
				PatternID: "",
				Count:     50,
			}
			updatedHubLights = append(updatedHubLights, light)

		} else {
			updatedLight, err = a.store.Light.Read(light)
		}

		data[i] = updatedLight
		updatedLight.Updated = time.Now()

		// The client of this API endpoint has aggresive retry logic if we fail to write then continue on
		// We'd only really want to be notified if this was failing continously
		err = a.store.Light.Write(light, updatedLight)
		if err != nil {
			dlogger.Debugf("light write fail: %v", err)
		}

		h.Lights = updatedHubLights
		err = a.store.Hub.Write(h.ID, h)
		if err != nil {
			dlogger.Debugf("light write fail: %v", err)
		}
	}
	return data, nil
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

	user, err := a.store.User.Read(c.User)
	if err != nil {
		return "", err
	}

	dlogger.Debug("read user from store")

	if !checkPasswordHash(c.Password, user.Hash) {
		return "", ErrorUnauthenticated
	}

	dlogger.Debug("matched password")

	session := betterguid.New()

	err = a.store.Session.Write(session, &storage.Session{
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

func (a *App) LookupUser(session string) (string, error) {
	sess, err := a.store.Session.Read(session)
	if err != nil {
		return "", err
	}
	return sess.UserID, nil
}

func (a *App) Logout(session string) error {
	return a.store.Session.Delete(session)
}

var ErrorInvalidUser = errors.New("invalid user")

func (a *App) Register(user *storage.User, session string) error {

	dlogger := a.Dlogger.WithFields(logrus.Fields{
		"app.Method": "register",
	})

	dlogger.Debug("start")
	defer dlogger.Debug("complete")

	if user.ID == "" {
		return ErrorInvalidUser
	}

	dlogger.Debug("user ID not empty")

	_, err := a.store.User.Read(user.ID)

	// We failed to look up existing user because of an storage issues
	// we can not create this user now
	if !errors.Is(err, storage.ErrorNotFound) {
		return err
	}

	dlogger.Debug("looked up user")

	// If we are attempt to register a new user and we succeeded reading
	// a user from the storage that means this user id is taken so we must error
	if session == "" && err == nil {
		return ErrorInvalidUser
	}

	dlogger.Debug("looked up user")

	// If someone is attempting to update an existing user we need to
	// verify the session token is valid
	if session != "" {
		sess, err := a.store.Session.Read(session)
		if err != nil {
			return err
		}

		if sess.UserID != user.ID {
			return ErrorUnauthenticated
		}

		dlogger.Debug("can update user")

	}

	// Rehash password if one is provided
	if user.Password != "" {
		user.Hash, err = HashPassword(user.Password)
		if err != nil {
			return err
		}
		dlogger.Debug("rehashed password")
	}

	user.Password = ""

	return a.store.User.Write(user.ID, user)
}

func (a *App) HasHub(session, hubName string) (bool, error) {
	dlogger := a.Dlogger.WithFields(logrus.Fields{
		"app.Method": "hasHub",
	})

	dlogger.Debug("start")
	defer dlogger.Debug("complete")

	dlogger.Debug("apple" + session)
	user, err := a.GetUser(session, true)
	if err != nil {
		return false, err
	}
	dlogger.Debug("found user by session id")

	for _, h := range user.Hubs {
		if h.ID == hubName {
			return true, nil
		}
	}
	dlogger.Debug("no matching hubName")
	return false, nil
}

func (a *App) HasLight(session, hubName, lightName string) (bool, error) {

	hasHub, err := a.HasHub(session, hubName)
	if err != nil {
		return false, err
	}

	if !hasHub {
		return false, nil
	}

	hub, err := a.GetHub(hubName)
	if err != nil {
		return false, err
	}

	for _, light := range hub.Lights {
		if light == lightName {
			return true, nil
		}
	}
	return false, nil

}

func (a *App) GetUser(session string, hidePrivate bool) (*storage.User, error) {

	sess, err := a.store.Session.Read(session)
	if err != nil {
		return nil, err
	}

	user, err := a.store.User.Read(sess.UserID)

	if hidePrivate {
		user.Password = ""
		user.Hash = ""
	}
	return user, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
