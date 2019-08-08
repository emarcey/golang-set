package mapsetuint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeUintSet map[uint]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  uint
	Second uint
}

func newThreadUnsafeUintSet() threadUnsafeUintSet {
	return make(threadUnsafeUintSet)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeUintSet) Add(i uint) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeUintSet) Contains(i ...uint) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUintSet) IsSubset(other UintSet) bool {
	_ = other.(*threadUnsafeUintSet)
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

func (set *threadUnsafeUintSet) IsProperSubset(other UintSet) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeUintSet) IsSuperset(other UintSet) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeUintSet) IsProperSuperset(other UintSet) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeUintSet) Union(other UintSet) UintSet {
	o := other.(*threadUnsafeUintSet)

	unionedSet := newThreadUnsafeUintSet()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeUintSet) Intersect(other UintSet) UintSet {
	o := other.(*threadUnsafeUintSet)

	intersection := newThreadUnsafeUintSet()
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

func (set *threadUnsafeUintSet) Difference(other UintSet) UintSet {
	_ = other.(*threadUnsafeUintSet)

	difference := newThreadUnsafeUintSet()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeUintSet) SymmetricDifference(other UintSet) UintSet {
	_ = other.(*threadUnsafeUintSet)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeUintSet) Clear() {
	*set = newThreadUnsafeUintSet()
}

func (set *threadUnsafeUintSet) Remove(i uint) {
	delete(*set, i)
}

func (set *threadUnsafeUintSet) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeUintSet) Each(cb func(uint) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeUintSet) Iter() <-chan uint {
	ch := make(chan uint)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeUintSet) Iterator() *UintIterator {
	iterator, ch, stopCh := newUintIterator()

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

func (set *threadUnsafeUintSet) Equal(other UintSet) bool {
	_ = other.(*threadUnsafeUintSet)

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

func (set *threadUnsafeUintSet) Clone() UintSet {
	clonedSet := newThreadUnsafeUintSet()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeUintSet) String() string {
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

func (set *threadUnsafeUintSet) Pop() uint {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeUintSet) PowerSet() UintSet {
	powSet := NewThreadUnsafeUintSet()
	nullset := newThreadUnsafeUintSet()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeUintSet()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeUintSet()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeUintSet)
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
func (set *threadUnsafeUintSet) CartesianProduct(other UintSet) UintSet {
	o := other.(*threadUnsafeUintSet)
	cartProduct := NewThreadUnsafeUintSet()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeUintSet) ToSlice() []uint {
	keys := make([]uint, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeUintSet) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeUintSet) UnmarshalJSON(b []byte) error {
	var i []uint

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
