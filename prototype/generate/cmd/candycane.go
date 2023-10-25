package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	count := 150

	l := lights.New(uint16(count))

	for i := 1; i < count; i += 2 {
		l.SetColor(uint16(i-1), 209, 5, 5)
		l.SetColor(uint16(i), 120, 120, 40)
	}
	l.SetDelay(1000)
	for i := 1; i < count; i += 2 {
		l.SetColor(uint16(i), 209, 5, 5)
		l.SetColor(uint16(i-1), 120, 120, 40)
	}
	l.SetDelay(1000)
	l.Print("candycane.bin")
}
