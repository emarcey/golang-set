package mapsetint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeIntSet map[int]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  int
	Second int
}

func newThreadUnsafeIntSet() threadUnsafeIntSet {
	return make(threadUnsafeIntSet)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeIntSet) Add(i int) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeIntSet) Contains(i ...int) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeIntSet) IsSubset(other IntSet) bool {
	_ = other.(*threadUnsafeIntSet)
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

func (set *threadUnsafeIntSet) IsProperSubset(other IntSet) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeIntSet) IsSuperset(other IntSet) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeIntSet) IsProperSuperset(other IntSet) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeIntSet) Union(other IntSet) IntSet {
	o := other.(*threadUnsafeIntSet)

	unionedSet := newThreadUnsafeIntSet()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeIntSet) Intersect(other IntSet) IntSet {
	o := other.(*threadUnsafeIntSet)

	intersection := newThreadUnsafeIntSet()
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

func (set *threadUnsafeIntSet) Difference(other IntSet) IntSet {
	_ = other.(*threadUnsafeIntSet)

	difference := newThreadUnsafeIntSet()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeIntSet) SymmetricDifference(other IntSet) IntSet {
	_ = other.(*threadUnsafeIntSet)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeIntSet) Clear() {
	*set = newThreadUnsafeIntSet()
}

func (set *threadUnsafeIntSet) Remove(i int) {
	delete(*set, i)
}

func (set *threadUnsafeIntSet) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeIntSet) Each(cb func(int) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeIntSet) Iter() <-chan int {
	ch := make(chan int)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeIntSet) Iterator() *IntIterator {
	iterator, ch, stopCh := newIntIterator()

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

func (set *threadUnsafeIntSet) Equal(other IntSet) bool {
	_ = other.(*threadUnsafeIntSet)

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

func (set *threadUnsafeIntSet) Clone() IntSet {
	clonedSet := newThreadUnsafeIntSet()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeIntSet) String() string {
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

func (set *threadUnsafeIntSet) Pop() int {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeIntSet) PowerSet() IntSet {
	powSet := NewThreadUnsafeIntSet()
	nullset := newThreadUnsafeIntSet()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeIntSet()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeIntSet()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeIntSet)
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
func (set *threadUnsafeIntSet) CartesianProduct(other IntSet) IntSet {
	o := other.(*threadUnsafeIntSet)
	cartProduct := NewThreadUnsafeIntSet()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeIntSet) ToSlice() []int {
	keys := make([]int, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeIntSet) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeIntSet) UnmarshalJSON(b []byte) error {
	var i []int

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
