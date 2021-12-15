package main

import "github.com/pbivrell/lights/generate/lights"

func main() {

	l := lights.New(10)

	for i := 0; i < 10; i++ {
		if i%5 == 0 {
			l.SetDelay(1000)
		}
		l.SetColor(uint16(i%5), uint8(128*(i/5)), uint8(128), 0)
	}
	l.SetDelay(1000)
	l.Print("./first.bin")

}
