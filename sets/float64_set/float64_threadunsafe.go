package mapsetfloat64

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeFloat64Set map[float64]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  float64
	Second float64
}

func newThreadUnsafeFloat64Set() threadUnsafeFloat64Set {
	return make(threadUnsafeFloat64Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeFloat64Set) Add(i float64) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeFloat64Set) Contains(i ...float64) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeFloat64Set) IsSubset(other Float64Set) bool {
	_ = other.(*threadUnsafeFloat64Set)
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

func (set *threadUnsafeFloat64Set) IsProperSubset(other Float64Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeFloat64Set) IsSuperset(other Float64Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeFloat64Set) IsProperSuperset(other Float64Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeFloat64Set) Union(other Float64Set) Float64Set {
	o := other.(*threadUnsafeFloat64Set)

	unionedSet := newThreadUnsafeFloat64Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeFloat64Set) Intersect(other Float64Set) Float64Set {
	o := other.(*threadUnsafeFloat64Set)

	intersection := newThreadUnsafeFloat64Set()
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

func (set *threadUnsafeFloat64Set) Difference(other Float64Set) Float64Set {
	_ = other.(*threadUnsafeFloat64Set)

	difference := newThreadUnsafeFloat64Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeFloat64Set) SymmetricDifference(other Float64Set) Float64Set {
	_ = other.(*threadUnsafeFloat64Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeFloat64Set) Clear() {
	*set = newThreadUnsafeFloat64Set()
}

func (set *threadUnsafeFloat64Set) Remove(i float64) {
	delete(*set, i)
}

func (set *threadUnsafeFloat64Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeFloat64Set) Each(cb func(float64) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeFloat64Set) Iter() <-chan float64 {
	ch := make(chan float64)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeFloat64Set) Iterator() *Float64Iterator {
	iterator, ch, stopCh := newFloat64Iterator()

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

func (set *threadUnsafeFloat64Set) Equal(other Float64Set) bool {
	_ = other.(*threadUnsafeFloat64Set)

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

func (set *threadUnsafeFloat64Set) Clone() Float64Set {
	clonedSet := newThreadUnsafeFloat64Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeFloat64Set) String() string {
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

func (set *threadUnsafeFloat64Set) Pop() float64 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeFloat64Set) PowerSet() Float64Set {
	powSet := NewThreadUnsafeFloat64Set()
	nullset := newThreadUnsafeFloat64Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeFloat64Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeFloat64Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeFloat64Set)
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
func (set *threadUnsafeFloat64Set) CartesianProduct(other Float64Set) Float64Set {
	o := other.(*threadUnsafeFloat64Set)
	cartProduct := NewThreadUnsafeFloat64Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeFloat64Set) ToSlice() []float64 {
	keys := make([]float64, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeFloat64Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeFloat64Set) UnmarshalJSON(b []byte) error {
	var i []float64

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
