package gostream

type Predicate[T any] func(T) bool
type Consumer[T any] func(T)
type Map[T any] func(T) T
type BinaryOperator[T any] func(a, b T) T
type Comparator[T any] func(a, b T) int
