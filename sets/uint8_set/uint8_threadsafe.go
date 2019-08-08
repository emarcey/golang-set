package mapsetuint8

import (
	"sync"
)

type threadSafeUint8Set struct {
	s threadUnsafeUint8Set
	sync.RWMutex
}

func newThreadSafeUint8Set() threadSafeUint8Set {
	return threadSafeUint8Set{s: newThreadUnsafeUint8Set()}
}

func (set *threadSafeUint8Set) Add(i uint8) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeUint8Set) Contains(i ...uint8) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) IsSubset(other Uint8Set) bool {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) IsProperSubset(other Uint8Set) bool {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeUint8Set) IsSuperset(other Uint8Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeUint8Set) IsProperSuperset(other Uint8Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeUint8Set) Union(other Uint8Set) Uint8Set {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeUint8Set)
	ret := &threadSafeUint8Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) Intersect(other Uint8Set) Uint8Set {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeUint8Set)
	ret := &threadSafeUint8Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) Difference(other Uint8Set) Uint8Set {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeUint8Set)
	ret := &threadSafeUint8Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) SymmetricDifference(other Uint8Set) Uint8Set {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeUint8Set)
	ret := &threadSafeUint8Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeUint8Set()
	set.Unlock()
}

func (set *threadSafeUint8Set) Remove(i uint8) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeUint8Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeUint8Set) Each(cb func(uint8) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeUint8Set) Iter() <-chan uint8 {
	ch := make(chan uint8)
	go func() {
		set.RLock()

		for elem := range set.s {
			ch <- elem
		}
		close(ch)
		set.RUnlock()
	}()

	return ch
}

func (set *threadSafeUint8Set) Iterator() *Uint8Iterator {
	iterator, ch, stopCh := newUint8Iterator()

	go func() {
		set.RLock()
	L:
		for elem := range set.s {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
		close(ch)
		set.RUnlock()
	}()

	return iterator
}

func (set *threadSafeUint8Set) Equal(other Uint8Set) bool {
	o := other.(*threadSafeUint8Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) Clone() Uint8Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeUint8Set)
	ret := &threadSafeUint8Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeUint8Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeUint8Set) PowerSet() Uint8Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeUint8Set)
    set.RUnlock()

    ret := &threadSafeUint8Set{s: newThreadUnsafeUint8Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeUint8Set)
        ret.Add(&threadSafeUint8Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeUint8Set) Pop() uint8 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeUint8Set) CartesianProduct(other Uint8Set) Uint8Set {
    o := other.(*threadSafeUint8Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeUint8Set)
    ret := &threadSafeUint8Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeUint8Set) ToSlice() []uint8 {
	keys := make([]uint8, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeUint8Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeUint8Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
