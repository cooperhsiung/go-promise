package go_promise

import (
	"sync/atomic"
)

// pending: initial state, neither fulfilled nor rejected.
// fulfilled: meaning that the operation was completed successfully.
// rejected: meaning that the operation failed.
type State string

const (
	StatePending   State = "pending"
	StateFulfilled State = "fulfilled"
	StateRejected  State = "rejected"
)

type Promise struct {
	val      interface{}
	err      error
	State    State
	resolver Resolver
	resolve  chan interface{}
	reject   chan error
}

type Resolver = func(resolve chan interface{}, reject chan error)

func New(resolver Resolver) *Promise {
	var p = &Promise{
		val:      nil,
		err:      nil,
		State:    StatePending,
		resolve:  make(chan interface{}, 1),
		reject:   make(chan error, 1),
		resolver: resolver,
	}
	// construct
	go resolver(p.resolve, p.reject)
	return p
}

func (p *Promise) Await() (interface{}, error) {
	// await once
	if p.val != nil || p.err != nil {
		return p.val, p.err
	}

	select {
	case val := <-p.resolve:
		p.val = val
		p.State = StateFulfilled
	case err := <-p.reject:
		p.State = StateRejected
		p.err = err
	}
	return p.val, p.err
}

func Resolve(val interface{}) *Promise {
	return &Promise{
		val:   val,
		err:   nil,
		State: StateFulfilled,
	}
}

func Reject(err error) *Promise {
	return &Promise{
		val:   nil,
		err:   err,
		State: StateRejected,
	}
}

// It rejects when any of the input's promises rejects, with this first rejection reason.
func All(futures ...*Promise) ([]interface{}, error) {
	var result = make([]interface{}, len(futures))
	var err error
	var done = make(chan bool, 1)

	// race
	var count int32 = 0
	for idx := range futures {
		go func(idx int) {
			p := futures[idx]
			select {
			case val := <-p.resolve:
				result[idx] = val
			case err = <-p.reject:
				// return immediate when error occur
				done <- true
			}
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) == int32(len(futures)) {
				done <- true
			}
		}(idx)
	}

	<-done
	return result, err
}

func AllSettled(futures ...*Promise) ([]interface{}, []error) {
	var result = make([]interface{}, len(futures))
	var reasons = make([]error, len(futures))
	// max time
	for idx := range futures {
		p := futures[idx]
		select {
		case val := <-p.resolve:
			result[idx] = val
		case err := <-p.reject:
			reasons[idx] = err
		}
	}

	return result, reasons
}

func Race(futures ...*Promise) (interface{}, error) {
	var err error
	var val interface{}
	var done = make(chan bool, 1)

	// race
	for idx := range futures {
		go func(idx int) {
			p := futures[idx]
			select {
			case val = <-p.resolve:
				// return immediate when resolved
				done <- true
			case err = <-p.reject:
				// return immediate when error occur
				done <- true
			}
		}(idx)
	}
	<-done
	return val, err
}

func Any(futures ...*Promise) (interface{}, error) {
	var err error
	var val interface{}
	var done = make(chan bool, 1)

	// race
	var count int32 = 0
	for idx := range futures {
		go func(idx int) {
			p := futures[idx]

			select {
			case val = <-p.resolve:
				// return immediate when resolved
				done <- true
			case err = <-p.reject:
			}
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) == int32(len(futures)) {
				done <- true
			}
		}(idx)
	}

	<-done
	if val != nil {
		err = nil
	}
	return val, err
}

func Map[T any](input []T, mapper func(T) *Promise, concurrency int) ([]interface{}, error) {
	var result = make([]interface{}, len(input))
	var err error
	var done = make(chan bool, 1)
	var jobs = make(chan int, len(input))
	var count int32 = 0

	for i := 0; i < concurrency; i++ {
		go func() {

			for idx := range jobs {
				p := mapper(input[idx])

				select {
				case val := <-p.resolve:
					result[idx] = val
				case err = <-p.reject:
					// return immediate when error occur
					done <- true
				}

				atomic.AddInt32(&count, 1)
				if atomic.LoadInt32(&count) == int32(len(input)) {
					done <- true
				}

			}
		}()
	}

	for idx := range input {
		jobs <- idx
	}
	close(jobs)
	<-done
	return result, err
}
