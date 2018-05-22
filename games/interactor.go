package games

import "fmt"

type GamesInteractor struct {
	store GameStore
}

func (inter GamesInteractor) CreateInstance(gameName, userId string) (GameInstance, error) {
	g, err := Registry.GetGameLatestVersion(gameName)

	if err != nil {
		return GameInstance{}, err
	}

	inst := NewInstance(g, userId)
	inst.AddPlayer(userId)
	err = inter.store.SaveInstance(inst)
	return *inst, err
}

func (inter GamesInteractor) JoinGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if inst.HasPlayer(userId) {
		return fmt.Errorf("%v is already in the game.", userId)
	}

	inst.AddPlayer(userId)

	return inter.store.SaveInstance(inst)
}

func (inter GamesInteractor) StartGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if inst.AdminUserId != userId {
		return fmt.Errorf("you're not the admin of this game")
	}

	inst.ShufflePlayers()
	inst.MetaState = InProgress
	inst.Game().InitializeState(&inst.State)

	return inter.store.SaveInstance(inst)
}

func (inter GamesInteractor) GetInstance(instanceId, userId string) (GameInstance, error) {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return GameInstance{}, nil
	}

	// If the game isn't in progress, we can return it as-is
	if inst.MetaState != InProgress {
		return *inst, nil
	}

	return transformInstanceForPlayer(*inst, userId)
}

func transformInstanceForPlayer(inst GameInstance, userId string) (GameInstance, error) {
	idx := int(inst.GetPlayerIndex(userId))
	// If the user isn't a player in the game,
	// act as if they're the user at position 0, keeping the list as it is.
	if idx < 0 {
		idx = 0
	}

	// Rotate the list of players so that the current user is first.
	// Also remove private state from other players
	ps := inst.State.Players
	n := len(ps)
	inst.State.Players = make([]PlayerState, n)
	for i, p := range ps {
		if p.UserId != userId {
			p.PrivateState = nil
		}
		inst.State.Players[(n + i - idx) % n] = p
	}
	return inst, nil
}