package gostream

import "sync"

type Stream[T any] interface {
	Filter(p Predicate[T]) Stream[T]
	Map(m Map[T]) Stream[T]
	ForEach(c Consumer[T])
	Reduce(bo BinaryOperator[T]) T
	ToSlice() []T
	Count() int
	AnyMatch(p Predicate[T]) bool
	AllMatch(p Predicate[T]) bool
}

type sequentialStream[T any] struct {
	slc []T
}

func (s sequentialStream[T]) Filter(p Predicate[T]) Stream[T] {
	res := []T{}
	for _, val := range s.slc {
		if p(val) {
			res = append(res, val)
		}
	}

	return sequentialStream[T]{
		slc: res,
	}

}

func (s sequentialStream[T]) Map(m Map[T]) Stream[T] {
	res := []T{}

	for _, val := range s.slc {
		res = append(res, m(val))
	}

	return sequentialStream[T]{
		slc: res,
	}
}

func (s sequentialStream[T]) ForEach(c Consumer[T]) {
	for _, val := range s.slc {
		c(val)
	}
}

func (s sequentialStream[T]) Reduce(bo BinaryOperator[T]) T {
	var res T

	for _, val := range s.slc {
		res = bo(res, val)
	}

	return res
}

func (s sequentialStream[T]) ToSlice() []T {
	return s.slc
}

func (s sequentialStream[T]) Count() int {
	return len(s.slc)
}

func (s sequentialStream[T]) AnyMatch(p Predicate[T]) bool {
	res := false

	for _, val := range s.slc {
		if p(val) {
			res = true
			break
		}
	}

	return res
}

func (s sequentialStream[T]) AllMatch(p Predicate[T]) bool {
	res := true

	for _, val := range s.slc {
		if !p(val) {
			res = false
			break
		}
	}

	if l := s.Count(); l == 0 {
		return false
	}
	return res
}

type concurrentStream[T any] struct {
	sync.Mutex
	slc []T
}

func (c *concurrentStream[T]) Filter(p Predicate[T]) Stream[T] {
	res := []T{}

	c.ForEach(func(t T) {
		if p(t) {
			c.Lock()
			res = append(res, t)
			c.Unlock()
		}
	})

	return &concurrentStream[T]{
		slc: res,
	}
}

func (c *concurrentStream[T]) Map(m Map[T]) Stream[T] {
	res := []T{}

	c.ForEach(func(t T) {
		c.Lock()
		defer c.Unlock()
		res = append(res, m(t))
	})

	return &concurrentStream[T]{
		slc: res,
	}
}

func (c *concurrentStream[T]) ForEach(cons Consumer[T]) {
	var wg sync.WaitGroup
	for _, val := range c.slc {
		wg.Add(1)
		go func(t T) {
			defer wg.Done()

			cons(t)

		}(val)
	}
	wg.Wait()
}

func (c *concurrentStream[T]) Reduce(bo BinaryOperator[T]) T {
	var res T

	c.ForEach(func(t T) {
		c.Lock()
		defer c.Unlock()
		res = bo(res, t)
	})

	return res
}

func (c *concurrentStream[T]) ToSlice() []T {
	return c.slc
}

func (c *concurrentStream[T]) Count() int {
	return len(c.slc)
}

func (c *concurrentStream[T]) AnyMatch(p Predicate[T]) bool {
	res := false

	c.ForEach(func(t T) {
		if p(t) {
			res = true
		}
	})

	return res
}

func (c *concurrentStream[T]) AllMatch(p Predicate[T]) bool {
	res := true

	c.ForEach(func(t T) {
		if !p(t) {
			res = false
		}
	})

	if l := c.Count(); l == 0 {
		return false
	}
	return res
}
