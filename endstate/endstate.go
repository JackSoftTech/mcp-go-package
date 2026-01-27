package victoryTools

type EndStateResult struct {
	opts  EndStateParams
	Reach bool
}

type EndState interface {
	Key() string
	Info() string
}

type EndStateCondition interface {
	EndState() EndState
	ReachEndState(score int) bool
}
