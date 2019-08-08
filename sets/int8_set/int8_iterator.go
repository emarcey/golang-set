package mapsetint8

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Int8Iterator struct {
	C    <-chan int8
	stop chan struct{}
}

// Stop stops the Int8Iterator, no further elements will be received on C, C will be closed.
func (i *Int8Iterator) Stop() {
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

// newInt8Iterator returns a new Int8Iterator instance together with its item and stop channels.
func newInt8Iterator() (*Int8Iterator, chan<- int8, <-chan struct{}) {
	itemChan := make(chan int8)
	stopChan := make(chan struct{})
	return &Int8Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
