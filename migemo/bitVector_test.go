package migemo_test

import (
	"math/rand"
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func bits_to_vector(bits []bool) []uint32 {
	vec := make([]uint32, (len(bits)+63)>>5)
	for i := 0; i < len(bits); i++ {
		if bits[i] {
			vec[i>>5] |= uint32(1) << (i & 31)
		}
	}
	return vec
}

func rank(bits []bool, position uint, b bool) uint {
	count := uint(0)
	for i := uint(0); i < position; i++ {
		if bits[i] == b {
			count++
		}
	}
	return count
}

func TestRank(t *testing.T) {
	size := uint(1000)
	bits := make([]bool, size)
	for i := uint(0); i < size; i++ {
		bits[i] = rand.Intn(2) == 1
	}
	vec := bits_to_vector(bits)
	bv := migemo.NewBitVector(vec, uint32(size))
	for i := uint(0); i < size; i++ {
		expected := rank(bits, i, true)
		actual := bv.Rank(uint32(i), true)
		if expected != uint(actual) {
			t.Error("actual: ", actual, "\nexpected: ", expected, "\n")
		}
	}
	for i := uint(0); i < size; i++ {
		expected := rank(bits, i, false)
		actual := bv.Rank(uint32(i), false)
		if expected != uint(actual) {
			t.Error("actual: ", actual, "\nexpected: ", expected, "\n")
		}
	}
}

func _select(bits []bool, count uint, b bool) int {
	for i := 0; i < len(bits); i++ {
		if bits[i] == b {
			count--
		}
		if count == 0 {
			return i
		}
	}
	return -1
}

func TestSelect(t *testing.T) {
	size := uint(1000)
	bits := make([]bool, size)
	for i := uint(0); i < size; i++ {
		bits[i] = rand.Intn(2) == 1
	}
	vec := bits_to_vector(bits)
	bv := migemo.NewBitVector(vec, uint32(size))
	count1 := uint(0)
	for i := 0; i < len(bits); i++ {
		if bits[i] {
			count1++
		}
	}
	count0 := uint(len(bits)) - count1
	for i := uint(1); i < count1; i++ {
		expected := _select(bits, i, true)
		actual := bv.Select(uint32(i), true)
		if expected != int(actual) {
			t.Error("actual: ", actual, "\nexpected: ", expected, "\n")
		}
	}
	for i := uint(1); i < count0; i++ {
		expected := _select(bits, i, false)
		actual := bv.Select(uint32(i), false)
		if expected != int(actual) {
			t.Error("actual: ", actual, "\nexpected: ", expected, "\n")
		}
	}
}

func TestNextClearBit(t *testing.T) {
	size := uint(100)
	bits := make([]bool, size)
	for i := uint(0); i < size; i++ {
		bits[i] = true // rand.Intn(2) == 1
	}
	bits[50] = false
	vec := bits_to_vector(bits)
	bv := migemo.NewBitVector(vec, uint32(size))
	if bv.NextClearBit(uint32(34)) != 50 {
		t.Error()
	}
}
