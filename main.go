package main

import (
  "flag"
  "os"

  "github.com/gdamore/tcell"
)

const (
  Playing = iota
  Won
  Lost
)

var Levels = []Level{
  {9, 9, 10},
  {16, 16, 40},
  {16, 30, 99},
}

var (
  isEasy   bool
  isMedium bool
  isHard   bool
)

type Level struct {
  m     int
  n     int
  mines int
}

type Game struct {
  board  *Board
  cursor Cell
  screen tcell.Screen
  state  int
}

func init() {
  flag.BoolVar(&isEasy, "easy", false, "Play easy mode")
  flag.BoolVar(&isMedium, "medium", false, "Play medium mode")
  flag.BoolVar(&isHard, "hard", false, "Play hard mode")
}

func main() {
  flag.Parse()

  game := NewGame()

  for {
		switch event := game.screen.PollEvent().(type) {
		case *tcell.EventResize:
			game.screen.Sync()
		case *tcell.EventKey:
      if ok := game.update(event.Rune()); !ok {
        game.screen.Fini()
        os.Exit(0)
      }
		}
		game.render()
	}
}

func NewGame() *Game {
  level := Levels[0]
  if isMedium {
    level = Levels[1]
  } else if isHard {
    level = Levels[2]
  }

  board := NewBoard(level.m, level.n, level.mines)

  cursor := Cell{board.m / 2, board.n / 2}

  screen, err := tcell.NewScreen()
  if err != nil {
    panic(err)
  }

  if err := screen.Init(); err != nil {
    panic(err)
  }

  return &Game{board, cursor, screen, Playing}
}

func (g *Game) update(key rune) bool {
  // If game is over, exit on any keypress
  if g.state == Won || g.state == Lost {
    return false;
  }

  switch key {
  case 'q': // quit
    return false
  case 'k': // up
    g.cursor.i = max(0, g.cursor.i - 1)
  case 'j': // down
    g.cursor.i = min(g.cursor.i + 1, g.board.m - 1)
  case 'h': // left
    g.cursor.j = max(0, g.cursor.j - 1)
  case 'l': // right
    g.cursor.j = min(g.cursor.j + 1, g.board.n - 1)
  case 'f': // flag
    g.board.ToggleFlag(g.cursor)
  case 'e': // explore
    g.state = g.board.Explore(g.cursor)
    switch g.state {
    case Won:
      g.cursor = Cell{-1, -1}
    case Lost:
      g.board.Reveal()
    }
  }

  return true
}

func (g *Game) render() {
	g.screen.Clear()
	w, h := g.screen.Size()
  g.renderBoard((w - (2 * g.board.n - 1)) / 2, (h - g.board.m) / 2)
	g.screen.Show()
}

func (g *Game) renderString(s string, x, y int) {
	for i, c := range s {
		var comb []rune
    g.screen.SetContent(x + i, y, c, comb, tcell.StyleDefault)
	}
}

func (g *Game) renderBoard(x, y int) {
  for i := 0; i < g.board.m; i++ {
    for j := 0; j < g.board.n; j++ {
      c := g.board.visible[i][j]

      style := tcell.StyleDefault

      colors := []tcell.Color{
        tcell.ColorWhite,
        tcell.ColorBlue,
        tcell.ColorGreen,
        tcell.ColorRed,
        tcell.ColorNavy,
        tcell.ColorFireBrick,
        tcell.ColorMediumTurquoise,
        tcell.ColorBlack,
        tcell.ColorDarkGray,
      }

      switch c {
      case ' ', '.', '*':
        style = style.Foreground(tcell.ColorBlack)
      case '>':
        style = style.Foreground(tcell.ColorRed)
      default:
        idx := int(c - '0')
        style = style.Foreground(colors[idx])
      }

      if i == g.cursor.i && j == g.cursor.j {
        if g.state == Lost {
          style = style.Background(tcell.ColorRed)
        } else {
          style = tcell.StyleDefault.Reverse(true)
        }
      }

      var comb []rune
      g.screen.SetContent(x + 2 * j, y + i, c, comb, style)
    }
  }
}
