package mapsetuint64

// Uint64Set is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type Uint64Set interface {
	// Adds an element to the set. Returns whether
	// the item was added.
	Add(i uint64) bool

	// Returns the number of elements in the set.
	Cardinality() int

	// Removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Uint64Set

	// Returns whether the given items
	// are all in the set.
	Contains(i ...uint64) bool

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Uint64Set) Uint64Set

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Uint64Set) bool

	// Returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Uint64Set) Uint64Set

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other Uint64Set) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other Uint64Set) bool

	// Determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other Uint64Set) bool

	// Determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other Uint64Set) bool

	// Iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(uint64) bool)

	// Returns a channel of elements that you can
	// range over.
	Iter() <-chan uint64

	// Returns an Iterator object that you can
	// use to range over the set.
	Iterator() *Uint64Iterator

	// Remove a single element from the set.
	Remove(i uint64)

	// Provides a convenient string representation
	// of the current state of the set.
	String() string

	// Returns a new set with all elements which are
	// in either this set or the other set but not in both.
	//
	// Note that the argument to SymmetricDifference
	// must be of the same type as the receiver
	// of the method. Otherwise, SymmetricDifference
	// will panic.
	SymmetricDifference(other Uint64Set) Uint64Set

	// Returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the

	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other Uint64Set) Uint64Set

	// Pop removes and returns an arbitrary item from the set.
	Pop() uint64

	// Returns all subsets of a given set (Power Uint64Set).
	// Not yet supported
	// PowerSet() Uint64Set

	// Returns the Cartesian Product of two sets.
	// Not yet supported
	// CartesianProduct(other Uint64Set) Uint64Set

	// Returns the members of the set as a slice.
	ToSlice() []uint64
}

// NewUint64Set creates and returns a reference to an empty set.  Operations
// on the resulting set are thread-safe.
func NewUint64Set(s ...uint64) Uint64Set {
	set := newThreadSafeUint64Set()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewUint64SetWith creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewUint64SetWith(elts ...uint64) Uint64Set {
	return NewUint64SetFromSlice(elts)
}

// NewUint64SetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewUint64SetFromSlice(s []uint64) Uint64Set {
	a := NewUint64Set(s...)
	return a
}

// NewThreadUnsafeSet creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeUint64Set() Uint64Set {
	set := newThreadUnsafeUint64Set()
	return &set
}

// NewThreadUnsafeSetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewThreadUnsafeUint64SetFromSlice(s []uint64) Uint64Set {
	a := NewThreadUnsafeUint64Set()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
