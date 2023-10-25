package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	URL    = "https://lights.paulbivrell.com"
	Creds  = "demo"
	Hub    = "Hub1"
	Key    = "theoriginallightskey"
	Light1 = "light1"
)

func main() {
	fmt.Println("login")
	session := login(Creds, Creds)
	fmt.Println("get hub")
	GetHub(session, Hub)
	fmt.Println("report hub")
	ReportHub(session, []string{Light1})
	fmt.Println("get light")
	GetLight(session, Hub, Light1)
	fmt.Println("get hub")
	GetHub(session, Hub)
	fmt.Println("toggle light")
	ToggleLight(session, Hub, Light1, true)
	//SetLight(session, Hub, Light1, "first.bin")
	fmt.Println("get light")
	GetLight(session, Hub, Light1)
	fmt.Println("report hub")
	ReportHub(session, []string{Light1})
	fmt.Println("get hub")
	GetHub(session, Hub)
	fmt.Println("insert pattern")
	InsertPattern(session, "rainbow.bin", "rainbow")
	fmt.Println("list patterns")
	ListPattern()

}

func login(user, pass string) string {

	creds := struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}{
		User:     user,
		Password: pass,
	}

	body := &bytes.Buffer{}

	err := json.NewEncoder(body).Encode(&creds)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(URL+"/user", "application/json", body)
	if err != nil {
		panic(err)
	}

	cookies, ok := resp.Header["Set-Cookie"]
	if !ok {
		panic("no cookie header")
	}

	return strings.SplitN(strings.SplitN(cookies[0], ";", 2)[0], "=", 2)[1]
}

func ToggleLight(session, hub, light string, state bool) {

	lightState := "off"
	if state {
		lightState = "on"
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/hub/%s/light/%s?s=%s", URL, hub, light, lightState), nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

}

func SetLight(session, hub, light, pattern string) {

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/hub/%s/light/%s?p=%s", URL, hub, light, pattern), nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func GetLight(session, hub, light string) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/hub/%s/light/%s", URL, hub, light), nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

}

func GetHub(session, hub string) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/hub/%s", URL, hub), nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

}

func ReportHub(session string, lights []string) {

	type HubData struct {
		Key    string   `json:"hubKey" dynamodbav:"key"`
		ID     string   `json:"hubID" dynamodbav:"id"`
		Lights []string `json:"lights" dynamodbav:"lights"`
	}

	body, err := json.Marshal(&HubData{
		Key:    Key,
		ID:     Hub,
		Lights: lights,
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hub", URL), strings.NewReader(string(body)))
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func InsertPattern(session, filePath, alias string) string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pattern?alias=%s", URL, alias), f)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	return string(data)

}

func ListPattern() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/pattern", URL), nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
