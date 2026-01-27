package victoryTools

import "fmt"

type ScoreEndState struct {
	ScoreTarget int
}

func (v *ScoreEndState) Key() string {
	return "score"
}

func (v *ScoreEndState) Info() string {
	if v.ScoreTarget > 0 {
		return fmt.Sprintf("Score EndState (target=%d)", v.ScoreTarget)
	}
	return "Score EndState (No target set or invalid target)"
}

func (v *ScoreEndState) EndState() EndState {
	return v
}

func (v *ScoreEndState) ReachEndState(score int) bool {
	if v.ScoreTarget <= 0 {
		return false
	}
	return score >= v.ScoreTarget
}

func SetScoreEndState(res EndStateResult) EndState {
	target := res.opts.ScoreTarget
	if target <= 0 {
		target = 10
	}
	return &ScoreEndState{ScoreTarget: target}
}
