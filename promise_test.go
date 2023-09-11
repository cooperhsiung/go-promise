package go_promise

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var p = NewSleep(100, false)
	assert.Equal(t, p.State, StatePending, "")

	val, _ := p.Await()
	assert.Equal(t, p.State, StateFulfilled, "")
	assert.Equal(t, p.val, "hello world", "p.val equal")
	assert.Equal(t, val, "hello world", "p.val equal")
}

func TestResolve(t *testing.T) {
	var p = Resolve(1)
	assert.Equal(t, p.State, StateFulfilled, "")

	val, _ := p.Await()
	assert.Equal(t, p.State, StateFulfilled, "")
	assert.Equal(t, p.val, 1, "p.val equal")
	assert.Equal(t, val, 1, "p.val equal")
}

func TestReject(t *testing.T) {
	var p = Reject(errors.New("abcd"))
	assert.Equal(t, p.State, StateRejected, "")

	_, err := p.Await()
	assert.Equal(t, p.State, StateRejected, "")
	assert.Equal(t, p.err.Error(), "abcd", "error equal")
	assert.Equal(t, err.Error(), "abcd", "error equal")
}

func TestAwait(t *testing.T) {
	var p = NewSleep(100, true)
	assert.Equal(t, p.State, StatePending, "")

	_, err := p.Await()
	assert.Equal(t, p.err.Error(), "SLEEP ERROR", "error equal")
	assert.Equal(t, err.Error(), "SLEEP ERROR", "error equal")
}

func TestAll(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(100, false)
	var p2 = NewSleep(100, false)

	val, err := All(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-100)
}

func TestAllErr(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(200, false)
	var p2 = NewSleep(100, true)

	val, err := All(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-100)
	assert.NotEqual(t, nil, err, "error occur")
}

func TestAllSettled(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(200, false)
	var p2 = NewSleep(100, true)

	val, err := AllSettled(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-200)
	assert.NotEqual(t, nil, err, "error occur")
}

func TestRace(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(200, false)
	var p2 = NewSleep(100, false)

	val, err := Race(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-100)
}

func TestRaceErr(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(200, false)
	var p2 = NewSleep(100, true)

	val, err := Race(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-100)
	assert.NotEqual(t, nil, err, "error occur")
}

func TestAny(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(200, false)
	var p2 = NewSleep(100, false)

	val, err := Any(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-100)

	assert.Equal(t, nil, err, "error nil")
}

func TestAnyErr(t *testing.T) {
	start := time.Now()
	var p1 = NewSleep(200, false)
	var p2 = NewSleep(100, true)

	val, err := Any(p1, p2)
	fmt.Println(val, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-200)

	assert.Equal(t, nil, err, "error nil")
}

func TestMap(t *testing.T) {
	start := time.Now()
	ret, err := Map([]int{100, 200, 400, 800}, func(i int) *Promise {
		return New(func(resolve chan interface{}, reject chan error) {
			time.Sleep(time.Duration(i) * time.Millisecond)
			resolve <- "hello world"
		})
	}, 2)

	fmt.Println(ret, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-1000)
	assert.Equal(t, nil, err, "error nil")
}

func TestMap_timer1(t *testing.T) {
	start := time.Now()
	ret, err := Map([]int{200, 300, 100, 400}, func(i int) *Promise {
		return New(func(resolve chan interface{}, reject chan error) {
			time.Sleep(time.Duration(i) * time.Millisecond)
			resolve <- "hello world"
		})
	}, 2)

	fmt.Println(ret, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-700)
}

func TestMapErr(t *testing.T) {
	start := time.Now()
	ret, err := Map([]int{100, 200, 400, 600, 100}, func(i int) *Promise {
		return New(func(resolve chan interface{}, reject chan error) {
			var delay = i

			time.Sleep(time.Duration(delay) * time.Millisecond)
			if delay > 500 {
				reject <- errors.New("SLEEP ERROR")
			} else {
				resolve <- "hello world"
			}
		})
	}, 2)

	fmt.Println(ret, err)
	var cost = int(time.Since(start) / 1e6)
	fmt.Println("cost:", cost)
	AssertNear(t, cost-800)
	assert.NotEqual(t, nil, err, "error occur")
}

func AssertNear(t *testing.T, cost int) {
	assert.Equal(t, true, cost < 10 && cost >= 0, "time delta "+strconv.Itoa(cost))
}

func NewSleep(delay int, throw bool) *Promise {
	return New(func(resolve chan interface{}, reject chan error) {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		if throw {
			reject <- errors.New("SLEEP ERROR")
		} else {
			resolve <- "hello world"
		}
	})
}
