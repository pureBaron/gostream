package gostream

var (
	SequentialType = 1
	ConcurrentType = 2
)

func New[T any](slc []T, t int) Stream[T] {
	factories := map[int]Factory[T]{
		SequentialType: sequentialStreamFactory[T]{},
		ConcurrentType: concurrentStreamFactory[T]{},
	}

	return factories[t].New(slc)
}

type Factory[T any] interface {
	New(slc []T) Stream[T]
}

type sequentialStreamFactory[T any] struct{}

func (s sequentialStreamFactory[T]) New(slc []T) Stream[T] {
	return &sequentialStream[T]{
		slc: slc,
	}
}

type concurrentStreamFactory[T any] struct{}

func (c concurrentStreamFactory[T]) New(slc []T) Stream[T] {
	return &concurrentStream[T]{
		slc: slc,
	}
}
