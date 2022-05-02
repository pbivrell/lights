package server

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gocarina/gocsv"
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

	r.HandleFunc("/user/hub/{hub}", getHub(a)).Methods("GET")
	//r.HandleFunc("/user/hub", getHubs(a)).Methods("POST")

	r.HandleFunc("/user/hub/{hub}/light/{light}", getLight(a)).Methods("GET")
	r.HandleFunc("/user/hub/{hub}/light/{light}", setLight(a)).Methods("POST")

	r.HandleFunc("/pattern", listPatterns(a)).Methods("GET")
	r.HandleFunc("/pattern", setPattern(a)).Methods("POST")
	r.HandleFunc("/pattern/{id}", getPattern(a)).Methods("GET")

	r.HandleFunc("/hub", hubReport(a)).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../light/build")))

	return r
}

func getHub(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "getHub",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		hubID := vars["hub"]

		hub, err := a.GetHub(hubID)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		session := parseCookie(r)

		hasHub, err := a.HasHub(session, hubID)
		if err != nil && !errors.Is(err, storage.ErrorNotFound) {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		if err != nil || !hasHub {
			respondError(w, a.Responder.Unauthorized())
			return

		}

		dlogger.Debug("retrieved hub")

		err = json.NewEncoder(w).Encode(&hub)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("encoded json")

		respondOK(w, a.Responder.OK())
	}
}

func listPatterns(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "listPatterns",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		patterns, err := a.GetPatterns()
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		err = json.NewEncoder(w).Encode(patterns)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		respondOK(w, a.Responder.OK())
	}

}

func getPattern(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "getPattern",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		patternID := vars["id"]

		pattern, err := a.GetPattern(patternID)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		err = binary.Write(w, binary.BigEndian, pattern.Data)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		respondOK(w, a.Responder.OK())
	}
}

func setPattern(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "setPattern",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		patterns, err := a.GetPatterns()
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		patternID := len(patterns)

		var alias string
		aliasQP, ok := r.URL.Query()["alias"]
		if ok {
			alias = aliasQP[0]
		}

		session := parseCookie(r)

		userID, err := a.LookupUser(session)
		if errors.Is(err, storage.ErrorNotFound) {
			respondError(w, a.Responder.Unauthorized())
			return
		}

		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		var pattern storage.Pattern

		pattern.Creator = userID
		pattern.Alias = alias
		pattern.ID = fmt.Sprintf("%d", patternID)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		pattern.Data = data

		err = a.SetPattern(&pattern)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		resp := struct {
			ID string `json:"id"`
		}{
			ID: pattern.ID,
		}

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		respondOK(w, a.Responder.OK())
	}
}

func setLight(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "getLight",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		hubID := vars["hub"]
		lightID := vars["light"]

		light, err := a.GetLight(lightID)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		dlogger.Debugf("found light: %v", light)

		session := parseCookie(r)

		hasLight, err := a.HasLight(session, hubID, lightID)
		if err != nil && !errors.Is(err, storage.ErrorNotFound) {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		if err != nil || !hasLight {
			respondError(w, a.Responder.Unauthorized())
			return

		}

		dlogger.Debug("retrieved light")

		var updatedLight storage.LightNullable

		err = json.NewDecoder(r.Body).Decode(&updatedLight)
		if err != nil {
			respondError(w, a.Responder.BadRequest())
			return
		}

		err = a.SetLight(light, &updatedLight)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		respondOK(w, a.Responder.OK())
	}
}

func getLight(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "getLight",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		hubID := vars["hub"]
		lightID := vars["light"]

		light, err := a.GetLight(lightID)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		session := parseCookie(r)

		hasLight, err := a.HasLight(session, hubID, lightID)
		if err != nil && !errors.Is(err, storage.ErrorNotFound) {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		if err != nil || !hasLight {
			respondError(w, a.Responder.Unauthorized())
			return

		}

		dlogger.Debug("retrieved light")

		err = json.NewEncoder(w).Encode(&light)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("encoded json")

		respondOK(w, a.Responder.OK())
	}
}

func hubReport(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var hub storage.Hub

		dlogger := a.Dlogger.WithFields(logrus.Fields{
			"handlerFunc": "hubReport",
		})

		dlogger.Debug("start")
		defer dlogger.Debug("complete")

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(data, &hub)
		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		dlogger.Debug("decoded json")

		changes, err := a.ReportHub(&hub)
		if errors.Is(err, app.ErrorInvalidHub) {
			respondError(w, a.Responder.Unauthorized())
			return
		}

		if err != nil {
			respondError(w, a.Responder.InternalError(err))
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		gocsv.MarshalWithoutHeaders(changes, w)

		respondOK(w, a.Responder.OK())

	}
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

	for k, v := range r.Headers {
		w.Header().Set(k, v)

	}
}

func respondError(w http.ResponseWriter, r handlers.ErrorResponse) {

	body, _ := json.Marshal(r)

	http.Error(w, string(body), r.Code)
}
