package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	l := lights.New(50)

	l.Fill(252, 182, 18, 0, 50)

	l.SetDelay(1000)
	l.Print("./setup.bin")

}
