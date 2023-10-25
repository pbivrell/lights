package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

var ip string = "192.168.86.25"
var count int = 100

func main() {
	random()
}

func rainbow() {
	off()
}

func halloween() {

	off()
	defer off()
	all(255, 69, 0)
	for r := 0; r < 1000; r++ {
		color(0, 0, 0, rand.Intn(1000))
		if rand.Intn(36) == 4 {
			all(255, 69, 0)
		}
	}
}

func random() {

	off()
	defer off()
	for r := 0; r < 1000; r++ {
		color(rand.Intn(254), rand.Intn(254), rand.Intn(254), rand.Intn(100))
	}
}

func off() {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://%s/off", ip))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func color(r, g, b, i int) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://%s/color?r=%d&g=%d&b=%d&i=%d", ip, r, g, b, i))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func all(r, g, b int) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://%s/all?r=%d&g=%d&b=%d", ip, r, g, b))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
