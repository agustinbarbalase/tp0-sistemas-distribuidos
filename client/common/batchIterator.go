package common

import (
	"bufio"
)

type BatchIterator struct {
	agencyID       string
	scanner        *bufio.Scanner
	maxAmount      int
	lastLineReaded string
}

// CreateIteratorByScanner reads lines from the provided bufio.Scanner, processes each line as a bet,
// and groups bets into batches according to the specified maxAmount and maxSize constraints.
// Invalid CSV lines are skipped and logged as errors.
func CreateIteratorByScanner(agencyID string, scanner *bufio.Scanner, maxAmount int) *BatchIterator {
	return &BatchIterator{
		agencyID:  agencyID,
		scanner:   scanner,
		maxAmount: maxAmount,
		lastLineReaded: "",
	}
}

// HasNext returns true if there are more batches to iterate over.
// It checks whether the current index is less than the total number of batches.
func (it *BatchIterator) HasNext() bool {
	return it.scanner.Text() != ""
}

// GetCurrent returns the current Batch in the iterator without advancing the iterator.
// If there are no more batches, it returns nil.
func (it *BatchIterator) GetCurrent() *Batch {
	batch := NewBatch(it.maxAmount)

	if it.lastLineReaded != "" {
		bet, err := ProcessCSVLine(it.lastLineReaded)
		if err != nil {
			log.Errorf("failed to process CSV line: %v", err)
		}
		batch.AddBet(it.agencyID, bet)
	}

	for it.scanner.Scan() {
		it.lastLineReaded = it.scanner.Text()
		bet, err := ProcessCSVLine(it.lastLineReaded)
		if err != nil {
			log.Errorf("failed to process CSV line: %v", err)
			continue
		}

		if !batch.AddBet(it.agencyID, bet) {
			batch = NewBatch(it.maxAmount)
			batch.AddBet(it.agencyID, bet)
			break
		}
	}

	return batch
}

// Next advances the iterator to the next batch by incrementing the index.
// If there are no more batches, it does nothing.
func (it *BatchIterator) Next() {
	it.scanner.Scan()
}
