package mapsetfloat32

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeFloat32Set map[float32]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  float32
	Second float32
}

func newThreadUnsafeFloat32Set() threadUnsafeFloat32Set {
	return make(threadUnsafeFloat32Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeFloat32Set) Add(i float32) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeFloat32Set) Contains(i ...float32) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeFloat32Set) IsSubset(other Float32Set) bool {
	_ = other.(*threadUnsafeFloat32Set)
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

func (set *threadUnsafeFloat32Set) IsProperSubset(other Float32Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeFloat32Set) IsSuperset(other Float32Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeFloat32Set) IsProperSuperset(other Float32Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeFloat32Set) Union(other Float32Set) Float32Set {
	o := other.(*threadUnsafeFloat32Set)

	unionedSet := newThreadUnsafeFloat32Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeFloat32Set) Intersect(other Float32Set) Float32Set {
	o := other.(*threadUnsafeFloat32Set)

	intersection := newThreadUnsafeFloat32Set()
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

func (set *threadUnsafeFloat32Set) Difference(other Float32Set) Float32Set {
	_ = other.(*threadUnsafeFloat32Set)

	difference := newThreadUnsafeFloat32Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeFloat32Set) SymmetricDifference(other Float32Set) Float32Set {
	_ = other.(*threadUnsafeFloat32Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeFloat32Set) Clear() {
	*set = newThreadUnsafeFloat32Set()
}

func (set *threadUnsafeFloat32Set) Remove(i float32) {
	delete(*set, i)
}

func (set *threadUnsafeFloat32Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeFloat32Set) Each(cb func(float32) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeFloat32Set) Iter() <-chan float32 {
	ch := make(chan float32)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeFloat32Set) Iterator() *Float32Iterator {
	iterator, ch, stopCh := newFloat32Iterator()

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

func (set *threadUnsafeFloat32Set) Equal(other Float32Set) bool {
	_ = other.(*threadUnsafeFloat32Set)

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

func (set *threadUnsafeFloat32Set) Clone() Float32Set {
	clonedSet := newThreadUnsafeFloat32Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeFloat32Set) String() string {
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

func (set *threadUnsafeFloat32Set) Pop() float32 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeFloat32Set) PowerSet() Float32Set {
	powSet := NewThreadUnsafeFloat32Set()
	nullset := newThreadUnsafeFloat32Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeFloat32Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeFloat32Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeFloat32Set)
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
func (set *threadUnsafeFloat32Set) CartesianProduct(other Float32Set) Float32Set {
	o := other.(*threadUnsafeFloat32Set)
	cartProduct := NewThreadUnsafeFloat32Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeFloat32Set) ToSlice() []float32 {
	keys := make([]float32, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeFloat32Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeFloat32Set) UnmarshalJSON(b []byte) error {
	var i []float32

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
