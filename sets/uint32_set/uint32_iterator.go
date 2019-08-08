package mapsetuint32

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Uint32Iterator struct {
	C    <-chan uint32
	stop chan struct{}
}

// Stop stops the Uint32Iterator, no further elements will be received on C, C will be closed.
func (i *Uint32Iterator) Stop() {
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

// newUint32Iterator returns a new Uint32Iterator instance together with its item and stop channels.
func newUint32Iterator() (*Uint32Iterator, chan<- uint32, <-chan struct{}) {
	itemChan := make(chan uint32)
	stopChan := make(chan struct{})
	return &Uint32Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
