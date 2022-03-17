package guid

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

var (
	fp           = int32(-1)
	globalPrefix = [2]byte{'n', 'w'}
	prefixOnce   sync.Once

	// TestGUID is a nonsense GUID used for testing.
	// Its printable value is:
	// test0test0test0test0test00
	TestGUID = GUID{
		0x74, 0x65, // prefix
		0xa8, 0xd9, 0xac, 0xde, 0xb2, 0x83, 0x1, 0x0, // ts
		0xda, 0xc0, 0xa7, 0x1, // fp
		0xc8, 0xd3, 0x4, 0x0, // incr
		0xc4, 0xa5, 0xa5, 0x1, // decr
		0xa0, 0x87, 0xa4, 0x1, // rnd
	}
)

// SetGlobalPrefixBytes is a global initializer for GUID prefixes.
// The default prefix bytes are 'n' and 'w'.  This function can only
// be called once per execution.  Subsequent calls are no-ops.
// If one or both of the input bytes are not printable ascii characters
// in the base36 charset, this function will panic.
// !!! Panics triggered here must not be caught, or suppressed if they are caught.
func SetGlobalPrefixBytes(b1, b2 byte) {
	prefixOnce.Do(func() {
		if !(isValidPrefixByte(b1) && isValidPrefixByte(b2)) {
			panic("guid.SetGlobalPrefixBytes: prefix bytes must be base36-compatible and lowercase")
		}
		globalPrefix[0] = b1
		globalPrefix[1] = b2
	})
}

// filter out of band integers.  The integers produced by the
// default generator will always be within band
func filter(v int32) int64 {
	if v < 0 {
		v = -v
	}
	if v >= maxInt {
		return int64(v % maxInt)
	}
	return int64(v)
}

// generate a random 32 bit int.
// heavily influenced by fastrand in the runtime package of the standard library
// https://www.jstatsoft.org/article/view/v008i14/xorshift.pdf
func randomInt32(reader io.Reader) (int32, error) {
	b := [16]byte{}
	_, err := reader.Read(b[:])
	if err != nil {
		return 0, err
	}

	s0 := int32(b[0]+b[1]+b[2]+b[3]+b[4]+b[5]+b[6]+b[7]) | i32Buff
	s1 := int32(b[8]+b[9]+b[10]+b[11]+b[12]+b[13]+b[14]+b[15]) | i32Buff

	s1 ^= s1 << 17
	s1 = s1 ^ s0 ^ s1>>7 ^ s0>>16
	// we do need to do a modulo reduction here, but the cost is negligible
	return s0 + s1%maxInt, nil
}

// pad a byte slice with the zero value until
// it is the required size
func leftPad(in string, size int) string {
	l := len(in)
	// if the size of the input gte size or size
	// is invalid, just return the input
	if l >= size || size <= 0 {
		return in
	}
	b := append(bytes.Repeat([]byte{'0'}, size-l), []byte(in)...)

	return string(b)
}

// get the default hostname of the device
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

// get the current process id
func defaultPid() int32 {
	return int32(os.Getpid())
}

// get the default pid of the device
func defaultFingerprint() int32 {
	lfp := atomic.LoadInt32(&fp)
	if lfp == -1 {
		lfp = defaultPid()<<2 | defaultHostname()>>2
		atomic.StoreInt32(&fp, lfp)
	}

	return lfp
}

// prefix bytes must be printable base36 ASCII chars
func isValidPrefixByte(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'z')
}
