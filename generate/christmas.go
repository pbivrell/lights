package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	count := 100

	l := lights.New(uint16(count))

	for i := 2; i < count; i += 3 {
		l.SetColor(uint16(i-1), 209, 5, 5)
		l.SetColor(uint16(i-2), 215, 215, 215)
		l.SetColor(uint16(i), 53, 203, 0)
	}
	l.SetDelay(500)
	for i := 2; i < count; i += 3 {
		l.SetColor(uint16(i), 209, 5, 5)
		l.SetColor(uint16(i-1), 215, 215, 215)
		l.SetColor(uint16(i-2), 53, 203, 0)
	}
	l.SetDelay(500)
	for i := 2; i < count; i += 3 {
		l.SetColor(uint16(i-2), 209, 5, 5)
		l.SetColor(uint16(i), 215, 215, 215)
		l.SetColor(uint16(i-1), 53, 203, 0)
	}
	l.SetDelay(500)
	l.Print()
}
