package mapsetuint16

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Uint16Iterator struct {
	C    <-chan uint16
	stop chan struct{}
}

// Stop stops the Uint16Iterator, no further elements will be received on C, C will be closed.
func (i *Uint16Iterator) Stop() {
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

// newUint16Iterator returns a new Uint16Iterator instance together with its item and stop channels.
func newUint16Iterator() (*Uint16Iterator, chan<- uint16, <-chan struct{}) {
	itemChan := make(chan uint16)
	stopChan := make(chan struct{})
	return &Uint16Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
