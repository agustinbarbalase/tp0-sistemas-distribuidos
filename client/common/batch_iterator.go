package common

import (
	"bufio"
)

type BatchIterator struct {
	agencyID        string
	maxAmount       int
	maxSize         int
	batches         []*Batch
	index           int
}

// CreateIteratorByScanner reads lines from the provided bufio.Scanner, processes each line as a bet,
// and groups bets into batches according to the specified maxAmount and maxSize constraints.
// Invalid CSV lines are skipped and logged as errors.
func CreateIteratorByScanner(scanner *bufio.Scanner, agencyID string, maxAmount int, maxSize int) *BatchIterator {
	var batches []*Batch

	batch := NewBatch(maxAmount, maxSize)

	for scanner.Scan() {
		line := scanner.Text()
		bet, err := ProcessCSVLine(line)
		if err != nil {
			log.Errorf("failed to process CSV line: %v", err)
			continue
		}

		if !batch.AddBet(agencyID, bet) {
			batch = NewBatch(maxAmount, maxSize)
			batch.AddBet(agencyID, bet)
		}
	}

	if batch.Amount > 0 {
		batches = append(batches, batch)
	}

	return &BatchIterator{
		agencyID:  agencyID,
		maxAmount: maxAmount,
		maxSize:   maxSize,
		batches:   batches,
		index:     0,
	}
}

// HasNext returns true if there are more batches to iterate over.
// It checks whether the current index is less than the total number of batches.
func (it *BatchIterator) HasNext() bool {
	return it.index < len(it.batches)
}

// GetCurrent returns the current Batch in the iterator without advancing the iterator.
// If there are no more batches, it returns nil.
func (it *BatchIterator) GetCurrent() *Batch {
	if !it.HasNext() {
		return nil
	}
	return it.batches[it.index]
}

// Next advances the iterator to the next batch by incrementing the index.
// If there are no more batches, it does nothing.
func (it *BatchIterator) Next() {
	if it.HasNext() {
		it.index++
	}
}
