package main

import (
	"fmt"
)

type List[T any] interface {
	IsEmpty() bool
	Head() T
	Tail() List[T]
	ForEach(func(T))
}

type EmptyList[T any] struct{}

func (e EmptyList[T]) IsEmpty() bool     { return true }
func (e EmptyList[T]) Head() T           { panic("EmptyList has no Head") }
func (e EmptyList[T]) Tail() List[T]     { panic("EmptyList has no Tail") }
func (e EmptyList[T]) ForEach(f func(T)) {}

// Node is a non-empty list
type Node[T any] struct {
	val  T
	next List[T]
}

func (n Node[T]) IsEmpty() bool { return false }
func (n Node[T]) Head() T       { return n.val }
func (n Node[T]) Tail() List[T] { return n.next }
func (n Node[T]) ForEach(f func(T)) {
	f(n.val)
	n.next.ForEach(f)
}

// Prepend returns a new list with val added to the front of the existing list.
func Prepend[T any](l List[T], val T) List[T] {
	return Node[T]{val: val, next: l}
}

// Append returns a new list with val added to the end of the existing list.
func Append[T any](l List[T], val T) List[T] {
	if l.IsEmpty() {
		return Node[T]{val: val, next: EmptyList[T]{}}
	}
	n := l.(Node[T])
	return Node[T]{val: n.val, next: Append(n.next, val)}
}

// FromSlice constructs a List[T] from a slice by prepending elements in reverse order.
func FromSlice[T any](items []T) List[T] {
	var list List[T] = EmptyList[T]{}
	for _, val := range items {
		list = Prepend(list, val)
	}
	return list
}

// Map applies a function f to each element in the input list and returns a new list of the results.
func Map[F any, T any](l List[F], f func(F) T) List[T] {
	if l.IsEmpty() {
		var empty EmptyList[T]
		return empty
	}
	n := l.(Node[F])
	return Node[T]{val: f(n.val), next: Map(n.next, f)}
}

// Filter returns a new list containing only the elements that satisfy the predicate p.
func Filter[T any](l List[T], p func(T) bool) List[T] {
	if l.IsEmpty() {
		return l
	}
	n := l.(Node[T])
	if p(n.val) {
		return Node[T]{val: n.val, next: Filter(n.next, p)}
	} else {
		return Filter(n.next, p)
	}
}

// Fold reduces the list using a binary function f and an initial accumulator acc.
func Fold[F any, T any](l List[F], f func(F, T) T, acc T) T {
	if l.IsEmpty() {
		return acc
	}
	n := l.(Node[F])
	newAcc := f(n.val, acc)
	return Fold(n.next, f, newAcc)
}

// Length returns the number of elements in the list.
func Length[T any](l List[T]) uint {
	return Fold(
		l,
		func(v T, acc uint) uint { return acc + 1 },
		0,
	)
}

// Contains returns true if elem is found in the list; otherwise false.
func Contains[T comparable](l List[T], elem T) bool {
	if l.IsEmpty() {
		return false
	}
	n := l.(Node[T])
	if elem == n.val {
		return true
	} else {
		return Contains(n.next, elem)
	}
}

func main() {
	l := FromSlice([]int{1, 2, 3})

	fmt.Println("Original:")
	l.ForEach(func(x int) { fmt.Print(x, " ") })
	fmt.Println()

	fmt.Println("After Prepend 0:")
	prepended := Prepend(l, 0)
	prepended.ForEach(func(x int) { fmt.Print(x, " ") })
	fmt.Println()

	fmt.Println("After Append 4:")
	appended := Append(l, 4)
	appended.ForEach(func(x int) { fmt.Print(x, " ") })
	fmt.Println()

	fmt.Println("After Map (*2):")
	mapped := Map(l, func(x int) int { return x * 2 })
	mapped.ForEach(func(x int) { fmt.Print(x, " ") })
	fmt.Println()

	fmt.Println("After Filter (even):")
	filtered := Filter(l, func(x int) bool { return x%2 == 0 })
	filtered.ForEach(func(x int) { fmt.Print(x, " ") })
	fmt.Println()

	fmt.Println("After Fold (sum):")
	sum := Fold(l, func(x, acc int) int { return x + acc }, 0)
	fmt.Println(sum)

	fmt.Println("Length:")
	len := Length(l)
	fmt.Println(len)

	fmt.Println("Contains 42:")
	contains := Contains(l, 42)
	fmt.Println(contains)

}
