package mapsetbool

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type BoolIterator struct {
	C    <-chan bool
	stop chan struct{}
}

// Stop stops the BoolIterator, no further elements will be received on C, C will be closed.
func (i *BoolIterator) Stop() {
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

// newBoolIterator returns a new BoolIterator instance together with its item and stop channels.
func newBoolIterator() (*BoolIterator, chan<- bool, <-chan struct{}) {
	itemChan := make(chan bool)
	stopChan := make(chan struct{})
	return &BoolIterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
