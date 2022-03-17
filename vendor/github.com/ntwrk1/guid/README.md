# GUID

**G**lobally **U**nique **Id**entifiers

This library generates globally unique 26-byte portable and serializable identifiers that are also highly collision-resistant and introspectable.

For a more thorough explanation, please see this Confluence doc: https://ntwrk.atlassian.net/wiki/spaces/ENG/pages/828112905/Globally+Unique+Identifiers

## Godoc
[here](./GODOC.md)

## Usage

### Default Behavior

#### Generating

The default signature for generating a GUID is `New() (GUID, error)`:

```go
1	package main
2
3	import (
4		"fmt"
5
6		"github.com/ntwrk1/guid"
7	)
8
9	func main() {
10		g, err := guid.New()
11		if err != nil {
12			panic(err)
13		}
14
15		fmt.Println(g)
16	}
```

The in the example above, the GUID generated would resemble something like:

`nwkkpkmep11m4l0000ykysk8xb`

This is the default behavior out of the box.  The odds of `guid.New` returning an error are extremely low and are bound to the particular implementation of `io.Reader` used when generating random bytes.  By default, the `crypto/rand` reader is used, so any errors returned would be due to a broken syscall.  In any event, the error should not be ignored.

#### Parsing

Parsing data into a GUID can be accomplished by calling either `Parse(b []byte)` or `ParseString(s string)`.  `ParseString` just calls `Parse` internally.

```go
1	package main
2
3	import (
4		"fmt"
5
6		"github.com/ntwrk1/guid"
7	)
8
9	func main() {
10		s := "nwkkpkmep11m4l0000ykysk8xb"
11		g, err := guid.ParseString(s)
12		if err != nil {
13			panic(err)
14		}
15
16		fmt.Println(g)
17	}
```

The above example would output `nwkkpkmep11m4l0000ykysk8xb`.

## Customization

### Prefix

By default, the first two prefix bytes of a GUID are `0x6e` and `0x77` (which serialize to `'n'` and `'w'`).  This behavior can be customized at application start up by calling the function `SetGlobalPrefixBytes(b1, b2 byte)`:

```go
1	package main
2
3	import (
4		"fmt"
5
6		"github.com/ntwrk1/guid"
7	)
8
9	func main() {
10		guid.SetGlobalPrefixBytes('4', '2')
11		g, err := guid.New()
12		if err != nil {
13			panic(err)
14		}
15
16		fmt.Println(g)
17	}
```

The above example would output something like `42kkpld1811pmd0000ylxbpvt6`.  Note that the prefix bytes should be printable ASCII characters, which can be conveniently wrapped with single quotes.  Calling `guid.SetGlobalPrefixBytes(4, 2)` will cause the application to panic.  Normal convention for a panicking library function is to prepend it with `Must...`, but this func does not, since it's internal to NTWRK.

Also note that subsequent calls to `guid.SetGlobalPrefixBytes` are no-ops.

#### Prefix per GUID

You can also easily set the prefix per GUID by either manually setting it:

```go
g := guid.New()
g[0], g[1] = '4', '2'
```

Or you can use the option `WithPrefixBytes`:

```go
g := guid.New(guid.WithPrefixBytes('4', '2'))
```

### BYO Generator

GUID uses a global instance of a `Generator` interface:

```go
// Generator defines the contract for generating GUIDs
type Generator interface {
    Generate() (GUID, error)
}
```

It is extremely unlikely that you'll ever need to use a custom generator.  If you do, just implement the `Generator` interface and call `SetGlobalGenerator`:

```go
1	package main
2
3	import (
4		"fmt"
5
6		"github.com/ntwrk1/guid"
7	)
8
9	type myGen struct{}
10
11	func (myGen) Generate() (guid.GUID, error) {
12		// hi000000010001000100010001
13		return guid.GUID{
14			'h', 'i', // prefix
15			0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ts
16			0x02, 0x00, 0x00, 0x00, // fp
17			0x02, 0x00, 0x00, 0x00, // incr counter
18			0x02, 0x00, 0x00, 0x00, // decr counter
19			0x02, 0x00, 0x00, 0x00, // rnd
20		}, nil
21	}
22
23	func main() {
24		guid.SetGlobalGenerator(myGen{})
25		g, err := guid.New()
26		if err != nil {
27			panic(err)
28		}
29
30		fmt.Println(g)
31	}
```

You can also just return `guid.TestGUID` from the generator if you like.  `guid.TestGUID` serializes to `test0test0test0test0test00`

## Introspection

I key feature of GUID is that it can be introspected.  The most common introspection will be to examine the prefix bytes, but every component of the GUID can be extracted:

```go
package main

import (
	"fmt"

	"github.com/ntwrk1/guid"
)

func main() {
	s := "nwkkpol4ar24bh0000yq3rus7i"
	g, err := guid.ParseString(s)
	if err != nil {
		panic(err)
	}

	// get prefix bytes
	b1, b2 := g.PrefixBytes()
	fmt.Printf("Prefix Bytes: %#x, %#x\n", b1, b2)

	// get timestamp
	t := g.Time()
	fmt.Printf("Time: %v\n", t)

	// get fingerprint
	fp := g.Fingerprint()
	fmt.Printf("Fingerprint: %#x\n", fp)

	// get counters
	incr, decr := g.Counters()
	fmt.Printf("Increment Counter: %d\n", incr)
	fmt.Printf("Decrement Counter: %d\n", decr)

	// get random noise
	rnd := g.Random()
	fmt.Printf("Random: %#x\n", rnd)
}
```

The above code will output the following:

```
Prefix Bytes: 0x6e, 0x77
Time: 2021-02-03 12:04:39.171 -0500 EST
Fingerprint: 0x1825d
Increment Counter: 0
Decrement Counter: 1620135
Random: 0x15ea4e
```
