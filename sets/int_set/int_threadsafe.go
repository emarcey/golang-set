package mapsetint

import (
	"sync"
)

type threadSafeIntSet struct {
	s threadUnsafeIntSet
	sync.RWMutex
}

func newThreadSafeIntSet() threadSafeIntSet {
	return threadSafeIntSet{s: newThreadUnsafeIntSet()}
}

func (set *threadSafeIntSet) Add(i int) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeIntSet) Contains(i ...int) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeIntSet) IsSubset(other IntSet) bool {
	o := other.(*threadSafeIntSet)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeIntSet) IsProperSubset(other IntSet) bool {
	o := other.(*threadSafeIntSet)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeIntSet) IsSuperset(other IntSet) bool {
	return other.IsSubset(set)
}

func (set *threadSafeIntSet) IsProperSuperset(other IntSet) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeIntSet) Union(other IntSet) IntSet {
	o := other.(*threadSafeIntSet)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeIntSet)
	ret := &threadSafeIntSet{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeIntSet) Intersect(other IntSet) IntSet {
	o := other.(*threadSafeIntSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeIntSet)
	ret := &threadSafeIntSet{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeIntSet) Difference(other IntSet) IntSet {
	o := other.(*threadSafeIntSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeIntSet)
	ret := &threadSafeIntSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeIntSet) SymmetricDifference(other IntSet) IntSet {
	o := other.(*threadSafeIntSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeIntSet)
	ret := &threadSafeIntSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeIntSet) Clear() {
	set.Lock()
	set.s = newThreadUnsafeIntSet()
	set.Unlock()
}

func (set *threadSafeIntSet) Remove(i int) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeIntSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeIntSet) Each(cb func(int) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeIntSet) Iter() <-chan int {
	ch := make(chan int)
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

func (set *threadSafeIntSet) Iterator() *IntIterator {
	iterator, ch, stopCh := newIntIterator()

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

func (set *threadSafeIntSet) Equal(other IntSet) bool {
	o := other.(*threadSafeIntSet)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeIntSet) Clone() IntSet {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeIntSet)
	ret := &threadSafeIntSet{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeIntSet) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeIntSet) PowerSet() IntSet {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeIntSet)
    set.RUnlock()

    ret := &threadSafeIntSet{s: newThreadUnsafeIntSet()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeIntSet)
        ret.Add(&threadSafeIntSet{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeIntSet) Pop() int {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeIntSet) CartesianProduct(other IntSet) IntSet {
    o := other.(*threadSafeIntSet)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeIntSet)
    ret := &threadSafeIntSet{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeIntSet) ToSlice() []int {
	keys := make([]int, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeIntSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeIntSet) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
