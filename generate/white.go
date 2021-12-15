package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	count := 150

	l := lights.New(uint16(count))

	for i := 0; i < count; i++ {
		l.SetColor(uint16(i), 120, 120, 30)
	}
	l.SetDelay(10000)
	l.Print("white.bin")
}
