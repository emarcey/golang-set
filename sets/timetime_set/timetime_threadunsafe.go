package mapsettimetime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"time"
)

type threadUnsafeTimeTimeSet map[time.Time]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  time.Time
	Second time.Time
}

func newThreadUnsafeTimeTimeSet() threadUnsafeTimeTimeSet {
	return make(threadUnsafeTimeTimeSet)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeTimeTimeSet) Add(i time.Time) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeTimeTimeSet) Contains(i ...time.Time) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeTimeTimeSet) IsSubset(other TimeTimeSet) bool {
	_ = other.(*threadUnsafeTimeTimeSet)
	if set.Cardinality() > other.Cardinality() {
		return false
	}
	for elem := range *set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *threadUnsafeTimeTimeSet) IsProperSubset(other TimeTimeSet) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeTimeTimeSet) IsSuperset(other TimeTimeSet) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeTimeTimeSet) IsProperSuperset(other TimeTimeSet) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeTimeTimeSet) Union(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadUnsafeTimeTimeSet)

	unionedSet := newThreadUnsafeTimeTimeSet()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeTimeTimeSet) Intersect(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadUnsafeTimeTimeSet)

	intersection := newThreadUnsafeTimeTimeSet()
	// loop over smaller set
	if set.Cardinality() < other.Cardinality() {
		for elem := range *set {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range *o {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return &intersection
}

func (set *threadUnsafeTimeTimeSet) Difference(other TimeTimeSet) TimeTimeSet {
	_ = other.(*threadUnsafeTimeTimeSet)

	difference := newThreadUnsafeTimeTimeSet()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeTimeTimeSet) SymmetricDifference(other TimeTimeSet) TimeTimeSet {
	_ = other.(*threadUnsafeTimeTimeSet)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeTimeTimeSet) Clear() {
	*set = newThreadUnsafeTimeTimeSet()
}

func (set *threadUnsafeTimeTimeSet) Remove(i time.Time) {
	delete(*set, i)
}

func (set *threadUnsafeTimeTimeSet) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeTimeTimeSet) Each(cb func(time.Time) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeTimeTimeSet) Iter() <-chan time.Time {
	ch := make(chan time.Time)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeTimeTimeSet) Iterator() *TimeTimeIterator {
	iterator, ch, stopCh := newTimeTimeIterator()

	go func() {
	L:
		for elem := range *set {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
		close(ch)
	}()

	return iterator
}

func (set *threadUnsafeTimeTimeSet) Equal(other TimeTimeSet) bool {
	_ = other.(*threadUnsafeTimeTimeSet)

	if set.Cardinality() != other.Cardinality() {
		return false
	}
	for elem := range *set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *threadUnsafeTimeTimeSet) Clone() TimeTimeSet {
	clonedSet := newThreadUnsafeTimeTimeSet()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeTimeTimeSet) String() string {
	items := make([]string, 0, len(*set))

	for elem := range *set {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// String outputs a 2-tuple in the form "(A, B)".
func (pair OrderedPair) String() string {
	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
}

func (set *threadUnsafeTimeTimeSet) Pop() time.Time {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return time.Time{}
}

/*
// Not yet supported
func (set *threadUnsafeTimeTimeSet) PowerSet() TimeTimeSet {
	powSet := NewThreadUnsafeTimeTimeSet()
	nullset := newThreadUnsafeTimeTimeSet()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeTimeTimeSet()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeTimeTimeSet()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeTimeTimeSet)
				for ek := range *(k) {
					p.Add(ek)
				}
			} else {
				p.Add(er)
			}
			p.Add(es)
			u.Add(&p)
		}

		powSet = powSet.Union(&u)
	}

	return powSet
}
*/

/*
// Not yet supported
func (set *threadUnsafeTimeTimeSet) CartesianProduct(other TimeTimeSet) TimeTimeSet {
	o := other.(*threadUnsafeTimeTimeSet)
	cartProduct := NewThreadUnsafeTimeTimeSet()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeTimeTimeSet) ToSlice() []time.Time {
	keys := make([]time.Time, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeTimeTimeSet) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, set.Cardinality())

	for elem := range *set {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON recreates a set from a JSON array, it only decodes
// primitive types. Numbers are decoded as json.Number.
func (set *threadUnsafeTimeTimeSet) UnmarshalJSON(b []byte) error {
	var i []time.Time

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}

	for _, v := range i {
		set.Add(v)
	}

	return nil
}
