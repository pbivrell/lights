package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	count := 200

	l := lights.New(uint16(count))

	for i := 0; i < count; i += 5 {
		l.SetColor(uint16(i), 255, 32, 45)
		l.SetColor(uint16(i+1), 255, 239, 31)
		l.SetColor(uint16(i+2), 50, 255, 32)
		l.SetColor(uint16(i+3), 32, 34, 255)
		l.SetColor(uint16(i+4), 255, 32, 219)
	}
	l.SetDelay(30000)
	l.Print("rainbow.bin")
}
