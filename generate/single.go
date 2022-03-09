package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	l := lights.New(2)

	l.Print("./single.bin")

}
