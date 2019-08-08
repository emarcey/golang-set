package mapsetuint64

import (
	"sync"
)

type threadSafeUint64Set struct {
	s threadUnsafeUint64Set
	sync.RWMutex
}

func newThreadSafeUint64Set() threadSafeUint64Set {
	return threadSafeUint64Set{s: newThreadUnsafeUint64Set()}
}

func (set *threadSafeUint64Set) Add(i uint64) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeUint64Set) Contains(i ...uint64) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) IsSubset(other Uint64Set) bool {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) IsProperSubset(other Uint64Set) bool {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeUint64Set) IsSuperset(other Uint64Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeUint64Set) IsProperSuperset(other Uint64Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeUint64Set) Union(other Uint64Set) Uint64Set {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeUint64Set)
	ret := &threadSafeUint64Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) Intersect(other Uint64Set) Uint64Set {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeUint64Set)
	ret := &threadSafeUint64Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) Difference(other Uint64Set) Uint64Set {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeUint64Set)
	ret := &threadSafeUint64Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) SymmetricDifference(other Uint64Set) Uint64Set {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeUint64Set)
	ret := &threadSafeUint64Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeUint64Set()
	set.Unlock()
}

func (set *threadSafeUint64Set) Remove(i uint64) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeUint64Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeUint64Set) Each(cb func(uint64) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeUint64Set) Iter() <-chan uint64 {
	ch := make(chan uint64)
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

func (set *threadSafeUint64Set) Iterator() *Uint64Iterator {
	iterator, ch, stopCh := newUint64Iterator()

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

func (set *threadSafeUint64Set) Equal(other Uint64Set) bool {
	o := other.(*threadSafeUint64Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) Clone() Uint64Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeUint64Set)
	ret := &threadSafeUint64Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeUint64Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeUint64Set) PowerSet() Uint64Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeUint64Set)
    set.RUnlock()

    ret := &threadSafeUint64Set{s: newThreadUnsafeUint64Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeUint64Set)
        ret.Add(&threadSafeUint64Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeUint64Set) Pop() uint64 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeUint64Set) CartesianProduct(other Uint64Set) Uint64Set {
    o := other.(*threadSafeUint64Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeUint64Set)
    ret := &threadSafeUint64Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeUint64Set) ToSlice() []uint64 {
	keys := make([]uint64, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeUint64Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeUint64Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
