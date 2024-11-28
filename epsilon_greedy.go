package bandit

import (
	"math/rand"
	"sync"
)

// EpsilonGreedy represents the epsilon greedy algorithm
type EpsilonGreedy struct {
	sync.RWMutex
	Epsilon float64   `json:"epsilon"`
	Counts  []int     `json:"counts"`
	Rewards []float64 `json:"values"`
}

// Init will initialise the counts and rewards with the provided number of arms
func (b *EpsilonGreedy) Init(nArms int) error {
	b.Lock()
	defer b.Unlock()

	if nArms < 1 {
		return ErrInvalidArms
	}
	b.Counts = make([]int, nArms)
	b.Rewards = make([]float64, nArms)
	return nil
}

// SelectArm chooses an arm that exploits if the value is more than the epsilon
// threshold, and explore if the value is less than epsilon
func (b *EpsilonGreedy) SelectArm(probability float64) int {
	b.RLock()
	defer b.RUnlock()

	// Exploit
	if probability > b.Epsilon {
		return maxMean(b.Counts, b.Rewards)
	}

	// Explore
	return rand.Intn(len(b.Rewards))
}

// Update will update an arm with some reward value,
// e.g. click = 1, no click = 0
func (b *EpsilonGreedy) Update(chosenArm int, reward float64) error {
	b.Lock()
	defer b.Unlock()

	if chosenArm < 0 || chosenArm >= len(b.Rewards) {
		return ErrArmsIndexOutOfRange
	}
	if reward < 0 {
		return ErrInvalidReward
	}

	b.Counts[chosenArm]++
	n := float64(b.Counts[chosenArm])

	oldRewards := b.Rewards[chosenArm]
	b.Rewards[chosenArm] = (oldRewards*(n-1) + reward) / n

	return nil
}

// GetCounts returns the counts
func (b *EpsilonGreedy) GetCounts() []int {
	b.RLock()
	defer b.RUnlock()

	sCopy := make([]int, len(b.Counts))
	copy(sCopy, b.Counts)
	return sCopy
}

// GetRewards returns the rewards
func (b *EpsilonGreedy) GetRewards() []float64 {
	b.RLock()
	defer b.RUnlock()

	sCopy := make([]float64, len(b.Rewards))
	copy(sCopy, b.Rewards)
	return sCopy
}

// NewEpsilonGreedy returns a pointer to the EpsilonGreedy struct
func NewEpsilonGreedy(epsilon float64, counts []int, rewards []float64) (*EpsilonGreedy, error) {
	if epsilon < 0 || epsilon > 1 {
		return nil, ErrInvalidEpsilon
	}
	if len(counts) != len(rewards) {
		return nil, ErrInvalidLength
	}

	return &EpsilonGreedy{
		Epsilon: epsilon,
		Rewards: rewards,
		Counts:  counts,
	}, nil
}
