package common

type Batch struct {
	Bets      []*Bet
	MaxAmount int
	Amount    int
	Size      int
}

func NewBatch(MaxAmount int) *Batch {
	return &Batch{
		Bets:      make([]*Bet, 0),
		MaxAmount: MaxAmount,
		Amount:    0,
		Size:      SIZE_HEADER_BYTES + SIZE_LENGTH_BYTES,
	}
}

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
