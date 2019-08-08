package mapsetuint8

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Uint8Iterator struct {
	C    <-chan uint8
	stop chan struct{}
}

// Stop stops the Uint8Iterator, no further elements will be received on C, C will be closed.
func (i *Uint8Iterator) Stop() {
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

// newUint8Iterator returns a new Uint8Iterator instance together with its item and stop channels.
func newUint8Iterator() (*Uint8Iterator, chan<- uint8, <-chan struct{}) {
	itemChan := make(chan uint8)
	stopChan := make(chan struct{})
	return &Uint8Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
