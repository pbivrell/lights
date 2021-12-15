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
			if r, _, _ := l.GetColor(light); r == 242 {
				l.SetColor(light, 172, 176, 175)
			} else if r, _, _ := l.GetColor(light); r == 172 {
				l.SetColor(light, 208, 211, 210)
			} else {
				l.SetColor(light, 242, 248, 247)
			}
		}
		l.SetDelay(500)
	}
	l.Print("crystal.bin")
}
