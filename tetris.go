package main

import (
	"math/rand"
)

type Tetris struct {
	Board    [22][10]MinoKind
	CurMino  *Mino
	Hold     MinoKind
	Next     []MinoKind
	CanHold  bool
	GameOver bool
}

func NewTetris() *Tetris {
	t := &Tetris{
		Hold:    MinoNone,
		CanHold: true,
	}
	t.ShuffleNext()
	t.newMino()

	return t
}

func (t *Tetris) newMino() {
	t.CurMino = NewMino(t.Next[0])
	t.Next = t.Next[1:]
	t.ShuffleNext()

	// 新しいミノは基本的に20段目からスタートするが、
	// スタート位置にブロックがあったら21段目からスタートとなり、
	// さらに置けなければ22段目からスタートとなる
	// 22段目にも置けない場合はゲームオーバーとなる

	// 20段目
	if !t.minoCollides(t.CurMino) {
		return
	}

	// 21段目
	t.CurMino.Y--
	if !t.minoCollides(t.CurMino) {
		return
	}

	// 22段目
	t.CurMino.Y--
	if t.minoCollides(t.CurMino) {
		t.GameOver = true
	}
}

func (t *Tetris) Shadow() *Mino {
	shadow := t.CurMino.Clone()
	// 衝突するまで落下させる
	for {
		for _, p := range shadow.BlocksPos() {
			if t.collide(p.X, p.Y) {
				shadow.Y--
				goto after
			}
		}
		shadow.Y++
	}

after:
	return shadow
}

func (t *Tetris) collide(x, y int) bool {
	if x < 0 || x >= 10 || y < 0 || y >= 22 {
		return true
	}
	return t.Board[y][x] != MinoNone
}

func (t *Tetris) minoCollides(m *Mino) bool {
	for _, p := range m.BlocksPos() {
		if t.collide(p.X, p.Y) {
			return true
		}
	}
	return false
}

func (t *Tetris) Put() {
	shadow := t.Shadow()
	for _, p := range shadow.BlocksPos() {
		if p.Y >= 0 {
			t.Board[p.Y][p.X] = shadow.Kind
		}
	}
	t.newMino()
	t.CanHold = true
}

func (t *Tetris) DeleteLines() (n int) {
	for y := 0; y < 22; y++ {
		canDelete := true
		for x := 0; x < 10; x++ {
			if t.Board[y][x] == MinoNone {
				canDelete = false
			}
		}
		if canDelete {
			n++
			for yy := y; yy > 0; yy-- {
				for x := 0; x < 10; x++ {
					t.Board[yy][x] = t.Board[yy-1][x]
				}
			}
			for x := 0; x < 10; x++ {
				t.Board[0][x] = MinoNone
			}
		}
	}
	return n
}

func (t *Tetris) MoveLeft() (ok bool) {
	t.CurMino.X--
	for _, p := range t.CurMino.BlocksPos() {
		if t.collide(p.X, p.Y) {
			t.CurMino.X++
			return false
		}
	}
	return true
}

func (t *Tetris) MoveRight() (ok bool) {
	t.CurMino.X++
	for _, p := range t.CurMino.BlocksPos() {
		if t.collide(p.X, p.Y) {
			t.CurMino.X--
			return false
		}
	}
	return true
}

func (t *Tetris) MoveDown() (ok bool) {
	t.CurMino.Y++
	for _, p := range t.CurMino.BlocksPos() {
		if t.collide(p.X, p.Y) {
			t.CurMino.Y--
			return false
		}
	}
	return true
}

func (t *Tetris) RotateRight() (ok bool) {
	var xs, ys []int
	if t.CurMino.Kind != I {
		switch t.CurMino.Dir {
		case North:
			xs = []int{0, -1, -1, 0, -1}
			ys = []int{0, 0, -1, 2, 2}
		case East:
			xs = []int{0, 1, 1, 0, 1}
			ys = []int{0, 0, 1, -2, -2}
		case South:
			xs = []int{0, 1, 1, 0, 1}
			ys = []int{0, 0, -1, 2, 2}
		case West:
			xs = []int{0, -1, -1, 0, -1}
			ys = []int{0, 0, 1, -2, -2}
		}
	} else {
		switch t.CurMino.Dir {
		case North:
			xs = []int{0, -2, 1, -2, 1}
			ys = []int{0, 0, 0, 1, -2}
		case East:
			xs = []int{0, -1, 2, -1, 2}
			ys = []int{0, 0, 0, -2, 1}
		case South:
			xs = []int{0, 2, -1, 2, -1}
			ys = []int{0, 0, 0, -1, 2}
		case West:
			xs = []int{0, -2, 1, 1, -2}
			ys = []int{0, 0, 0, 2, -1}
		}
	}
	for i := range xs {
		dx := xs[i]
		dy := ys[i]
		clone := t.CurMino.Clone()
		clone.RotateRight()
		clone.X += dx
		clone.Y += dy
		if !t.minoCollides(clone) {
			t.CurMino = clone
			return true
		}
	}
	return true
}

func (t *Tetris) RotateLeft() (ok bool) {
	var xs, ys []int
	if t.CurMino.Kind != I {
		switch t.CurMino.Dir {
		case North:
			xs = []int{0, 1, 1, 0, 1}
			ys = []int{0, 0, -1, 2, 2}
		case East:
			xs = []int{0, 1, 1, 0, 1}
			ys = []int{0, 0, 1, -2, -2}
		case South:
			xs = []int{0, -1, -1, 0, -1}
			ys = []int{0, 0, -1, 2, 2}
		case West:
			xs = []int{0, -1, -1, 0, -1}
			ys = []int{0, 0, 1, -2, -2}
		}
	} else {
		switch t.CurMino.Dir {
		case North:
			xs = []int{0, -1, 2, -1, 2}
			ys = []int{0, 0, 0, -2, 1}
		case East:
			xs = []int{0, 2, -1, 2, -1}
			ys = []int{0, 0, 0, -1, 2}
		case South:
			xs = []int{0, 1, -2, 1, -2}
			ys = []int{0, 0, 0, 2, -1}
		case West:
			xs = []int{0, 1, -2, -2, 1}
			ys = []int{0, 0, 0, 1, -2}
		}
	}
	for i := range xs {
		dx := xs[i]
		dy := ys[i]
		clone := t.CurMino.Clone()
		clone.RotateLeft()
		clone.X += dx
		clone.Y += dy
		if !t.minoCollides(clone) {
			t.CurMino = clone
			return true
		}
	}
	return true
}

func (t *Tetris) DoHold() (ok bool) {
	if !t.CanHold {
		return false
	}

	t.CanHold = false
	if t.Hold == MinoNone {
		t.Hold = t.CurMino.Kind
		t.newMino()
		return true
	}

	t.CurMino, t.Hold = NewMino(t.Hold), t.CurMino.Kind
	return true
}

func (t *Tetris) ShuffleNext() {
	if len(t.Next) > 6 {
		return
	}
	minos := []MinoKind{I, O, S, Z, J, L, T}
	rand.Shuffle(7, func(i, j int) { minos[i], minos[j] = minos[j], minos[i] })
	for _, mino := range minos {
		t.Next = append(t.Next, mino)
	}
}
