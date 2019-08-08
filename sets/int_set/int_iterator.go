package mapsetint

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type IntIterator struct {
	C    <-chan int
	stop chan struct{}
}

// Stop stops the IntIterator, no further elements will be received on C, C will be closed.
func (i *IntIterator) Stop() {
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

// newIntIterator returns a new IntIterator instance together with its item and stop channels.
func newIntIterator() (*IntIterator, chan<- int, <-chan struct{}) {
	itemChan := make(chan int)
	stopChan := make(chan struct{})
	return &IntIterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
