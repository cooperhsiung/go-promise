
## go-promise

[![Build Status](https://github.com/cooperhsiung/go-promise/workflows/Run%20Tests/badge.svg?branch=master)](https://github.com/cooperhsiung/go-promise/actions?query=branch%3Amaster)
[![codecov](https://codecov.io/cooperhsiung/go-promise/branch/master/graph/badge.svg)](https://codecov.io/cooperhsiung/go-promise)
[![Go Report Card](https://goreportcard.com/badge/github.com/cooperhsiung/go-promise)](https://goreportcard.com/report/github.com/cooperhsiung/go-promise)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/cooperhsiung/go-promise/blob/master/LICENSE.md)

Promise / Future library for Go.

### import

```
go get github.com/cooperhsiung/go-promise
```

```go
import (
    Promise "github.com/cooperhsiung/go-promise"
)
```

### Usage


- #### Promise.New(resolver Resolver)

Return the Promise object
```go
func NewSleep(delay int, throw bool) *Promise {
	return Promise.New(func(resolve chan interface{}, reject chan error) {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		if throw {
			reject <- errors.New("SLEEP ERROR")
		} else {
			resolve <- "hello world"
		}
	})
}
```

- #### Promise.Await()

Retrieve the result
```go
var p = NewSleep(100, false)
assert.Equal(t, p.State, StatePending, "")

val, _ := p.Await()
assert.Equal(t, p.State, StateFulfilled, "")
assert.Equal(t, p.val, "hello world", "p.val equal")
assert.Equal(t, val, "hello world", "p.val equal")
```

- #### Promise.All(futures ...*Promise) ([]interface{}, error)

This returned promise fulfills when all of the input's promises fulfill (including when an empty iterable is passed), with an array of the fulfillment values. It rejects when any of the input's promises rejects, with this first rejection reason.
```go
var p1 = NewSleep(100, false)
var p2 = NewSleep(100, false)

val, err := All(p1, p2)
fmt.Println(val, err)
```

- #### Promise.AllSettled(futures ...*Promise) ([]interface{}, []error)

This returned promise fulfills when all of the input's promises settle (including when an empty iterable is passed), with an array of objects that describe the outcome of each promise.
```go
var p1 = NewSleep(200, false)
var p2 = NewSleep(100, true)

val, err := AllSettled(p1, p2)
fmt.Println(val, err)
```

- #### Promise.Race(futures ...*Promise) (interface{}, error)

This returned promise settles with the eventual state of the first promise that settles.
```go
var p1 = NewSleep(200, false)
var p2 = NewSleep(100, false)

val, err := Race(p1, p2)
fmt.Println(val, err)
```

- #### Promise.Any(futures ...*Promise) (interface{}, error)

This returned promise fulfills when any of the input's promises fulfills, with this first fulfillment value. It rejects when all of the input's promises reject
```go
var p1 = NewSleep(200, false)
var p2 = NewSleep(100, false)

val, err := Any(p1, p2)
fmt.Println(val, err)
```

- #### Promise.Map[T any](input []T, mapper func(interface{}) *Promise, concurrency int) ([]interface{}, error)

Promises returned by the mapper function are awaited for and the returned promise doesn't fulfill until all mapped promises have fulfilled as well. If any promise in the array is rejected, or any promise returned by the mapper function is rejected, the returned promise is rejected as well.

The concurrency limit applies to Promises returned by the mapper function and it basically limits the number of Promises created

```go
ret, err := Map([]int{100, 200, 400, 800}, func(i interface{}) *Promise {
    return New(func(resolve chan interface{}, reject chan error) {
        time.Sleep(time.Duration(i.(int)) * time.Millisecond)
        resolve <- "hello world"
    })
}, 2)

fmt.Println(ret, err)
```


### Tests

```
go test ./... -count=1 -v -covermode=atomic
# coverage: 100.0% of statements
```