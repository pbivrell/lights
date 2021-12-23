package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	count := 100

	l := lights.New(uint16(count))

	for i := 0; i < count; i += 6 {
		l.SetColor(uint16(i), 85, 92, 1)    // yellow
		l.SetColor(uint16(i+1), 15, 83, 1)  // green
		l.SetColor(uint16(i+2), 3, 2, 49)   // blue
		l.SetColor(uint16(i+3), 49, 2, 6)   // pink
		l.SetColor(uint16(i+4), 154, 9, 5)  // re
		l.SetColor(uint16(i+5), 154, 36, 5) // orange
	}
	l.SetDelay(30000)
	l.Print("rainbow.bin")
}
