package main

const (
	BoardWidth  = 6
	BoardHeight = 6
)

type Board struct {
	Tiles         [BoardWidth][BoardHeight]int
	OpenPositions int
}

func NewBoard() *Board {
	return &Board{
		OpenPositions: 36,
	}
}

func (Board) IsLegalPosition(position *Position) bool {
	return position != nil &&
		position.X >= 0 && position.X < BoardWidth &&
		position.Y >= 0 && position.Y < BoardHeight
}

func (board *Board) SetPosition(position *Position) {

}
