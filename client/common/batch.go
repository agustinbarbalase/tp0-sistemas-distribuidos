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
		Amount:    SIZE_HEADER_BYTES + SIZE_LENGTH_BYTES,
		Size:      0,
	}
}

func (b *Batch) AddBet(ID string, bet *Bet) bool {
	if b.Amount+1 < b.MaxAmount && SIZE_LENGTH_BYTES+bet.PackagedBetLength(ID) < b.MaxSize {
		b.Size += bet.PackagedBetLength(ID)
		b.Amount++
		b.Bets = append(b.Bets, bet)
		return true
	}
	return false
}
