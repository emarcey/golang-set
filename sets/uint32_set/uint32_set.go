package mapsetuint32

// Uint32Set is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type Uint32Set interface {
	// Adds an element to the set. Returns whether
	// the item was added.
	Add(i uint32) bool

	// Returns the number of elements in the set.
	Cardinality() int

	// Removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Uint32Set

	// Returns whether the given items
	// are all in the set.
	Contains(i ...uint32) bool

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Uint32Set) Uint32Set

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Uint32Set) bool

	// Returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Uint32Set) Uint32Set

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other Uint32Set) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other Uint32Set) bool

	// Determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other Uint32Set) bool

	// Determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other Uint32Set) bool

	// Iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(uint32) bool)

	// Returns a channel of elements that you can
	// range over.
	Iter() <-chan uint32

	// Returns an Iterator object that you can
	// use to range over the set.
	Iterator() *Uint32Iterator

	// Remove a single element from the set.
	Remove(i uint32)

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
	SymmetricDifference(other Uint32Set) Uint32Set

	// Returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the

	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other Uint32Set) Uint32Set

	// Pop removes and returns an arbitrary item from the set.
	Pop() uint32

	// Returns all subsets of a given set (Power Uint32Set).
	// Not yet supported
	// PowerSet() Uint32Set

	// Returns the Cartesian Product of two sets.
	// Not yet supported
	// CartesianProduct(other Uint32Set) Uint32Set

	// Returns the members of the set as a slice.
	ToSlice() []uint32
}

// NewUint32Set creates and returns a reference to an empty set.  Operations
// on the resulting set are thread-safe.
func NewUint32Set(s ...uint32) Uint32Set {
	set := newThreadSafeUint32Set()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewUint32SetWith creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewUint32SetWith(elts ...uint32) Uint32Set {
	return NewUint32SetFromSlice(elts)
}

// NewUint32SetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewUint32SetFromSlice(s []uint32) Uint32Set {
	a := NewUint32Set(s...)
	return a
}

// NewThreadUnsafeSet creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeUint32Set() Uint32Set {
	set := newThreadUnsafeUint32Set()
	return &set
}

// NewThreadUnsafeSetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewThreadUnsafeUint32SetFromSlice(s []uint32) Uint32Set {
	a := NewThreadUnsafeUint32Set()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
