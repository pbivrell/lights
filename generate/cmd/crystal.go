package main

import (
	"math/rand"

	"github.com/pbivrell/lights/generate/lights"
)

func main() {

	count := 500

	l := lights.New(uint16(count))

	randPercent := 0.20

	for i := 0; i < count; i++ {
		l.SetColor(uint16(i), 242, 248, 247)
	}

	r := rand.New(rand.NewSource(99))

	for i := 0; i < 10; i++ {
		for j := int(float64(count) * randPercent); j > 0; j-- {
			light := uint16(r.Intn(count))
			if r, _, _ := l.GetColor(light); r == 2 {
				l.SetColor(light, 140, 140, 140)
			} else if r, _, _ := l.GetColor(light); r == 50 {
				l.SetColor(light, 50, 50, 50)
			} else {
				l.SetColor(light, 20, 20, 20)
			}
		}
		l.SetDelay(500)
	}
	l.Print("crystal.bin")
}
