package main

import (
	"reflect"
	"strings"
	"testing"
)

func minoKindFromRune(c rune) MinoKind {
	switch c {
	case '.':
		return MinoNone
	case 'I':
		return I
	case 'O':
		return O
	case 'S':
		return S
	case 'Z':
		return Z
	case 'J':
		return J
	case 'L':
		return L
	case 'T':
		return T
	default:
		panic("not reach")
	}
}

func boardFromStr(s string) (ret [22][10]MinoKind) {
	for i, line := range strings.Split(s, "\n") {
		for j, c := range line {
			ret[i][j] = minoKindFromRune(c)
		}
	}
	return
}

func boardToStr(board [22][10]MinoKind) string {
	s := ""
	for y := 0; y < 22; y++ {
		line := ""
		for x := 0; x < 10; x++ {
			switch board[y][x] {
			case MinoNone:
				line += "."
			case I:
				line += "I"
			case O:
				line += "O"
			case S:
				line += "S"
			case Z:
				line += "Z"
			case J:
				line += "J"
			case L:
				line += "L"
			case T:
				line += "T"
			}
		}
		s += "\n" + line
	}
	return s
}

func TestRotateLeft(t *testing.T) {
	tests := []struct {
		msg      string
		tetris   *Tetris
		expected *Mino
	}{
		// TODO: 他にも色々テストする
		{
			msg: "T-Spin",
			tetris: &Tetris{
				Board: boardFromStr(`..........
..........
..........
..........
..........
..........
..........
..........
..........
..........
..........
..........
..........
..........
..........
..LL......
...L.....I
JJ.LT..OOI
J..TTTIOOI
J...SSIZZI
OO.SSLIJZZ
OO.LLLIJJJ`),
				CurMino: &Mino{
					X:    -1,
					Y:    14,
					Kind: T,
					Dir:  East,
				},
			},
			expected: &Mino{
				X:    1,
				Y:    17,
				Kind: T,
				Dir:  West,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			test.tetris.RotateLeft()
			test.tetris.RotateLeft()
			if !reflect.DeepEqual(test.expected, test.tetris.CurMino) {
				t.Errorf("should be %v, but got %v", test.expected, test.tetris.CurMino)
			}
		})
	}
}
