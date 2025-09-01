package common

import (
	"bufio"
)

type BatchIterator struct {
	agencyID   string
	scanner    *bufio.Scanner
	maxAmount  int
	currBatch  *Batch
	pendingBet *Bet
	done       bool
}

// CreateIteratorByScanner initializes and returns a new BatchIterator for the specified agency.
// It uses the provided bufio.Scanner to read input and processes batches up to maxAmount.
// The iterator maintains its internal state for batch processing and pending bets.
func CreateIteratorByScanner(agencyID string, scanner *bufio.Scanner, maxAmount int) *BatchIterator {
	return &BatchIterator{
		agencyID:   agencyID,
		scanner:    scanner,
		maxAmount:  maxAmount,
		currBatch:  nil,
		pendingBet: nil,
		done:       false,
	}
}

// Next reads the next batch from the scanner and returns it.
// Returns nil and false when there are no more batches.
func (it *BatchIterator) Next() (*Batch, bool) {
	if it.done {
		return nil, false
	}

	batch := NewBatch(it.maxAmount)
	if it.pendingBet != nil {
		batch.AddBet(it.agencyID, it.pendingBet)
		it.pendingBet = nil
	}

	for it.scanner.Scan() {
		line := it.scanner.Text()
		bet, err := ProcessCSVLine(line)
		if err != nil {
			log.Errorf("failed to process CSV line: %v", err)
			continue
		}

		if !batch.AddBet(it.agencyID, bet) {
			it.pendingBet = bet
			break
		}
	}

	if batch.Amount == 0 {
		it.done = true
		return nil, false
	}

	return batch, true
}
