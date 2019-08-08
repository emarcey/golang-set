package mapsettimetime

import (
	"time"
)

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type TimeTimeIterator struct {
	C    <-chan time.Time
	stop chan struct{}
}

// Stop stops the TimeTimeIterator, no further elements will be received on C, C will be closed.
func (i *TimeTimeIterator) Stop() {
	// Allows for Stop() to be called multiple times
	// (close() panics when called on already closed channel)
	defer func() {
		recover()
	}()

	close(i.stop)

	// Exhaust any remaining elements.
	for range i.C {
	}
}

// newTimeTimeIterator returns a new TimeTimeIterator instance together with its item and stop channels.
func newTimeTimeIterator() (*TimeTimeIterator, chan<- time.Time, <-chan struct{}) {
	itemChan := make(chan time.Time)
	stopChan := make(chan struct{})
	return &TimeTimeIterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
