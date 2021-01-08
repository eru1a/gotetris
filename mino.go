package main

type MinoKind int

const (
	MinoNone MinoKind = iota
	I
	O
	S
	Z
	J
	L
	T
)

type Dir int

const (
	North Dir = iota
	East
	South
	West
)

func (d Dir) RotateRight() Dir {
	if d != West {
		return Dir(d + 1)
	}
	return North
}

func (d Dir) RotateLeft() Dir {
	if d != North {
		return Dir(d - 1)
	}
	return West
}

type Mino struct {
	X, Y int // 左上の座標
	Kind MinoKind
	Dir  Dir
}

func NewMino(k MinoKind) *Mino {
	return &Mino{
		X:    3,
		Y:    2, // 20段目
		Kind: k,
		Dir:  North,
	}
}

func (m *Mino) RotateRight() {
	m.Dir = m.Dir.RotateRight()
}

func (m *Mino) RotateLeft() {
	m.Dir = m.Dir.RotateLeft()
}

func (m *Mino) Clone() *Mino {
	return &Mino{
		X:    m.X,
		Y:    m.Y,
		Kind: m.Kind,
		Dir:  m.Dir,
	}
}

func (m *Mino) Blocks() [4][4]bool {
	return minoBlocksMap[m.Kind][m.Dir]
}

type Pos struct {
	X, Y int
}

func (m *Mino) BlocksPos() []Pos {
	ret := []Pos{}
	if m.Kind == MinoNone {
		return ret
	}

	for i, blocks := range m.Blocks() {
		for j, b := range blocks {
			if !b {
				continue
			}
			x := j + m.X
			y := i + m.Y
			ret = append(ret, Pos{x, y})
		}
	}
	return ret
}

const (
	o = true
	x = false
)

var minoBlocksMap = map[MinoKind][4][4][4]bool{
	I: {
		{
			{x, x, x, x},
			{o, o, o, o},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, x, o, x},
			{x, x, o, x},
			{x, x, o, x},
			{x, x, o, x},
		},
		{
			{x, x, x, x},
			{x, x, x, x},
			{o, o, o, o},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{x, o, x, x},
			{x, o, x, x},
			{x, o, x, x},
		},
	},
	O: {
		{
			{x, o, o, x},
			{x, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, o, x},
			{x, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, o, x},
			{x, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, o, x},
			{x, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
	},
	S: {
		{
			{x, o, o, x},
			{o, o, x, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{x, o, o, x},
			{x, x, o, x},
			{x, x, x, x},
		},
		{
			{x, x, x, x},
			{x, o, o, x},
			{o, o, x, x},
			{x, x, x, x},
		},
		{
			{o, x, x, x},
			{o, o, x, x},
			{x, o, x, x},
			{x, x, x, x},
		},
	},
	Z: {
		{
			{o, o, x, x},
			{x, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, x, o, x},
			{x, o, o, x},
			{x, o, x, x},
			{x, x, x, x},
		},
		{
			{x, x, x, x},
			{o, o, x, x},
			{x, o, o, x},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{o, o, x, x},
			{o, x, x, x},
			{x, x, x, x},
		},
	},
	J: {
		{
			{o, x, x, x},
			{o, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, o, x},
			{x, o, x, x},
			{x, o, x, x},
			{x, x, x, x},
		},
		{
			{x, x, x, x},
			{o, o, o, x},
			{x, x, o, x},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{x, o, x, x},
			{o, o, x, x},
			{x, x, x, x},
		},
	},
	L: {
		{
			{x, x, o, x},
			{o, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{x, o, x, x},
			{x, o, o, x},
			{x, x, x, x},
		},
		{
			{x, x, x, x},
			{o, o, o, x},
			{o, x, x, x},
			{x, x, x, x},
		},
		{
			{o, o, x, x},
			{x, o, x, x},
			{x, o, x, x},
			{x, x, x, x},
		},
	},
	T: {
		{
			{x, o, x, x},
			{o, o, o, x},
			{x, x, x, x},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{x, o, o, x},
			{x, o, x, x},
			{x, x, x, x},
		},
		{
			{x, x, x, x},
			{o, o, o, x},
			{x, o, x, x},
			{x, x, x, x},
		},
		{
			{x, o, x, x},
			{o, o, x, x},
			{x, o, x, x},
			{x, x, x, x},
		},
	},
}
