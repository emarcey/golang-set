package mapsetbool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeBoolSet map[bool]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  bool
	Second bool
}

func newThreadUnsafeBoolSet() threadUnsafeBoolSet {
	return make(threadUnsafeBoolSet)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeBoolSet) Add(i bool) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeBoolSet) Contains(i ...bool) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeBoolSet) IsSubset(other BoolSet) bool {
	_ = other.(*threadUnsafeBoolSet)
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

func (set *threadUnsafeBoolSet) IsProperSubset(other BoolSet) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeBoolSet) IsSuperset(other BoolSet) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeBoolSet) IsProperSuperset(other BoolSet) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeBoolSet) Union(other BoolSet) BoolSet {
	o := other.(*threadUnsafeBoolSet)

	unionedSet := newThreadUnsafeBoolSet()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeBoolSet) Intersect(other BoolSet) BoolSet {
	o := other.(*threadUnsafeBoolSet)

	intersection := newThreadUnsafeBoolSet()
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

func (set *threadUnsafeBoolSet) Difference(other BoolSet) BoolSet {
	_ = other.(*threadUnsafeBoolSet)

	difference := newThreadUnsafeBoolSet()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeBoolSet) SymmetricDifference(other BoolSet) BoolSet {
	_ = other.(*threadUnsafeBoolSet)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeBoolSet) Clear() {
	*set = newThreadUnsafeBoolSet()
}

func (set *threadUnsafeBoolSet) Remove(i bool) {
	delete(*set, i)
}

func (set *threadUnsafeBoolSet) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeBoolSet) Each(cb func(bool) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeBoolSet) Iter() <-chan bool {
	ch := make(chan bool)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeBoolSet) Iterator() *BoolIterator {
	iterator, ch, stopCh := newBoolIterator()

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

func (set *threadUnsafeBoolSet) Equal(other BoolSet) bool {
	_ = other.(*threadUnsafeBoolSet)

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

func (set *threadUnsafeBoolSet) Clone() BoolSet {
	clonedSet := newThreadUnsafeBoolSet()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeBoolSet) String() string {
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

func (set *threadUnsafeBoolSet) Pop() bool {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return false
}

/*
// Not yet supported
func (set *threadUnsafeBoolSet) PowerSet() BoolSet {
	powSet := NewThreadUnsafeBoolSet()
	nullset := newThreadUnsafeBoolSet()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeBoolSet()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeBoolSet()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeBoolSet)
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
func (set *threadUnsafeBoolSet) CartesianProduct(other BoolSet) BoolSet {
	o := other.(*threadUnsafeBoolSet)
	cartProduct := NewThreadUnsafeBoolSet()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeBoolSet) ToSlice() []bool {
	keys := make([]bool, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeBoolSet) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeBoolSet) UnmarshalJSON(b []byte) error {
	var i []bool

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
