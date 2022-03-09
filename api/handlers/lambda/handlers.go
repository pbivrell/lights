package lambda

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pbivrell/lights/api/app"
	"github.com/pbivrell/lights/api/storage"
	"github.com/pbivrell/lights/api/version"
)

var AllowedOrigins = []string{"http://polylights-react-app.s3-website-us-east-1.amazonaws.com"}

//var AllowedOrigins = []string{"http://localhost:3000"}

func GetRouter(a *app.App) func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		if req.HTTPMethod == "OPTIONS" {
			return corsResponse(), nil
		}

		switch path.Base(req.Path) {

		case "version":
			switch req.HTTPMethod {
			case "GET":
				return versionResponse(), nil
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Body:       http.StatusText(http.StatusMethodNotAllowed),
				}, nil

			}
		case "light":
			switch req.HTTPMethod {
			case "GET":
				return getLight(a, req)
			case "POST":
				return setLight(a, req)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Body:       http.StatusText(http.StatusMethodNotAllowed),
				}, nil
			}

		case "user":
			switch req.HTTPMethod {
			case "GET":
				return getUser(a, req)
			case "POST":
				return login(a, req)
			case "PUT":
				return register(a, req)
			case "DELETE":
				return logout(a, req)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Body:       http.StatusText(http.StatusMethodNotAllowed),
				}, nil
			}
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       http.StatusText(http.StatusNotFound),
			}, nil

		}
	}
}

func corsResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(http.StatusOK),
		Headers: map[string]string{
			"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
			"Access-Control-Allow-Methods":     "POST,GET,PATCH,PUT,DELETE",
			"Access-Control-Allow-Credentials": "true",
		},
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
	}
}

func versionResponse() events.APIGatewayProxyResponse {

	v, _ := json.Marshal(version.Version)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(v),
	}

}

func getLight(a *app.App, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	session := parseCookie(req)

	if session == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(http.StatusBadRequest),
			Headers: map[string]string{
				"Access-Control-Allow-Credentials": "true",
			},
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	user, err := a.GetUser(session, true)

	if errors.Is(err, storage.ErrorNotFound) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       fmt.Sprintf("%v", err),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	id := req.QueryStringParameters["light"]

	ownsLight := false

	for _, light := range user.Lights {
		if id == light.ID {
			ownsLight = true
			break
		}
	}

	if !ownsLight {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "you don't own that light",
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	errEvent := events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
	}

	light, err := a.GetLights(id)
	if err != nil {
		return errEvent, nil
	}

	jsonString, err := json.Marshal(light)
	if err != nil {
		return errEvent, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonString),
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func setLight(a *app.App, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       http.StatusText(http.StatusBadRequest),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	errEvent := events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
	}

	var light storage.Light

	err := json.Unmarshal([]byte(req.Body), &light)
	if err != nil {
		errEvent.Body = fmt.Sprintf("%v", err)
		return errEvent, nil
	}

	err = a.SetLights(&light)
	if err != nil {
		errEvent.Body = fmt.Sprintf("%v", err)
		return errEvent, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func login(a *app.App, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	errEvent := events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
	}

	var creds app.Credentials

	err := json.Unmarshal([]byte(req.Body), &creds)
	if err != nil {
		errEvent.Body = fmt.Sprintf("%v", err)
		return errEvent, nil
	}

	session, err := a.Login(&creds)
	if err != nil {
		errEvent.Body = fmt.Sprintf("%v", err)
		return errEvent, nil
	}

	expires := time.Now().Add(10 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(http.StatusOK),
		Headers: map[string]string{
			"Set-Cookie":                       fmt.Sprintf("session=%s; SameSite=None; Secure; Expires=%s", session, expires),
			"Access-Control-Allow-Credentials": "true",
		},
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
	}, nil
}

func processCookie(cookie string) string {

	cookies := strings.Split(cookie, ";")
	for _, v := range cookies {
		temp := strings.SplitN(v, "=", 2)
		if temp[0] == "session" {
			return temp[1]
		}
	}
	return ""
}

func parseCookie(req events.APIGatewayProxyRequest) string {

	const cookieHeader = "cookie"

	cookie, _ := req.Headers[cookieHeader]

	session := processCookie(cookie)

	if session != "" {
		return session
	}

	cookies, _ := req.MultiValueHeaders[cookieHeader]
	for _, cookie := range cookies {
		session = processCookie(cookie)
		if session != "" {
			return session
		}
	}

	return ""
}

func logout(a *app.App, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	expires := time.Now().Add(-10 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	session := parseCookie(req)
	if session == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	err := a.Logout(session)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Set-Cookie":                       fmt.Sprintf("session=%s; SameSite=None; Secure; Expires=%s", session, expires),
		},
	}, nil

}

func register(a *app.App, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user storage.User

	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	session := parseCookie(req)

	err = a.Register(&user, session)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func getUser(a *app.App, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	reqJSON, _ := json.Marshal(req)
	session := parseCookie(req)

	if session == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(reqJSON),
			Headers: map[string]string{
				"Access-Control-Allow-Credentials": "true",
			},
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	user, err := a.GetUser(session, true)

	if errors.Is(err, storage.ErrorNotFound) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       fmt.Sprintf("%v", err),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	jsonString, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(http.StatusInternalServerError),
			MultiValueHeaders: map[string][]string{
				"Access-Control-Allow-Origin": AllowedOrigins,
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonString),
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
		},
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Origin": AllowedOrigins,
		},
	}, nil
}
