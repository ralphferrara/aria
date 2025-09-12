package crypto

import (
	"crypto/sha256"
	"io"
)

//||------------------------------------------------------------------------------------------------||
//|| NewDeterministicReader
//||------------------------------------------------------------------------------------------------||

type deterministicReader struct {
	seed   []byte
	buffer []byte
	index  int
}

//||------------------------------------------------------------------------------------------------||
//|| Read generates pseudo-random bytes using repeated SHA256 hashing of the seed
//||------------------------------------------------------------------------------------------------||

func (d *deterministicReader) Read(p []byte) (int, error) {
	for len(d.buffer) < len(p) {
		hash := sha256.Sum256(append(d.seed, byte(d.index)))
		d.buffer = append(d.buffer, hash[:]...)
		d.index++
	}
	n := copy(p, d.buffer[:len(p)])
	d.buffer = d.buffer[len(p):]
	return n, nil
}

//||------------------------------------------------------------------------------------------------||
//|| NewDeterministicReader returns an io.Reader that produces deterministic bytes from a seed
//||------------------------------------------------------------------------------------------------||

func NewDeterministicReader(seed []byte) io.Reader {
	seedCopy := make([]byte, len(seed))
	copy(seedCopy, seed)
	return &deterministicReader{
		seed:   seedCopy,
		buffer: make([]byte, 0, 64),
		index:  0,
	}
}
