package mapsetint16

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeInt16Set map[int16]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  int16
	Second int16
}

func newThreadUnsafeInt16Set() threadUnsafeInt16Set {
	return make(threadUnsafeInt16Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeInt16Set) Add(i int16) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeInt16Set) Contains(i ...int16) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeInt16Set) IsSubset(other Int16Set) bool {
	_ = other.(*threadUnsafeInt16Set)
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

func (set *threadUnsafeInt16Set) IsProperSubset(other Int16Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt16Set) IsSuperset(other Int16Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeInt16Set) IsProperSuperset(other Int16Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt16Set) Union(other Int16Set) Int16Set {
	o := other.(*threadUnsafeInt16Set)

	unionedSet := newThreadUnsafeInt16Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeInt16Set) Intersect(other Int16Set) Int16Set {
	o := other.(*threadUnsafeInt16Set)

	intersection := newThreadUnsafeInt16Set()
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

func (set *threadUnsafeInt16Set) Difference(other Int16Set) Int16Set {
	_ = other.(*threadUnsafeInt16Set)

	difference := newThreadUnsafeInt16Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeInt16Set) SymmetricDifference(other Int16Set) Int16Set {
	_ = other.(*threadUnsafeInt16Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeInt16Set) Clear() {
	*set = newThreadUnsafeInt16Set()
}

func (set *threadUnsafeInt16Set) Remove(i int16) {
	delete(*set, i)
}

func (set *threadUnsafeInt16Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeInt16Set) Each(cb func(int16) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeInt16Set) Iter() <-chan int16 {
	ch := make(chan int16)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeInt16Set) Iterator() *Int16Iterator {
	iterator, ch, stopCh := newInt16Iterator()

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

func (set *threadUnsafeInt16Set) Equal(other Int16Set) bool {
	_ = other.(*threadUnsafeInt16Set)

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

func (set *threadUnsafeInt16Set) Clone() Int16Set {
	clonedSet := newThreadUnsafeInt16Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeInt16Set) String() string {
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

func (set *threadUnsafeInt16Set) Pop() int16 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeInt16Set) PowerSet() Int16Set {
	powSet := NewThreadUnsafeInt16Set()
	nullset := newThreadUnsafeInt16Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeInt16Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeInt16Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeInt16Set)
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
func (set *threadUnsafeInt16Set) CartesianProduct(other Int16Set) Int16Set {
	o := other.(*threadUnsafeInt16Set)
	cartProduct := NewThreadUnsafeInt16Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeInt16Set) ToSlice() []int16 {
	keys := make([]int16, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeInt16Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeInt16Set) UnmarshalJSON(b []byte) error {
	var i []int16

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
