package mapsetuint64

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Uint64Iterator struct {
	C    <-chan uint64
	stop chan struct{}
}

// Stop stops the Uint64Iterator, no further elements will be received on C, C will be closed.
func (i *Uint64Iterator) Stop() {
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

// newUint64Iterator returns a new Uint64Iterator instance together with its item and stop channels.
func newUint64Iterator() (*Uint64Iterator, chan<- uint64, <-chan struct{}) {
	itemChan := make(chan uint64)
	stopChan := make(chan struct{})
	return &Uint64Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
