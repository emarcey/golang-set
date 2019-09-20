package mapsetstring

import (
	"sync"
)

type threadSafeStringSet struct {
	s threadUnsafeStringSet
	sync.RWMutex
}

func newThreadSafeStringSet() threadSafeStringSet {
	return threadSafeStringSet{s: newThreadUnsafeStringSet()}
}

func (set *threadSafeStringSet) Add(i string) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeStringSet) Contains(i ...string) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeStringSet) IsSubset(other StringSet) bool {
	o := other.(*threadSafeStringSet)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeStringSet) IsProperSubset(other StringSet) bool {
	o := other.(*threadSafeStringSet)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeStringSet) IsSuperset(other StringSet) bool {
	return other.IsSubset(set)
}

func (set *threadSafeStringSet) IsProperSuperset(other StringSet) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeStringSet) Union(other StringSet) StringSet {
	o := other.(*threadSafeStringSet)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeStringSet)
	ret := &threadSafeStringSet{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeStringSet) Intersect(other StringSet) StringSet {
	o := other.(*threadSafeStringSet)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeStringSet)
	ret := &threadSafeStringSet{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeStringSet) Difference(other StringSet) StringSet {
	o := other.(*threadSafeStringSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeStringSet)
	ret := &threadSafeStringSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeStringSet) SymmetricDifference(other StringSet) StringSet {
	o := other.(*threadSafeStringSet)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeStringSet)
	ret := &threadSafeStringSet{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeStringSet) Clear() {
	set.Lock()
	set.s = newThreadUnsafeStringSet()
	set.Unlock()
}

func (set *threadSafeStringSet) Remove(i string) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeStringSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeStringSet) Each(cb func(string) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeStringSet) Iter() <-chan string {
	ch := make(chan string)
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

func (set *threadSafeStringSet) Iterator() *StringIterator {
	iterator, ch, stopCh := newStringIterator()

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

func (set *threadSafeStringSet) Equal(other StringSet) bool {
	o := other.(*threadSafeStringSet)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeStringSet) Clone() StringSet {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeStringSet)
	ret := &threadSafeStringSet{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeStringSet) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeStringSet) PowerSet() StringSet {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeStringSet)
    set.RUnlock()

    ret := &threadSafeStringSet{s: newThreadUnsafeStringSet()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeStringSet)
        ret.Add(&threadSafeStringSet{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeStringSet) Pop() string {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeStringSet) CartesianProduct(other StringSet) StringSet {
    o := other.(*threadSafeStringSet)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeStringSet)
    ret := &threadSafeStringSet{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeStringSet) ToSlice() []string {
	keys := make([]string, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeStringSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeStringSet) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
