package mapsetint16

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Int16Iterator struct {
	C    <-chan int16
	stop chan struct{}
}

// Stop stops the Int16Iterator, no further elements will be received on C, C will be closed.
func (i *Int16Iterator) Stop() {
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

// newInt16Iterator returns a new Int16Iterator instance together with its item and stop channels.
func newInt16Iterator() (*Int16Iterator, chan<- int16, <-chan struct{}) {
	itemChan := make(chan int16)
	stopChan := make(chan struct{})
	return &Int16Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
