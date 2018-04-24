package bufferpool

import (
	"sync"
)

// BufferPool provides a pool of byte slice, made from a single huge byte array.
type BufferPool struct {
	mx                      sync.RWMutex
	bufferCount, bufferSize int
	pool                    [][]byte
}

// New makes a new *BufferPool
func New(bufferSize, bufferCount int) *BufferPool {
	pool, _ := newBufferPool(bufferSize, bufferCount)
	return pool
}

func newBufferPool(bufferSize, bufferCount int) (pool *BufferPool, back []byte) {
	var res BufferPool
	res.bufferSize = bufferSize
	res.bufferCount = bufferCount
	res.pool, back = makePartitions(bufferSize, bufferCount)
	pool = &res
	return
}

func makePartitions(bufferSize, bufferCount int) (pool [][]byte, back []byte) {
	size := bufferCount * bufferSize
	back = make([]byte, size, size)
	for i := 0; i < bufferCount; i++ {
		low := i * bufferSize
		high := low + bufferSize
		pool = append(pool, back[low:high:high])
	}
	return
}

// Len returns the length of the pool.
func (bf *BufferPool) Len() int {
	bf.mx.RLock()
	defer bf.mx.RUnlock()
	l := len(bf.pool)
	return l
}

// Take returns a byte slice from the pool or nil if the pool is depleted.
func (bf *BufferPool) Take() (buffer []byte) {
	bf.mx.Lock()
	defer bf.mx.Unlock()
	l := len(bf.pool)
	if l == 0 {
		return nil
	}
	buffer, bf.pool = bf.pool[l-1], bf.pool[:l-1]
	return buffer
}

// Put puts back a []byte into the pool unless the pool if full or
// the provided buffer has a different length than the initial buffer size.
func (bf *BufferPool) Put(buffer []byte) bool {
	bf.mx.Lock()
	defer bf.mx.Unlock()
	l := len(bf.pool)
	if l == bf.bufferCount {
		return false
	}
	if len(buffer) != bf.bufferSize {
		return false
	}
	bf.pool = append(bf.pool, buffer)
	return true
}

// Expand expands the pool bufferCount times, creating a new underlying array.
func (bf *BufferPool) Expand(bufferCount int) {
	bf.mx.Lock()
	defer bf.mx.Unlock()
	pool, _ := makePartitions(bf.bufferSize, bufferCount)
	bf.pool = append(bf.pool, pool...)
	bf.bufferCount += bufferCount
}
