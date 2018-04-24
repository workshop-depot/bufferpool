package bufferpool

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makePartitions_1(t *testing.T) {
	assert := assert.New(t)

	pool, back := makePartitions(10, 10)
	assert.Equal(100, len(back))
	assert.Equal(10, len(pool))
	for _, v := range pool {
		assert.Equal(10, len(v))
	}

	copy(back, []byte(strings.Repeat("-", 200)))
	assert.Equal(100, len(back))
	for _, v := range back {
		assert.Equal("-"[0], v)
	}

	for _, vp := range pool {
		for _, v := range vp {
			assert.Equal("-"[0], v)
		}
	}

	for kp := range pool {
		pool[kp] = append(pool[kp][:0], []byte(strings.Repeat("*", 10))...)
	}

	for _, v := range back {
		assert.Equal("*"[0], v)
	}

	for kp := range pool {
		pool[kp] = append(pool[kp][:0], []byte(strings.Repeat("-", 11))...)
	}

	for _, v := range back {
		assert.Equal("*"[0], v)
	}
}

func Test_makePartitions_2(t *testing.T) {
	assert := assert.New(t)

	pool, back := makePartitions(10, 10)
	assert.Equal(100, len(back))
	assert.Equal(10, len(pool))
	for _, v := range pool {
		assert.Equal(10, len(v))
	}

	copy(back, []byte(strings.Repeat("-", 200)))
	assert.Equal(100, len(back))
	for _, v := range back {
		assert.Equal("-"[0], v)
	}

	for _, vp := range pool {
		for _, v := range vp {
			assert.Equal("-"[0], v)
		}
	}

	for kp := range pool {
		copy(pool[kp], []byte(strings.Repeat(string(kp), 200)))
	}
	for _, v := range pool {
		assert.Equal(10, len(v))
	}

	for k, v := range pool {
		check := []byte(strings.Repeat(string(k), 10))
		assert.Equal(0, bytes.Compare(v, check))
	}
	for k, v := range back {
		assert.Equal([]byte(string(k / 10))[0], v)
	}

	for _, vp := range pool {
		copy(vp, []byte(strings.Repeat("-", 10)))
	}
	for _, v := range back {
		assert.Equal("-"[0], v)
	}
}

func Test_Smoke_BufferPool(t *testing.T) {
	assert := assert.New(t)

	pool := New(10, 10)
	b := pool.Take()
	assert.NotNil(b)
	assert.Equal(9, pool.Len())
	assert.True(pool.Put(b))
	assert.Equal(10, pool.Len())
	pool.Expand(10)
	assert.Equal(20, pool.Len())

	var (
		next = true
		pile [][]byte
	)
	for next {
		b := pool.Take()
		if b == nil {
			next = false
			continue
		}
		pile = append(pile, b)
	}

	big := make([]byte, 100)
	assert.False(pool.Put(big))

	for _, v := range pile {
		pool.Put(v)
	}

	big = make([]byte, 10)
	assert.False(pool.Put(big))

	pool.Take()

	big = make([]byte, 100)
	assert.False(pool.Put(big))
}
