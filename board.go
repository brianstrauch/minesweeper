package main

import (
  "math/rand"
  "time"
)

type Board struct {
  m        int
  n        int
  mines    int
  explored int
  visible  [][]rune
  answers  [][]rune
}

type Cell struct {
  i int
  j int
}

func NewBoard(m, n, mines int) *Board {
  b := &Board{m: m, n: n, mines: mines, explored: 0}
  b.generateVisible()
  return b
}

// Toggle between '.' and '>'
func (b *Board) ToggleFlag(c Cell) {
  switch b.visible[c.i][c.j] {
  case '.':
    b.visible[c.i][c.j] = '>'
  case '>':
    b.visible[c.i][c.j] = '.'
  }
}

func (b *Board) Explore(c Cell) int {
  if b.answers == nil {
    b.generateAnswers(c)
  }

  if b.answers[c.i][c.j] == '*' {
    return Lost
  }

  // Already visited
  if b.visible[c.i][c.j] != '.' {
    return Playing
  }

  b.visible[c.i][c.j] = b.answers[c.i][c.j]

  b.explored++
  if b.explored == b.m * b.n - b.mines {
    return Won
  }

  state := Playing
  if b.answers[c.i][c.j] == ' ' {
    for _, n := range b.getNeighbors(c) {
      state = max(state, b.Explore(n))
    }
  }

  return state
}

// Used at the end of the game to reveal the positions of the mines
func (b *Board) Reveal() {
  for i := 0; i < b.m; i++ {
    for j := 0; j < b.n; j++ {
      if b.visible[i][j] == '.' && b.answers[i][j] == '*' {
        b.visible[i][j] = '*'
      }
    }
  }
}

// Start the game with a completely unexplored board
func (b *Board) generateVisible() {
  b.visible = make([][]rune, b.m)
  for i := range b.visible {
    b.visible[i] = make([]rune, b.n)
    for j := range b.visible[i] {
      b.visible[i][j] = '.'
    }
  }
}

// Randomly place mines, leaving an empty 3x3 grid around the cursor
// Compute numbers for the remaining cells
func (b *Board) generateAnswers(c Cell) {
  b.answers = make([][]rune, b.m)
  counts := make([][]int, b.m)
  for i := range b.answers {
    b.answers[i] = make([]rune, b.n)
    counts[i] = make([]int, b.n)
    for j := range b.answers[i] {
      b.answers[i][j] = ' '
      counts[i][j] = 0
    }
  }

  rand.Seed(time.Now().UnixNano())
  remaining := b.mines

  for remaining > 0 {
    i := rand.Intn(b.m)
    j := rand.Intn(b.n)

    if b.answers[i][j] == '*' {
      continue
    }
    if abs(i - c.i) <= 1 && abs(j - c.j) <= 1 {
      continue
    }

    b.answers[i][j] = '*'
    remaining--

    for _, n := range b.getNeighbors(Cell{i, j}) {
      counts[n.i][n.j]++
    }
  }

  for i := 0; i < b.m; i++ {
    for j := 0; j < b.n; j++ {
      if b.answers[i][j] != '*' && counts[i][j] > 0 {
        b.answers[i][j] = rune('0' + counts[i][j])
      }
    }
  }
}

func (b Board) getNeighbors(c Cell) []Cell {
  var neighbors []Cell
  for i := max(0, c.i - 1); i <= min(c.i + 1, b.m - 1); i++ {
    for j := max(0, c.j - 1); j <= min(c.j + 1, b.n - 1); j++ {
      neighbors = append(neighbors, Cell{i, j})
    }
  }
  return neighbors
}
