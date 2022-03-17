package guid

import (
	"crypto/rand"
	"io"
	"sync"
	"time"
)

// Generator defines the contract for generating GUIDs
type Generator interface {
	Generate() (GUID, error)
}

// stdGenerator generates GUIDs
type stdGenerator struct {
	Fingerprint int32
	Random      io.Reader
	Now         func() time.Time
	IncrCounter int32
	DecrCounter int32

	incrLock        chan struct{}
	decrLock        chan struct{}
	createLocksOnce sync.Once
}

var (
	// Normally global variables are bad practice because they
	// introduce global state.  In this case, we want global state
	// because the generator counters must be locked globally, and
	// the fingerprint is global as well.
	// nolint: gochecknoglobals
	globalGenerator Generator = &stdGenerator{
		Random:      rand.Reader,
		Now:         time.Now,
		Fingerprint: defaultFingerprint(),
		IncrCounter: 0,
		DecrCounter: int32(time.Now().Unix() % int64(maxInt)),
	}

	setOnce sync.Once
)

// SetGlobalGenerator allows for the manual assignment of the GUID generator.
// The main usefulness of this function is primarily for testing, but
// this function can also be used to inject custom time and randomness
// providers.
// Note that this function can be called only once per runtime.
// Subsequent calls are no-ops.
func SetGlobalGenerator(g Generator) {
	setOnce.Do(func() {
		globalGenerator = g
	})
}

func (g *stdGenerator) randomInt32() (int32, error) {
	return randomInt32(g.Random)
}

// Generate will create a new GUID.
func (g *stdGenerator) Generate() (GUID, error) {
	g.createLocksOnce.Do(func() {
		g.incrLock = make(chan struct{}, 1)
		g.decrLock = make(chan struct{}, 1)
	})

	// increment counter lock
	g.incrLock <- struct{}{}
	incr := g.IncrCounter
	g.IncrCounter += 1
	if g.IncrCounter > maxInt {
		g.IncrCounter = 0
	}
	<-g.incrLock

	// get a random int32.  we sandwich this between
	// the locks to increase entropy
	r, err := g.randomInt32()
	if err != nil {
		return GUID{}, err
	}

	// decrement counter lock
	g.decrLock <- struct{}{}
	decr := g.DecrCounter
	g.DecrCounter -= 1
	if g.DecrCounter < 0 {
		g.DecrCounter = int32(maxInt)
	}
	<-g.decrLock

	v := (GUID{}).SetTime(g.Now()).SetCounters(incr, decr).SetFingerprint(g.Fingerprint).SetRandom(r)
	// set prefix bytes
	v[0] = globalPrefix[0]
	v[1] = globalPrefix[1]

	return v, nil
}
