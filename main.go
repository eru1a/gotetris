package main

import (
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	tetris               *Tetris
	hold                 *ebiten.Image
	board                *ebiten.Image
	next                 *ebiten.Image
	counter              int
	needLockdown         bool
	lockdownCounter      int // 床についてからのカウンター
	lockdownDelayCounter int // 床についた状態の遊び時間
}

func NewGame() *Game {
	return &Game{
		tetris: NewTetris(),
		hold:   ebiten.NewImage(4*gs, 3*gs),
		board:  ebiten.NewImage(10*gs, 22*gs),
		next:   ebiten.NewImage(4*gs, 3*gs*6),
	}
}

func (g *Game) MoveLeft() (ok bool) {
	ok = g.tetris.MoveLeft()
	if ok && g.needLockdown {
		g.counter = 0
		g.lockdownCounter = 0
		g.lockdownDelayCounter++
	}
	return
}

func (g *Game) MoveRight() (ok bool) {
	ok = g.tetris.MoveRight()
	if ok && g.needLockdown {
		g.counter = 0
		g.lockdownCounter = 0
		g.lockdownDelayCounter++
	}
	return
}

func (g *Game) MoveDown() (ok bool) {
	ok = g.tetris.MoveDown()
	if ok {
		g.counter = 0
		g.lockdownCounter = 0
		g.lockdownDelayCounter = 0
		g.needLockdown = false
	}
	return
}

func (g *Game) RotateLeft() (ok bool) {
	ok = g.tetris.RotateLeft()
	if ok && g.needLockdown {
		g.counter = 0
		g.lockdownCounter = 0
		g.lockdownDelayCounter++
	}
	return
}

func (g *Game) RotateRight() (ok bool) {
	ok = g.tetris.RotateRight()
	if ok && g.needLockdown {
		g.counter = 0
		g.lockdownCounter = 0
		g.lockdownDelayCounter++
	}
	return
}

func (g *Game) Put() {
	g.tetris.Put()
	g.tetris.DeleteLines()
	g.counter = 0
	g.lockdownCounter = 0
	g.lockdownDelayCounter = 0
	g.needLockdown = false
	return
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if g.tetris.GameOver {
		return nil
	}

	// TODO: 色々ダメ

	if g.tetris.CurMino.Y == g.tetris.Shadow().Y {
		g.needLockdown = true
	}
	if g.needLockdown {
		g.lockdownCounter++
	}
	if g.needLockdown && g.tetris.CurMino.Y == g.tetris.Shadow().Y {
		if g.lockdownCounter >= 30 || g.lockdownDelayCounter >= 15 {
			g.Put()
			return nil
		}
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		g.MoveLeft()
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		g.MoveRight()
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		g.MoveDown()
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		g.RotateRight()
	case inpututil.IsKeyJustPressed(ebiten.KeyZ):
		g.RotateLeft()
	case inpututil.IsKeyJustPressed(ebiten.KeyX):
		g.RotateRight()
	case inpututil.IsKeyJustPressed(ebiten.KeyC):
		g.tetris.DoHold()
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		g.Put()
		return nil
	}

	switch {
	case inpututil.KeyPressDuration(ebiten.KeyLeft) > 12 &&
		inpututil.KeyPressDuration(ebiten.KeyLeft)%3 == 0:
		g.MoveLeft()
	case inpututil.KeyPressDuration(ebiten.KeyRight) > 12 &&
		inpututil.KeyPressDuration(ebiten.KeyRight)%3 == 0:
		g.MoveRight()
	case inpututil.KeyPressDuration(ebiten.KeyDown) > 12 &&
		inpututil.KeyPressDuration(ebiten.KeyDown)%3 == 0:
		g.MoveDown()
	}

	// 自然落下
	g.counter++
	if g.counter%60 == 0 {
		if g.tetris.CurMino.Y != g.tetris.Shadow().Y || !g.needLockdown {
			if ok := g.MoveDown(); !ok {
				g.Put()
			}
		}
	}

	return nil
}

func MinoColor(m MinoKind) color.RGBA {
	switch m {
	case MinoNone:
		return color.RGBA{0, 0, 0, 255}
	case I:
		return color.RGBA{204, 255, 255, 255}
	case O:
		return color.RGBA{255, 204, 102, 255}
	case S:
		return color.RGBA{102, 204, 102, 255}
	case Z:
		return color.RGBA{226, 38, 58, 255}
	case J:
		return color.RGBA{0, 51, 204, 255}
	case L:
		return color.RGBA{232, 117, 40, 255}
	case T:
		return color.RGBA{102, 51, 153, 255}
	default:
		panic("not reach")
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.hold.Clear()
	g.board.Clear()
	g.next.Clear()

	// hold
	if g.tetris.Hold != MinoNone {
		mino := NewMino(g.tetris.Hold)
		mino.X = 0
		mino.Y = 0
		for _, p := range mino.BlocksPos() {
			ebitenutil.DrawRect(g.hold, float64(p.X*gs), float64(p.Y*gs),
				float64(gs), float64(gs), MinoColor(g.tetris.Hold))
		}
	}

	// board
	for y := 0; y < 22; y++ {
		for x := 0; x < 10; x++ {
			ebitenutil.DrawRect(g.board, float64(x*gs), float64(y*gs),
				float64(gs), float64(gs), MinoColor(g.tetris.Board[y][x]))
		}
	}
	if !g.tetris.GameOver {
		shadow := g.tetris.Shadow()
		shadowColor := MinoColor(shadow.Kind)
		shadowColor.A = 50
		for _, p := range shadow.BlocksPos() {
			ebitenutil.DrawRect(g.board, float64(p.X*gs), float64(p.Y*gs),
				float64(gs), float64(gs), shadowColor)
		}

		for _, p := range g.tetris.CurMino.BlocksPos() {
			ebitenutil.DrawRect(g.board, float64(p.X*gs), float64(p.Y*gs),
				float64(gs), float64(gs), MinoColor(g.tetris.CurMino.Kind))
		}
	}

	for i, next := range g.tetris.Next {
		mino := NewMino(next)
		mino.X = 0
		mino.Y = 0
		for _, p := range mino.BlocksPos() {
			ebitenutil.DrawRect(g.next, float64(p.X*gs), float64(p.Y*gs+3*gs*i),
				float64(gs), float64(gs), MinoColor(next))
		}
	}

	// 21~22段目
	gray1 := color.RGBA{204, 204, 204, 30}
	ebitenutil.DrawRect(g.board, 0, 0, 10*gs, 2*gs, gray1)

	// 21段目の赤線
	ebitenutil.DrawLine(g.board, 0, 2*gs, 10*gs, 2*gs, color.RGBA{200, 0, 0, 255})

	// 枠
	gray2 := color.RGBA{204, 204, 204, 255}
	ebitenutil.DrawLine(g.board, 1, 0, 1, 22*gs, gray2)
	ebitenutil.DrawLine(g.board, 10*gs, 0, 10*gs, 22*gs, gray2)
	ebitenutil.DrawLine(g.board, 0, 22*gs-1, 10*gs, 22*gs-1, gray2)

	// hold
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(gs, 0)
	screen.DrawImage(g.hold, op)

	// board
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(margin+4*gs+margin, 0)
	screen.DrawImage(g.board, op)

	// next
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(margin+4*gs+margin+10*gs+margin, 0)
	screen.DrawImage(g.next, op)

	if g.tetris.GameOver {
		ebitenutil.DebugPrint(screen, "GAME OVER")
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	gs           = 32
	margin       = gs
	screenWidth  = margin + 4*gs + margin + gs*10 + margin + 4*gs + margin
	screenHeight = 22 * gs
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
