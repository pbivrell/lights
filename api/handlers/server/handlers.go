package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pbivrell/lights/api/app"
	"github.com/pbivrell/lights/api/handlers"
	"github.com/pbivrell/lights/api/storage"
	"github.com/sirupsen/logrus"
)

func RegisterHandlers(a *app.App) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/user", getUser(a)).Methods("GET")
	r.HandleFunc("/user", login(a)).Methods("POST")
	r.HandleFunc("/user", register(a)).Methods("PUT")
	r.HandleFunc("/user", logout(a)).Methods("DELETE")

	return r
}

func login(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var creds app.Credentials

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "login",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		defer r.Body.Close()

		dlogger.Debug("decoded json")

		session, err := a.Login(&creds)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		dlogger.Debug("logged in")

		expires := time.Now().Add(10 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

		respondOK(w, a.Responder.OK(map[string]string{
			"Set-Cookie": fmt.Sprintf("session=%s; SameSite=None; Secure; Expires=%s", session, expires),
		}))
	}
}

func parseCookie(req *http.Request) string {

	cookie, _ := req.Cookie("session")
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

func logout(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "logout",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		session := parseCookie(r)
		err := a.Logout(session)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("logged out")

		expires := time.Now().Add(-10 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

		respondOK(w, a.Responder.OK(map[string]string{
			"Set-Cookie": fmt.Sprintf("session=%s; SameSite=None; Secure; Expires=%s", session, expires),
		}))
	}
}

func register(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "register",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		var user storage.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		defer r.Body.Close()

		dlogger.Debug("decoded json")

		session := parseCookie(r)

		err = a.Register(&user, session)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("registered")

		respondOK(w, a.Responder.OK())
	}
}

func getUser(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "register",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		session := parseCookie(r)

		user, err := a.GetUser(session, true)
		if errors.Is(err, storage.ErrorNotFound) {
			respondError(w, a.Responder.Unauthorized())
			return
		}
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("got user")

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("encoded json")

		respondOK(w, a.Responder.OK())
	}
}

func respondOK(w http.ResponseWriter, r handlers.Response) {

	w.WriteHeader(r.Code)
}

func respondError(w http.ResponseWriter, r handlers.ErrorResponse) {

	body, _ := json.Marshal(r)

	http.Error(w, string(body), r.Code)
}
