package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	CellSize     = 4
	CellsX       = WindowWidth / CellSize
	CellsY       = WindowHeight / CellSize
	PhysicsFPS   = 50
)

type CellType int

const (
	Empty CellType = iota
	Sand
	Water
	Stone
)

type Game struct {
	cells        [CellsX][CellsY]CellType
	actualChoice CellType
	physicsTick  int
}

// Update handles logic (Input + Physics)
func (g *Game) Update() error {
	// Handle Input: Selection
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.actualChoice = Sand
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.actualChoice = Water
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.actualChoice = Stone
	} else if inpututil.IsKeyJustPressed(ebiten.Key0) {
		g.actualChoice = Empty
	}

	// Handle Input: Drawing
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		g.placePixels(mx/CellSize, my/CellSize, 4, g.actualChoice)
	}

	// Physics Timing (Ebitengine runs at 60 TPS by default)
	g.physicsTick++
	if g.physicsTick >= (60 / PhysicsFPS) {
		g.physicsTick = 0
		g.updatePhysics()
	}

	return nil
}

func (g *Game) placePixels(cx, cy, size int, cellType CellType) {
	half := size / 2
	for x := -half; x < half; x++ {
		for y := -half; y < half; y++ {
			tx, ty := cx+x, cy+y
			if tx >= 0 && tx < CellsX && ty >= 0 && ty < CellsY {
				g.cells[tx][ty] = cellType
			}
		}
	}
}

func (g *Game) updatePhysics() {
	// Iterate backwards for Y so falling particles don't move multiple times in one tick
	for y := CellsY - 1; y >= 0; y-- {
		for x := 0; x < CellsX; x++ {
			cell := g.cells[x][y]

			if cell == Sand {
				if y+1 < CellsY {
					if g.cells[x][y+1] == Empty {
						g.cells[x][y+1] = Sand
						g.cells[x][y] = Empty
					} else if g.cells[x][y+1] == Water {
						g.cells[x][y+1] = Sand
						g.cells[x][y] = Water
					} else if g.cells[x][y+1] != Stone {
						left := x-1 >= 0 && g.cells[x-1][y+1] == Empty
						right := x+1 < CellsX && g.cells[x+1][y+1] == Empty

						if left && right {
							dir := -1
							if rand.Intn(2) == 0 {
								dir = 1
							}
							g.cells[x+dir][y+1] = Sand
							g.cells[x][y] = Empty
						} else if left {
							g.cells[x-1][y+1] = Sand
							g.cells[x][y] = Empty
						} else if right {
							g.cells[x+1][y+1] = Sand
							g.cells[x][y] = Empty
						}
					}
				}
			} else if cell == Water {
				if y+1 < CellsY && g.cells[x][y+1] == Empty {
					g.cells[x][y+1] = Water
					g.cells[x][y] = Empty
				} else {
					// Check diagonals
					left := x-1 >= 0 && y+1 < CellsY && g.cells[x-1][y+1] == Empty
					right := x+1 < CellsX && y+1 < CellsY && g.cells[x+1][y+1] == Empty

					if left && right {
						dir := -1
						if rand.Intn(2) == 0 {
							dir = 1
						}
						g.cells[x+dir][y+1] = Water
						g.cells[x][y] = Empty
					} else if left {
						g.cells[x-1][y+1] = Water
						g.cells[x][y] = Empty
					} else if right {
						g.cells[x+1][y+1] = Water
						g.cells[x][y] = Empty
					} else {
						// Horizontal flow
						leftH := x-1 >= 0 && g.cells[x-1][y] == Empty
						rightH := x+1 < CellsX && g.cells[x+1][y] == Empty

						if leftH && rightH {
							dir := -1
							if rand.Intn(2) == 0 {
								dir = 1
							}
							g.cells[x+dir][y] = Water
							g.cells[x][y] = Empty
						} else if leftH {
							g.cells[x-1][y] = Water
							g.cells[x][y] = Empty
						} else if rightH {
							g.cells[x+1][y] = Water
							g.cells[x][y] = Empty
						}
					}
				}
			}
		}
	}
}

// Draw handles rendering
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 64, 255})

	for x := 0; x < CellsX; x++ {
		for y := 0; y < CellsY; y++ {
			var c color.RGBA
			switch g.cells[x][y] {
			case Sand:
				c = color.RGBA{255, 165, 0, 255} // Orange
			case Water:
				c = color.RGBA{0, 0, 255, 255} // Blue
			case Stone:
				c = color.RGBA{128, 128, 128, 255} // Gray
			default:
				continue
			}

			vector.DrawFilledRect(
				screen,
				float32(x*CellSize),
				float32(y*CellSize),
				float32(CellSize),
				float32(CellSize),
				c,
				false,
			)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

func main() {
	game := &Game{
		actualChoice: Sand,
	}

	ebiten.SetWindowTitle("Sandbox - Ebitengine")
	ebiten.SetWindowSize(WindowWidth, WindowHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
