package glowlib

import (
	"regexp"
	"strconv"

	"github.com/jonnyhopper/murmur3"
)

// this must be consistent with client
const cMurmurHashSeed = 0xa9401e9f

var cHashPattern = "^([a-fA-F0-9]{32})$"
var cHashValidator = regexp.MustCompile(cHashPattern)

// when uploading a resource, the bare hash name is all that is supplied. One should never be uploading a dev resource!
func HashIsValid(hashname string) bool {
	return len(hashname) == 32 && cHashValidator.MatchString(hashname)
}

// Hash128 :we really want them to be uint64 but it's more convenient for MongoDB, which doesn't support uint64
type Hash128 struct {
	H1 int64
	H2 int64
}

// IsValid : if both H1 and H2 are non-zero
func (h Hash128) IsValid() bool {
	return h.H1 != 0 && h.H2 != 0
}

// convert a 128 bit hash to a 32 character string
// the swizzling is so that the byte order matches
func (h Hash128) String() string {
	h1 := h.H1
	h2 := h.H2
	h1a := ((h1 >> 56) & 0xff) |
		((h1 & 0x00ff000000000000) >> 40) |
		((h1 & 0x0000ff0000000000) >> 24) |
		((h1 & 0x000000ff00000000) >> 8) |
		((h1 & 0x00000000ff000000) << 8) |
		((h1 & 0x0000000000ff0000) << 24) |
		((h1 & 0x000000000000ff00) << 40) |
		((h1 & 0x00000000000000ff) << 56)

	h2a := ((h2 >> 56) & 0xff) |
		((h2 & 0x00ff000000000000) >> 40) |
		((h2 & 0x0000ff0000000000) >> 24) |
		((h2 & 0x000000ff00000000) >> 8) |
		((h2 & 0x00000000ff000000) << 8) |
		((h2 & 0x0000000000ff0000) << 24) |
		((h2 & 0x000000000000ff00) << 40) |
		((h2 & 0x00000000000000ff) << 56)

	sh1 := strconv.FormatUint(uint64(h1a), 16)
	sh2 := strconv.FormatUint(uint64(h2a), 16)
	sh1Pad := 16 - len(sh1)
	sh2Pad := 16 - len(sh2)

	// this method is actually faster (!!) than using strings.Repeat() or fmt.Printf
	for j := 0; j < sh1Pad; j++ {
		sh1 = "0" + sh1
	}
	for j := 0; j < sh2Pad; j++ {
		sh2 = "0" + sh2
	}

	return sh1 + sh2
}

// Hash128FromHex converts a 32 character hex string to two uints, which is a hash 128
// the swizzling is so that the byte order matches
func Hash128FromHex(hashHex string) Hash128 {
	h1, _ := strconv.ParseUint(hashHex[:16], 16, 64)
	h2, _ := strconv.ParseUint(hashHex[16:], 16, 64)

	h1a := ((h1 >> 56) & 0xff) |
		((h1 & 0x00ff000000000000) >> 40) |
		((h1 & 0x0000ff0000000000) >> 24) |
		((h1 & 0x000000ff00000000) >> 8) |
		((h1 & 0x00000000ff000000) << 8) |
		((h1 & 0x0000000000ff0000) << 24) |
		((h1 & 0x000000000000ff00) << 40) |
		((h1 & 0x00000000000000ff) << 56)

	h2a := ((h2 >> 56) & 0xff) |
		((h2 & 0x00ff000000000000) >> 40) |
		((h2 & 0x0000ff0000000000) >> 24) |
		((h2 & 0x000000ff00000000) >> 8) |
		((h2 & 0x00000000ff000000) << 8) |
		((h2 & 0x0000000000ff0000) << 24) |
		((h2 & 0x000000000000ff00) << 40) |
		((h2 & 0x00000000000000ff) << 56)

	return Hash128{int64(h1a), int64(h2a)}
}

// HashData takes a byte stream
func HashData(data []byte) Hash128 {
	h1, h2 := murmur3.Sum128WithSeed(data, cMurmurHashSeed)
	return Hash128{int64(h1), int64(h2)}
}
