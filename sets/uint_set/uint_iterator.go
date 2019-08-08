package mapsetuint

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type UintIterator struct {
	C    <-chan uint
	stop chan struct{}
}

// Stop stops the UintIterator, no further elements will be received on C, C will be closed.
func (i *UintIterator) Stop() {
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

// newUintIterator returns a new UintIterator instance together with its item and stop channels.
func newUintIterator() (*UintIterator, chan<- uint, <-chan struct{}) {
	itemChan := make(chan uint)
	stopChan := make(chan struct{})
	return &UintIterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
