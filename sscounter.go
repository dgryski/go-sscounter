// Package sscounter implements scalable statistics counters
/*

https://blogs.oracle.com/dave/resource/spaa13-dice-ScalableStatsCounters.pdf

*/
package sscounter

import (
	"math"
	"sync/atomic"
	"unsafe"
)

const fmaxint32 float64 = math.MaxUint32

// Counter is a scalable counter
type Counter struct {
	threshold float64

	// these are read-only after constructions
	a          float64
	probFactor float64
}

// New returns a new scalable counter. rstdv is the relative standard
// deviation, i.e., the ratio of the standard deviation of the projected value
// and the actual count.
func New(rstdv float64) *Counter {
	a := 1 / (2 * rstdv * rstdv)
	probFactor := a / (a + 1)

	return &Counter{
		a:          a,
		probFactor: probFactor,
		threshold:  fmaxint32,
	}
}

// Val returns the current projected value of the counter
func (c *Counter) Val() int {
	pr := c.threshold / fmaxint32
	val := (1.0/pr - 1.0) * c.a
	return int(val)
}

// Increment the counter by 1.  Each call must provide r, a uniformly random uint32.
func (c *Counter) Inc(r uint32) {

	for {
		seenT := math.Float64frombits(atomic.LoadUint64((*uint64)(unsafe.Pointer(&c.threshold))))
		if r > uint32(seenT) {
			return
		}
		overflow := (seenT < c.a+1.0)
		newT := seenT * c.probFactor
		if overflow {
			newT = fmaxint32
		}
		if atomic.CompareAndSwapUint64((*uint64)(unsafe.Pointer(&c.threshold)), math.Float64bits(seenT), math.Float64bits(newT)) {
			return
		}
	}
}
