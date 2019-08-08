package mapsetuint

import (
	"sync"
)

type threadSafeUintSet struct {
	s threadUnsafeUintSet
	sync.RWMutex
}

func newThreadSafeUintSet() threadSafeUintSet {
	return threadSafeUintSet{s: newThreadUnsafeUintSet()}
}

func (set *threadSafeUintSet) Add(i uint) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeUintSet) Contains(i ...uint) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeUintSet) IsSubset(other UintSet) bool {
	o := other.(*threadSafeUintSet)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUintSet) IsProperSubset(other UintSet) bool {
	o := other.(*threadSafeUintSet)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeUintSet) IsSuperset(other UintSet) bool {
	return other.IsSubset(set)
}

func (set *threadSafeUintSet) IsProperSuperset(other UintSet) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeUintSet) Union(other UintSet) UintSet {
	o := other.(*threadSafeUintSet)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeUintSet)
	ret := &threadSafeUintSet{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUintSet) Intersect(other UintSet) UintSet {
	o := other.(*threadSafeUintSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeUintSet)
	ret := &threadSafeUintSet{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUintSet) Difference(other UintSet) UintSet {
	o := other.(*threadSafeUintSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeUintSet)
	ret := &threadSafeUintSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUintSet) SymmetricDifference(other UintSet) UintSet {
	o := other.(*threadSafeUintSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeUintSet)
	ret := &threadSafeUintSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUintSet) Clear() {
	set.Lock()
	set.s = newThreadUnsafeUintSet()
	set.Unlock()
}

func (set *threadSafeUintSet) Remove(i uint) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeUintSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeUintSet) Each(cb func(uint) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeUintSet) Iter() <-chan uint {
	ch := make(chan uint)
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

func (set *threadSafeUintSet) Iterator() *UintIterator {
	iterator, ch, stopCh := newUintIterator()

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

func (set *threadSafeUintSet) Equal(other UintSet) bool {
	o := other.(*threadSafeUintSet)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUintSet) Clone() UintSet {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeUintSet)
	ret := &threadSafeUintSet{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeUintSet) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeUintSet) PowerSet() UintSet {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeUintSet)
    set.RUnlock()

    ret := &threadSafeUintSet{s: newThreadUnsafeUintSet()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeUintSet)
        ret.Add(&threadSafeUintSet{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeUintSet) Pop() uint {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeUintSet) CartesianProduct(other UintSet) UintSet {
    o := other.(*threadSafeUintSet)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeUintSet)
    ret := &threadSafeUintSet{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeUintSet) ToSlice() []uint {
	keys := make([]uint, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeUintSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeUintSet) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
