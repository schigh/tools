# guid


## Usage

```go
var (

	// TestGUID is a nonsense GUID used for testing.
	// Its printable value is:
	// test0test0test0test0test00
	TestGUID = GUID{
		0x74, 0x65,
		0xa8, 0xd9, 0xac, 0xde, 0xb2, 0x83, 0x1, 0x0,
		0xda, 0xc0, 0xa7, 0x1,
		0xc8, 0xd3, 0x4, 0x0,
		0xc4, 0xa5, 0xa5, 0x1,
		0xa0, 0x87, 0xa4, 0x1,
	}
)
```

#### func  SetGlobalGenerator

```go
func SetGlobalGenerator(g Generator)
```
SetGlobalGenerator allows for the manual assignment of the GUID generator. The
main usefulness of this function is primarily for testing, but this function can
also be used to inject custom time and randomness providers. Note that this
function can be called only once per runtime. Subsequent calls are no-ops.

#### func  SetGlobalPrefixBytes

```go
func SetGlobalPrefixBytes(b1, b2 byte)
```
SetGlobalPrefixBytes is a global initializer for GUID prefixes. The default
prefix bytes are 'n' and 'w'. This function can only be called once per
execution. Subsequent calls are no-ops. If one or both of the input bytes are
not printable ascii characters in the base36 charset, this function will panic.
!!! Panics triggered here must not be caught, or suppressed if they are caught.

#### type GUID

```go
type GUID [byteSize]byte
```

GUID is a globally unique identifier

    prefix  timestamp                 fingerprint   incr          decr          random

[[b, b], [b, b, b, b, b, b, b, b], [b, b, b, b], [b, b, b, b], [b, b, b, b], [b,
b, b, b]]

#### func  New

```go
func New(opts ...Option) GUID
```
New will return only a GUID using the global generator or panic (less safe way,
but unlikely to fail)

#### func  NewRandom

```go
func NewRandom(opts ...Option) (GUID, error)
```
NewRandom will create a GUID using the global generator or return an error (this
is the safer way)

#### func  Parse

```go
func Parse(in []byte) (GUID, error)
```
Parse the byte slice into a guid

#### func  ParseString

```go
func ParseString(s string) (GUID, error)
```
ParseString is a convenience func for parsing GUID strings

#### func (GUID) Counters

```go
func (g GUID) Counters() (int32, int32)
```
Counters returns the incrementer and decrementer counters.

#### func (GUID) Fingerprint

```go
func (g GUID) Fingerprint() int32
```
Fingerprint returns the device fingerprint

#### func (*GUID) GobDecode

```go
func (g *GUID) GobDecode(data []byte) error
```
GobDecode implements gob.GobDecoder

#### func (GUID) GobEncode

```go
func (g GUID) GobEncode() ([]byte, error)
```
GobEncode implements gob.GobEncoder

#### func (GUID) MarshalBinary

```go
func (g GUID) MarshalBinary() (data []byte, err error)
```
MarshalBinary implements encoding.BinaryMarshaler

#### func (GUID) MarshalJSON

```go
func (g GUID) MarshalJSON() ([]byte, error)
```
MarshalJSON implements json.Marshaler

#### func (GUID) MarshalText

```go
func (g GUID) MarshalText() (text []byte, err error)
```
MarshalText implements encoding.TextMarshaler

#### func (GUID) PrefixBytes

```go
func (g GUID) PrefixBytes() (byte, byte)
```
PrefixBytes returns the two GUID prefix bytes individually

#### func (GUID) Random

```go
func (g GUID) Random() int32
```
Random returns the random value encoded in the GUID.

#### func (*GUID) Scan

```go
func (g *GUID) Scan(v interface{}) error
```
Scan implements sql.Scanner

#### func (GUID) SetCounters

```go
func (g GUID) SetCounters(incr, decr int32) GUID
```
SetCounters sets the increment and decrement counters

#### func (GUID) SetFingerprint

```go
func (g GUID) SetFingerprint(v int32) GUID
```
SetFingerprint adds the device fingerprint Glyph to the GUID

#### func (GUID) SetRandom

```go
func (g GUID) SetRandom(v int32) GUID
```
SetRandom adds an 8-byte random Word to the GUID

#### func (GUID) SetTime

```go
func (g GUID) SetTime(t time.Time) GUID
```
SetTime inserts the unix timestamp into the GUID

#### func (GUID) Slug

```go
func (g GUID) Slug() string
```
Slug returns a shortened version of the GUID that may be used as a
disambiguation key in small documents or in URLs. Note that this is a ONE WAY
PROCESS. Generating a slug is lossy such that the original GUID cannot be
recreated.

#### func (GUID) String

```go
func (g GUID) String() string
```

#### func (GUID) Time

```go
func (g GUID) Time() time.Time
```
Time returns the timestamp embedded in the GUID

#### func (*GUID) UnmarshalBinary

```go
func (g *GUID) UnmarshalBinary(data []byte) error
```
UnmarshalBinary implements encoding.BinaryUnmarshaler

#### func (*GUID) UnmarshalJSON

```go
func (g *GUID) UnmarshalJSON(b []byte) error
```
UnmarshalJSON implements json.Unmarshaler

#### func (*GUID) UnmarshalText

```go
func (g *GUID) UnmarshalText(text []byte) error
```
UnmarshalText implements encoding.TextUnmarshaler

#### func (GUID) Value

```go
func (g GUID) Value() (driver.Value, error)
```
Value implements driver.Valuer

#### type Generator

```go
type Generator interface {
	Generate() (GUID, error)
}
```

Generator defines the contract for generating GUIDs

#### type Option

```go
type Option func(*GUID)
```

Option is a function that allows for mutating a GUID

#### func  WithPrefixBytes

```go
func WithPrefixBytes(b1, b2 byte) Option
```
WithPrefixBytes sets the prefix bytes for a single GUID. To set GUID prefix
bytes globally, use the SetGlobalPrefixBytes function.
