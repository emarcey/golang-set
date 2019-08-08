package mapsetint32

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Int32Iterator struct {
	C    <-chan int32
	stop chan struct{}
}

// Stop stops the Int32Iterator, no further elements will be received on C, C will be closed.
func (i *Int32Iterator) Stop() {
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

// newInt32Iterator returns a new Int32Iterator instance together with its item and stop channels.
func newInt32Iterator() (*Int32Iterator, chan<- int32, <-chan struct{}) {
	itemChan := make(chan int32)
	stopChan := make(chan struct{})
	return &Int32Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
