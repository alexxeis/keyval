package cluster

import (
	"hash"
)

type djb2aStringHash32 uint32

// newDjb32a returns a new hash.Hash32 object, computing a variant of Daniel J. Bernstein's hash that uses xor instead of +
func newDjb32a() hash.Hash32 {
	sh := djb2aStringHash32(0)
	sh.Reset()
	return &sh
}

func (sh *djb2aStringHash32) Size() int {
	return 4
}

func (sh *djb2aStringHash32) BlockSize() int {
	return 1
}

func (sh *djb2aStringHash32) Sum32() uint32 {
	return uint32(*sh)
}

func (sh *djb2aStringHash32) Reset() {
	*sh = djb2aStringHash32(5381)
}

func (sh *djb2aStringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *djb2aStringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = 33*h ^ uint32(c)
	}
	*sh = djb2aStringHash32(h)
	return len(b), nil
}
