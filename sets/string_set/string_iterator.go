package mapsetstring

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type StringIterator struct {
	C    <-chan string
	stop chan struct{}
}

// Stop stops the StringIterator, no further elements will be received on C, C will be closed.
func (i *StringIterator) Stop() {
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

// newStringIterator returns a new StringIterator instance together with its item and stop channels.
func newStringIterator() (*StringIterator, chan<- string, <-chan struct{}) {
	itemChan := make(chan string)
	stopChan := make(chan struct{})
	return &StringIterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
