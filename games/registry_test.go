package games

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testGame struct {
	info GameInfo
}

func (tg testGame) GetInitialStat(players []PlayerState) {
	panic("implement me")
}

func (tg testGame) HandleUpdate(g GameState) (GameState, error) {
	panic("implement me")
}

func (tg testGame) CanPlayerMove(playerIndex int, g GameState) bool {
	panic("implement me")
}

func (tg testGame) IsGameOver(g GameState) {
	panic("implement me")
}

func (tg testGame) Info() GameInfo {
	return tg.info
}

func makeGame(name string, version int) testGame {
	return testGame{info: GameInfo{
		Name: name,
		Version: version,
	}}
}


func TestGameRegistry_Register(t *testing.T) {
	gr := newRegistry()
	g1 := makeGame("myGame", 1)
	g2 := makeGame("mySecondGame", 1)
	g3 := makeGame("mySecondGame", 2)

	gr.Register(g1)
	gr.Register(g2)
	gr.Register(g3)
}

func TestGameRegistry_GetGame(t *testing.T) {
	req := require.New(t)
	gr := newRegistry()
	g1 := makeGame("myGame", 1)
	g2 := makeGame("mySecondGame", 1)
	g3 := makeGame("mySecondGame", 2)

	gr.Register(g1)
	gr.Register(g2)
	gr.Register(g3)

	res1, e1 := gr.GetGame("myGame", 1)
	req.NoError(e1)
	req.Equal(g1, res1)

	res2, e2 := gr.GetGame("mySecondGame", 1)
	req.NoError(e2)
	req.Equal(g2, res2)

	res3, e3 := gr.GetGame("mySecondGame", 2)
	req.NoError(e3)
	req.Equal(g3, res3)

	_, e4 := gr.GetGame("gameDoesNotExist", 1)
	req.Error(e4, "The game doesn't exist.")

	_, e5 := gr.GetGame("myGame", 2)
	req.Error(e5, "The game's version does not exist")
}