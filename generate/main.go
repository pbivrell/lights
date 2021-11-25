package main

import "fmt"

func main() {

	s := make(Sequence, 0)

	(&s).setCount(10)
	for i := 0; i < 10; i++ {
		if i%5 == 0 {
			(&s).setDelay(1000)
		}
		(&s).setColor(uint16(i%5), uint8(128*(i/5)), uint8(128), 0)
	}
	(&s).setDelay(1000)
	(&s).print()

}

type Sequence [][8]byte

func (s *Sequence) print() {
	for i := range *s {
		for j := 0; j < 8; j++ {
			fmt.Printf("%08b", (*s)[i][j])
		}
	}
	fmt.Println()
}

func (s *Sequence) setColor(index uint16, r, g, b uint8) {
	bytes := [8]byte{0x1, r, g, b, byte(index & 0xFF), byte(index >> 8), 0, 0}
	*s = append(*s, bytes)
}

func (s *Sequence) setDelay(delay uint32) {
	bytes := [8]byte{0x2, byte(0xFF & delay), byte(delay & 0xFF00 >> 8), byte(delay & 0xFF0000 >> 16), byte(delay >> 24), 0, 0, 0}
	*s = append(*s, bytes)
}

func (s *Sequence) setCount(count uint16) {
	bytes := [8]byte{0x3, byte(count & 0xFF), byte(count >> 8), 0, 0, 0, 0, 0}
	*s = append(*s, bytes)
}

func (s *Sequence) clear(index uint16) {
	bytes := [8]byte{0x3, byte(index & 0xFF), byte(index >> 8), 0, 0, 0, 0, 0}
	*s = append(*s, bytes)
}

func (s *Sequence) clearAll() {
	bytes := [8]byte{0x4, 0, 0, 0, 0, 0, 0, 0}
	*s = append(*s, bytes)
}
