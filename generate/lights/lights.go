package lights

import (
	"encoding/binary"
	"os"
)

func New(count uint16) *Sequence {
	s := &Sequence{
		current: make([]Light, count),
		history: make([][8]byte, 0),
	}
	s.setCount(count)
	return s
}

type Light struct {
	R uint8
	G uint8
	B uint8
}

type Sequence struct {
	current []Light
	history [][8]byte
}

func (s *Sequence) Print(filepath string) error {

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	for i := range s.history {
		for j := 0; j < 8; j++ {
			err = binary.Write(file, binary.LittleEndian, s.history[i][j])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sequence) GetColor(index uint16) (r, g, b uint8) {
	c := s.current[index]
	return c.R, c.G, c.B
}

func (s *Sequence) SetColor(index uint16, r, g, b uint8) {
	if int(index) >= len(s.current) {
		return
	}
	bytes := [8]byte{0x1, r, g, b, byte(index & 0xFF), byte(index >> 8), 0, 0}
	s.current[index] = Light{R: r, G: g, B: b}
	s.history = append(s.history, bytes)
}

func (s *Sequence) SetDelay(delay uint32) {
	bytes := [8]byte{0x2, byte(0xFF & delay), byte(delay & 0xFF00 >> 8), byte(delay & 0xFF0000 >> 16), byte(delay >> 24), 0, 0, 0}
	s.history = append(s.history, bytes)
}

func (s *Sequence) setCount(count uint16) {
	bytes := [8]byte{0x3, byte(count & 0xFF), byte(count >> 8), 0, 0, 0, 0, 0}
	s.history = append(s.history, bytes)
}

func (s *Sequence) Fill(r, g, b uint8, index, count uint16) {
	bytes := [8]byte{0x4, r, g, b, byte(index & 0xFF), byte(index >> 8), byte(count & 0xFF), byte(count >> 8)}
	for i := index; i < count; i++ {
		s.current[i] = Light{R: r, G: g, B: b}
	}
	s.history = append(s.history, bytes)
}

/*func (s *Sequence) Clear(index uint16) {
	bytes := [8]byte{0x3, byte(index & 0xFF), byte(index >> 8), 0, 0, 0, 0, 0}
	s.history = append(s.history, bytes)
}

func (s *Sequence) ClearAll() {
	bytes := [8]byte{0x4, 0, 0, 0, 0, 0, 0, 0}
	s.history = append(s.history, bytes)
}*/
