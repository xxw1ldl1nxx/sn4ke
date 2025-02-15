package sn4ke

import (
	"math/rand/v2"
	"slices"
)

var dirButton = map[string]Direction{
	"e": NONE,
	"у": NONE,

	"w": UP,
	"s": DOWN,
	"a": LEFT,
	"d": RIGHT,

	"ц": UP,
	"ы": DOWN,
	"ф": LEFT,
	"в": RIGHT,
}

type Direction = int
type gameMap = [2]int

// Directions for snake.
const (
	NONE Direction = iota
	UP
	DOWN
	LEFT
	RIGHT
)

// Fields of game matrix.
const (
	EMPTY int = iota
	TAIL
	HEAD
	APPLE
)

type Session struct {
	Snake  *Snake
	Apple  *Apple
	gMap   gameMap
	Input  func() Direction
	Output func([][]int, int)
	Score  int
}

func NewSession(x, y int, input func() Direction, output func([][]int, int)) *Session {
	gMap := [2]int{x, y}
	snake := NewSnake(gMap)
	apple := &Apple{
		gMap: gMap,
	}
	apple.setPos(*snake)
	return &Session{
		Snake:  snake,
		Apple:  apple,
		gMap:   gMap,
		Input:  input,
		Output: output,
		Score:  1, // because snake init length is 1
	}
}

func (s *Session) Update() error {
	dir := s.Input()

	eat, err := s.Snake.Move(dir, *s.Apple)
	if err != nil {
		return err
	}
	s.Score = s.Snake.Length()
	if eat {
		if err := s.Apple.setPos(*s.Snake); err == ErrBigSnake {
			return err
		}
	}
	return nil
}

func (s *Session) Draw() {
	matrix := s.newMatrix()

	s.Output(matrix, s.Score)

}

func (s Session) newMatrix() [][]int {
	matrix := make([][]int, 0, s.gMap[0])
	for i := 0; i < s.gMap[0]; i++ {
		innerSlice := make([]int, s.gMap[1])
		matrix = append(matrix, innerSlice)
	}

	head := s.Snake.Head()
	x, y := head[0], head[1]
	matrix[x][y] = HEAD

	for _, sn := range s.Snake.tail[:len(s.Snake.tail)-1] {
		x, y := sn[0], sn[1]
		matrix[x][y] = TAIL
	}

	ap := s.Apple.pos
	x, y = ap[0], ap[1]
	matrix[x][y] = APPLE

	return matrix
}

type Snake struct {
	tail    [][2]int
	bodyIdx map[[2]int]bool
	gMap    gameMap
	dir     Direction
}

func NewSnake(gMap gameMap) *Snake {
	x, y := gMap[0], gMap[1]
	tail := [][2]int{{x / 2, y / 2}}
	sIdx := make(map[[2]int]bool, x*y)
	sIdx[[2]int{x / 2, y / 2}] = true
	return &Snake{
		tail:    tail,
		bodyIdx: sIdx,
		gMap:    gMap,
		dir:     NONE,
	}
}

func (s Snake) Head() [2]int {
	return s.tail[len(s.tail)-1]
}

func (s Snake) Length() int {
	return len(s.tail)
}

func (s *Snake) Move(dir Direction, apple Apple) (bool, error) {

	if dir == NONE || s.oppositiveDir(dir) {
		dir = s.dir
	}

	newHead := s.newHead(dir)

	if slices.Contains(s.tail[1:], newHead) {
		return false, ErrSnakeDied
	}

	s.dir = dir
	s.bodyIdx[newHead] = true

	s.tail = append(s.tail, newHead)
	if newHead != apple.pos {
		delete(s.bodyIdx, s.tail[0])
		s.tail = s.tail[1:]
		return false, nil
	}
	return true, nil
}

func (s Snake) newHead(dir Direction) (newHead [2]int) {
	head := s.Head()
	x, y := head[0], head[1]
	switch dir {
	case UP:
		newHead = [2]int{(x - 1 + s.gMap[0]) % s.gMap[0], y}
	case DOWN:
		newHead = [2]int{(x + 1) % s.gMap[0], y}
	case LEFT:
		newHead = [2]int{x, (y - 1 + s.gMap[1]) % s.gMap[1]}
	case RIGHT:
		newHead = [2]int{x, (y + 1) % s.gMap[1]}
	case NONE:
		newHead = [2]int{x, y}
	}
	return
}

func (s Snake) oppositiveDir(dir Direction) bool {
	opDirs := map[Direction]Direction{
		UP:    DOWN,
		DOWN:  UP,
		LEFT:  RIGHT,
		RIGHT: LEFT,
	}
	return s.dir == opDirs[dir]
}

type Apple struct {
	pos  [2]int
	gMap gameMap
}

func (a *Apple) setPos(snake Snake) error {
	maxPos := (a.gMap[0] * a.gMap[1]) - snake.Length()
	if maxPos == 0 {
		return ErrBigSnake
	}
	occupied := snake.bodyIdx
	aviable := make([][2]int, 0, maxPos)

	for x := 0; x < a.gMap[0]; x++ {
		for y := 0; y < a.gMap[1]; y++ {
			pos := [2]int{x, y}
			if !occupied[pos] {
				aviable = append(aviable, pos)
			}
		}
	}

	a.pos = aviable[rand.IntN(maxPos)]
	return nil
}

// func main() {
// 	ss := NewSession(6, 6, Stdin, Stdout)
// 	for {
// 		ss.Draw()
// 		if err := ss.Update(); err == ErrSnakeDied || err == ErrBigSnake {
// 			break
// 		}
// 	}
// }
