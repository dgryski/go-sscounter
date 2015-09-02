package sscounter

import (
	"math/rand"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {

	const err = 0.01
	const iters = 1e7

	c := New(err)

	r := uint32(rand.Int63())
	for i := 0; i < iters; i++ {

		// xorshift rng, using the values from the paper
		r ^= r << 6
		r ^= r >> 21
		r ^= r << 7

		c.Inc(r)
	}

	v := c.Val()

	if v < iters*(1-err) || iters*(1+err) < v {
		t.Errorf("c.Val() fall outside error bounds: %v", v)
	}
}

func TestCounterN(t *testing.T) {

	const err = 0.01
	const iters = 1e7

	c := New(err)

	const goroutines = 4

	var wg sync.WaitGroup

	for j := 0; j < goroutines; j++ {
		wg.Add(1)
		go func() {
			r := uint32(rand.Int63())
			for i := 0; i < iters/goroutines; i++ {

				// xorshift rng, using the values from the paper
				r ^= r << 6
				r ^= r >> 21
				r ^= r << 7

				c.Inc(r)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	v := c.Val()

	if v < iters*(1-err) || iters*(1+err) < v {
		t.Errorf("c.Val() fall outside error bounds: %v", v)
	}
}
