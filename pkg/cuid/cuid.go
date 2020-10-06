// Package cuid contains a Go implementation of the CUID project defined at
// https://github.com/ericelliott/cuid.
package cuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	prefix     = "c"
	prefixByte = 'c'
	// padChar is used as a prefix value for any base36 encoding that does not
	// fit a certain length. For example, encoding the 32bit integer value of
	// 1 would result in a string of "1". However, the string must fit a four
	// character block in the string so we add the padChar to make "0001".
	padChar = "0"
	// padding matches https://github.com/ericelliott/cuid/blob/master/lib/pad.js
	// and is used to ensure that all values are preceded by zeros if they do
	// not meet the string length requirements. The original code does this to
	// avoid needing complex logic in the padding function by, instead, always
	// adding a large amount of padding to every string and then trimming away
	// any portions that are not needed. See the leftpad() method for an
	// example.
	padding = padChar + padChar + padChar + padChar + padChar + padChar + padChar + padChar + padChar
	// blockSize is the standard size string for each 32bit integer field.
	blockSize = 4
	// base is used for all encoding operations. CUIDs use a base36 encoding of
	// the binary data to generate a string.
	base = 36
	// maxInt is the highest value allowed for the counter and random fields.
	// This acts as a guard to ensure that all integer values encoded in the
	// CUID are within the acceptable range of the base36 encoding that is used
	// to generate a string representation. Reference:
	// https://github.com/ericelliott/cuid/blob/master/index.js#L20.
	maxInt = 1679616
)

var (
	// NOTE: I don't recommend this pattern of having global singletons and
	// locks. This is done only to provide better parity with existing CUID
	// and UUID libraries that also do this.
	globalLock      = &sync.RWMutex{} //nolint:gochecknoglobals
	globalGenerator = &Generator{     //nolint:gochecknoglobals
		Fingerprint: defaultFingerprint(),
		Random:      rand.Reader,
		Now:         time.Now,
		Counter:     0,
		Locker:      &sync.Mutex{},
	}
)

// SetGenerator changes the global Generator instance.
func SetGenerator(g *Generator) {
	globalLock.Lock()
	defer globalLock.Unlock()
	globalGenerator = g
}

// New generates a CUID using the global generator.
func New() (CUID, error) {
	globalLock.RLock()
	g := globalGenerator
	globalLock.RUnlock()
	return g.Generate()
}

// IsCUID determines if a given string is a valid CUID. The original version
// of this at https://github.com/ericelliott/cuid/blob/master/index.js#L69
// only checks the initial 'c' character which is not an accurate check for
// CUID validity. For example, the counter and random values could be out of
// range for what can be decoded using a base36 scheme but the original check
// would still consider those to be valid CUIDs.
//
// This version of IsCUID enforces all constraints of a CUID by parsing it.
// This is provided for parity with the JavaScript library API. If you plan
// to modify the CUID after validating then it is better to use the Parse
// variants which apply the same validation.
func IsCUID(s string) bool {
	_, err := ParseString(s)
	return err == nil
}

// ParseString attempts to create a CUID from the given string.
func ParseString(s string) (CUID, error) {
	return ParseBytes([]byte(s))
}

// ParseBytes attempts to create a CUID from the given byte string.
func ParseBytes(b []byte) (CUID, error) {
	if len(b) != 25 {
		return CUID{}, fmt.Errorf("CUID must be 25 characters. got %d", len(b))
	}
	if b[0] != prefixByte {
		return CUID{}, fmt.Errorf("CUID must start with 'c'. got %s", string(b[0]))
	}
	c := CUID{}
	c[0] = prefixByte

	t, err := strconv.ParseInt(string(b[1:9]), base, 64)
	if err != nil {
		return CUID{}, fmt.Errorf("invalid time string %s: %w", b[1:9], err)
	}
	c = c.SetTime(time.Unix(0, t*1e6))

	counter, err := strconv.ParseInt(string(b[9:13]), base, 64)
	if err != nil {
		return CUID{}, fmt.Errorf("invalid counter %s: %w", b[9:13], err)
	}
	c = c.SetCounter(int32(counter))

	fprint, err := strconv.ParseInt(string(b[13:17]), base, 64)
	if err != nil {
		return CUID{}, fmt.Errorf("invalid fingerprint %s: %w", b[13:17], err)
	}
	c = c.SetFingerprint(int32(fprint))

	rand1, err := strconv.ParseInt(string(b[17:21]), base, 64)
	if err != nil {
		return CUID{}, fmt.Errorf("invalid random %s: %w", b[17:21], err)
	}
	rand2, err := strconv.ParseInt(string(b[21:25]), base, 64)
	if err != nil {
		return CUID{}, fmt.Errorf("invalid random %s: %w", b[21:25], err)
	}
	c = c.SetRandom(int32(rand1), int32(rand2))
	return c, nil
}

// CUID is a 200 bit, or 25 byte, string value. These are defined by the
// https://github.com/ericelliott/cuid project. See
// https://github.com/ericelliott/cuid#motivation for details.
type CUID [25]byte

// Time returns the timestamp encoded in the CUID.
func (c CUID) Time() time.Time {
	msec, _ := binary.Varint(c[1:9])
	return time.Unix(0, msec*1e6)
}

// SetTime changes the timestamp of the CUID.
func (c CUID) SetTime(t time.Time) CUID {
	_ = binary.PutVarint(c[1:9], time.Duration(t.UnixNano()).Milliseconds())
	return c
}

// Counter returns the current sequence number.
func (c CUID) Counter() int32 {
	v, _ := binary.Varint(c[9:13])
	return int32(v)
}

// SetCounter changes the sequence number.
func (c CUID) SetCounter(v int32) CUID {
	_ = binary.PutVarint(c[9:13], int64(v))
	return c
}

// Fingerprint returns the ID value for the generator of the CUID.
func (c CUID) Fingerprint() int32 {
	v, _ := binary.Varint(c[13:17])
	return int32(v)
}

// SetFingerprint changes the generator ID.
func (c CUID) SetFingerprint(v int32) CUID {
	_ = binary.PutVarint(c[13:17], int64(v))
	return c
}

// Random returns the two random values encoded in the CUID.
func (c CUID) Random() (int32, int32) {
	rand1, _ := binary.Varint(c[17:21])
	rand2, _ := binary.Varint(c[21:25])
	return int32(rand1), int32(rand2)
}

// SetRandom changes the values of the random slots.
func (c CUID) SetRandom(first int32, second int32) CUID {
	_ = binary.PutVarint(c[17:21], int64(first))
	_ = binary.PutVarint(c[21:25], int64(second))
	return c
}

// String generates the canonical form of the CUID.
func (c CUID) String() string {
	nanos, _ := binary.Varint(c[1:9])
	counter, _ := binary.Varint(c[9:13])
	fingerprint, _ := binary.Varint(c[13:17])
	rand1, _ := binary.Varint(c[17:21])
	rand2, _ := binary.Varint(c[21:25])
	return prefix +
			leftpad(strconv.FormatInt(nanos, base), blockSize*2) +
			leftpad(strconv.FormatInt(counter, base), blockSize) +
			leftpad(strconv.FormatInt(fingerprint, base), blockSize) +
			leftpad(strconv.FormatInt(rand1, base), blockSize) +
			leftpad(strconv.FormatInt(rand2, base), blockSize)
}

// Slug generates a shortened version of the CUID that may be used as a
// disambiguation key in small documents or in URLs. Note, however, that this
// is a one-way process. Generating a slug is lossy such that the original
// CUID cannot be recreated.
func (c CUID) Slug() string {
	nanos, _ := binary.Varint(c[1:9])
	counter, _ := binary.Varint(c[9:13])
	fingerprint, _ := binary.Varint(c[13:17])
	rand1, _ := binary.Varint(c[17:21])

	nanosStr := strconv.FormatInt(nanos, base)
	nanosStr = nanosStr[len(nanosStr)-2:]
	counterStr := strconv.FormatInt(counter, base)
	if len(counterStr) > 4 {
		counterStr = counterStr[len(counterStr)-4:]
	}
	fingerprintStr := strconv.FormatInt(fingerprint, base)
	fingerprintStr = fingerprintStr[0:1] + fingerprintStr[len(fingerprintStr)-1:]
	randStr := strconv.FormatInt(rand1, base)
	if len(randStr) > 2 {
		randStr = randStr[len(randStr)-2:]
	}

	return (nanosStr + counterStr + fingerprintStr + randStr)
}

// Generator is a stateful producer of CUID values. This implements the logic
// defined by the JavaScript implementation found at
// https://github.com/ericelliott/cuid.
//
// You can create instances of this struct for cases where you want to inject
// CUID generation as a dependency rather than relying on the global convenience
// functions. All calls to the global functions use the default Generator
// instance which can be set with SetGenerator().
type Generator struct {
	Fingerprint int32
	Random      io.Reader
	Now         func() time.Time
	Counter     int32
	Locker      sync.Locker
}

// randomInt32 replaces the Rand.Int31N feature from math/rand that is lost when
// using crypto/rand. This function works by reading a random block of 4 bytes
// and then converting the sequence into an integer. Then the integer is
// modified to guarantee that it is less than or equal to the maxInt value in
// order to ensure the value can be encoded as base36. The reduction of the
// random integer uses an optimized form of modulo that is implemented with
// multiplication. See
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
// for details on the algorithm.
func (g *Generator) randomInt32() (int32, error) {
	b := make([]byte, 4)
	_, err := g.Random.Read(b)
	if err != nil {
		return 0, err
	}
	v, _ := binary.Varint(b)
	if v < 0 {
		v = -v
	}
	return int32((v * int64(maxInt)) >> 32), nil
}

// Generate a new CUID.
func (g *Generator) Generate() (CUID, error) {
	// NOTE: Due to the constrained space of the base36 encoding we must roll
	// over the integer before it would do so naturally. This requires
	// synchonrization across any concurrent usage so we must guard the
	// increment and roll over in a mutex.
	g.Locker.Lock()
	c := g.Counter
	g.Counter = g.Counter + 1
	if g.Counter > maxInt {
		g.Counter = 0
	}
	g.Locker.Unlock()

	rand1, err := g.randomInt32()
	if err != nil {
		return CUID{}, err
	}
	rand2, err := g.randomInt32()
	if err != nil {
		return CUID{}, err
	}
	v := (CUID{}).SetTime(g.Now()).SetCounter(c).SetFingerprint(g.Fingerprint).SetRandom(rand1, rand2)
	v[0] = prefixByte
	return v, err
}

// leftpad implements the behavior of
// https://github.com/ericelliott/cuid/blob/master/lib/pad.js. The purpose of
// this method is to ensure that all strings conform to some given size by
// padding the left-hand side with zeroes.
func leftpad(value string, size int) string {
	if len(value) == size {
		return value
	}
	if len(value) < size {
		value = padding + value
	}
	return value[len(value)-size:]
}

// defaultHostname matches the logic of
// https://github.com/ericelliott/cuid/blob/master/lib/fingerprint.js for
// converting a hostname into an int32 value. This is done by converting the
// hostname string into a byte sequence and then summing all the bytes into a
// single value.
func defaultHostname() int32 {
	h, err := os.Hostname()
	if err != nil {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		h = string(b)
	}
	hb := []byte(h)
	final := len(hb) + 36
	for _, b := range hb {
		final = final + int(b)
	}
	return int32(final)
}

func defaultPid() int32 {
	return int32(os.Getpid())
}

func defaultFingerprint() int32 {
	return defaultPid()<<2 | defaultHostname()>>2
}
