package mapsetint8

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeInt8Set map[int8]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  int8
	Second int8
}

func newThreadUnsafeInt8Set() threadUnsafeInt8Set {
	return make(threadUnsafeInt8Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeInt8Set) Add(i int8) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeInt8Set) Contains(i ...int8) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeInt8Set) IsSubset(other Int8Set) bool {
	_ = other.(*threadUnsafeInt8Set)
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

func (set *threadUnsafeInt8Set) IsProperSubset(other Int8Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt8Set) IsSuperset(other Int8Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeInt8Set) IsProperSuperset(other Int8Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt8Set) Union(other Int8Set) Int8Set {
	o := other.(*threadUnsafeInt8Set)

	unionedSet := newThreadUnsafeInt8Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeInt8Set) Intersect(other Int8Set) Int8Set {
	o := other.(*threadUnsafeInt8Set)

	intersection := newThreadUnsafeInt8Set()
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

func (set *threadUnsafeInt8Set) Difference(other Int8Set) Int8Set {
	_ = other.(*threadUnsafeInt8Set)

	difference := newThreadUnsafeInt8Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeInt8Set) SymmetricDifference(other Int8Set) Int8Set {
	_ = other.(*threadUnsafeInt8Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeInt8Set) Clear() {
	*set = newThreadUnsafeInt8Set()
}

func (set *threadUnsafeInt8Set) Remove(i int8) {
	delete(*set, i)
}

func (set *threadUnsafeInt8Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeInt8Set) Each(cb func(int8) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeInt8Set) Iter() <-chan int8 {
	ch := make(chan int8)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeInt8Set) Iterator() *Int8Iterator {
	iterator, ch, stopCh := newInt8Iterator()

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

func (set *threadUnsafeInt8Set) Equal(other Int8Set) bool {
	_ = other.(*threadUnsafeInt8Set)

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

func (set *threadUnsafeInt8Set) Clone() Int8Set {
	clonedSet := newThreadUnsafeInt8Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeInt8Set) String() string {
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

func (set *threadUnsafeInt8Set) Pop() int8 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeInt8Set) PowerSet() Int8Set {
	powSet := NewThreadUnsafeInt8Set()
	nullset := newThreadUnsafeInt8Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeInt8Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeInt8Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeInt8Set)
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
func (set *threadUnsafeInt8Set) CartesianProduct(other Int8Set) Int8Set {
	o := other.(*threadUnsafeInt8Set)
	cartProduct := NewThreadUnsafeInt8Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeInt8Set) ToSlice() []int8 {
	keys := make([]int8, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeInt8Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeInt8Set) UnmarshalJSON(b []byte) error {
	var i []int8

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
