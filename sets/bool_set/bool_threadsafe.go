package mapsetbool

import (
	"sync"
)

type threadSafeBoolSet struct {
	s threadUnsafeBoolSet
	sync.RWMutex
}

func newThreadSafeBoolSet() threadSafeBoolSet {
	return threadSafeBoolSet{s: newThreadUnsafeBoolSet()}
}

func (set *threadSafeBoolSet) Add(i bool) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeBoolSet) Contains(i ...bool) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) IsSubset(other BoolSet) bool {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) IsProperSubset(other BoolSet) bool {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeBoolSet) IsSuperset(other BoolSet) bool {
	return other.IsSubset(set)
}

func (set *threadSafeBoolSet) IsProperSuperset(other BoolSet) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeBoolSet) Union(other BoolSet) BoolSet {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeBoolSet)
	ret := &threadSafeBoolSet{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) Intersect(other BoolSet) BoolSet {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeBoolSet)
	ret := &threadSafeBoolSet{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) Difference(other BoolSet) BoolSet {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeBoolSet)
	ret := &threadSafeBoolSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) SymmetricDifference(other BoolSet) BoolSet {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeBoolSet)
	ret := &threadSafeBoolSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) Clear() {
	set.Lock()
	set.s = newThreadUnsafeBoolSet()
	set.Unlock()
}

func (set *threadSafeBoolSet) Remove(i bool) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeBoolSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeBoolSet) Each(cb func(bool) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeBoolSet) Iter() <-chan bool {
	ch := make(chan bool)
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

func (set *threadSafeBoolSet) Iterator() *BoolIterator {
	iterator, ch, stopCh := newBoolIterator()

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

func (set *threadSafeBoolSet) Equal(other BoolSet) bool {
	o := other.(*threadSafeBoolSet)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) Clone() BoolSet {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeBoolSet)
	ret := &threadSafeBoolSet{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeBoolSet) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeBoolSet) PowerSet() BoolSet {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeBoolSet)
    set.RUnlock()

    ret := &threadSafeBoolSet{s: newThreadUnsafeBoolSet()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeBoolSet)
        ret.Add(&threadSafeBoolSet{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeBoolSet) Pop() bool {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeBoolSet) CartesianProduct(other BoolSet) BoolSet {
    o := other.(*threadSafeBoolSet)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeBoolSet)
    ret := &threadSafeBoolSet{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeBoolSet) ToSlice() []bool {
	keys := make([]bool, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeBoolSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeBoolSet) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
