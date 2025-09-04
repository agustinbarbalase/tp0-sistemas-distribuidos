package common

type Batch struct {
	Bets      []*Bet
	MaxAmount int
	Amount    int
	Size      int
}

// NewBatch creates and returns a new Batch instance with the specified maximum amount of bets.
// It initializes the Bets slice, sets the MaxAmount, resets the Amount to zero, and sets the initial Size
// to the sum of SIZE_HEADER_BYTES and SIZE_LENGTH_BYTES.
func NewBatch(MaxAmount int) *Batch {
	return &Batch{
		Bets:      make([]*Bet, 0),
		MaxAmount: MaxAmount,
		Amount:    0,
		Size:      SIZE_HEADER_BYTES + SIZE_LENGTH_BYTES,
	}
}

// AddBet attempts to add a bet to the batch with the specified ID.
// It checks if adding the bet would exceed the batch's maximum amount or size constraints.
// If the bet can be added, it updates the batch's size and amount, appends the bet, and returns true.
// Otherwise, it returns false without modifying the batch.
func (b *Batch) AddBet(ID string, bet *Bet) bool {
	lengthPackagedBet := bet.PackagedBetLength(ID)
	if b.Amount+1 <= b.MaxAmount && b.Size+lengthPackagedBet <= MAX_SIZE_PACKAGE_IN_BYTES {
		b.Size += lengthPackagedBet
		b.Amount++
		b.Bets = append(b.Bets, bet)
		return true
	}
	return false
}
