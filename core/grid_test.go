package core

import (
	"testing"
	"fmt"
)

func TestGrid(t *testing.T) {
	player1 := "玩家1"
	player2 := "玩家2"

	g := NewGid(1,1,2,10,20)

	g.Add(1,player1)
	g.Add(2,player2)

	fmt.Println(g)

}
