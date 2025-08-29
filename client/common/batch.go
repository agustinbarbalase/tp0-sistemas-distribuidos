package common

type Batch struct {
	Bets      []*Bet
	MaxAmount int
	MaxSize   int
	Amount    int
	Size      int
}

func NewBatch(MaxAmount int, MaxSize int) *Batch {
	return &Batch{
		Bets:      make([]*Bet, 0),
		MaxAmount: MaxAmount,
		MaxSize:   MaxSize,
		Amount:    0,
		Size:      SIZE_HEADER_BYTES + SIZE_LENGTH_BYTES,
	}
}

func (b *Batch) AddBet(ID string, bet *Bet) bool {
	lengthPackagedBet := bet.PackagedBetLength(ID)
	if b.Amount+1 <= b.MaxAmount && b.Size+lengthPackagedBet <= b.MaxSize {
		b.Size += lengthPackagedBet
		b.Amount++
		b.Bets = append(b.Bets, bet)
		return true
	}
	return false
}
