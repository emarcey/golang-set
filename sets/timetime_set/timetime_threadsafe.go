package mapsettimetime

import (
	"sync"
	"time"
)

type threadSafeTimeTimeSet struct {
	s threadUnsafeTimeTimeSet
	sync.RWMutex
}

func newThreadSafeTimeTimeSet() threadSafeTimeTimeSet {
	return threadSafeTimeTimeSet{s: newThreadUnsafeTimeTimeSet()}
}

func (set *threadSafeTimeTimeSet) Add(i time.Time) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeTimeTimeSet) Contains(i ...time.Time) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) IsSubset(other TimeTimeSet) bool {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) IsProperSubset(other TimeTimeSet) bool {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeTimeTimeSet) IsSuperset(other TimeTimeSet) bool {
	return other.IsSubset(set)
}

func (set *threadSafeTimeTimeSet) IsProperSuperset(other TimeTimeSet) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeTimeTimeSet) Union(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeTimeTimeSet)
	ret := &threadSafeTimeTimeSet{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) Intersect(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeTimeTimeSet)
	ret := &threadSafeTimeTimeSet{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) Difference(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeTimeTimeSet)
	ret := &threadSafeTimeTimeSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) SymmetricDifference(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeTimeTimeSet)
	ret := &threadSafeTimeTimeSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) Clear() {
	set.Lock()
	set.s = newThreadUnsafeTimeTimeSet()
	set.Unlock()
}

func (set *threadSafeTimeTimeSet) Remove(i time.Time) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeTimeTimeSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeTimeTimeSet) Each(cb func(time.Time) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeTimeTimeSet) Iter() <-chan time.Time {
	ch := make(chan time.Time)
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

func (set *threadSafeTimeTimeSet) Iterator() *TimeTimeIterator {
	iterator, ch, stopCh := newTimeTimeIterator()

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

func (set *threadSafeTimeTimeSet) Equal(other TimeTimeSet) bool {
	o := other.(*threadSafeTimeTimeSet)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) Clone() TimeTimeSet {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeTimeTimeSet)
	ret := &threadSafeTimeTimeSet{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeTimeTimeSet) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeTimeTimeSet) PowerSet() TimeTimeSet {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeTimeTimeSet)
    set.RUnlock()

    ret := &threadSafeTimeTimeSet{s: newThreadUnsafeTimeTimeSet()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeTimeTimeSet)
        ret.Add(&threadSafeTimeTimeSet{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeTimeTimeSet) Pop() time.Time {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeTimeTimeSet) CartesianProduct(other TimeTimeSet) TimeTimeSet {
    o := other.(*threadSafeTimeTimeSet)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeTimeTimeSet)
    ret := &threadSafeTimeTimeSet{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeTimeTimeSet) ToSlice() []time.Time {
	keys := make([]time.Time, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeTimeTimeSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeTimeTimeSet) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
