package mapsetint64

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Int64Iterator struct {
	C    <-chan int64
	stop chan struct{}
}

// Stop stops the Int64Iterator, no further elements will be received on C, C will be closed.
func (i *Int64Iterator) Stop() {
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

// newInt64Iterator returns a new Int64Iterator instance together with its item and stop channels.
func newInt64Iterator() (*Int64Iterator, chan<- int64, <-chan struct{}) {
	itemChan := make(chan int64)
	stopChan := make(chan struct{})
	return &Int64Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
