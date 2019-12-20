package cmap

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const keyA = 0

func getHash(x, y uint8) uint16 {
	return uint16(x)<<8 | uint16(y)
}

func BenchmarkItems(b *testing.B) {
	m, _ := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Items()
	}
}

func BenchmarkMarshalJson(b *testing.B) {
	m, _ := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.MarshalJSON()
	}
}

func BenchmarkStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}

func BenchmarkSingleInsertAbsent(b *testing.B) {
	m, _ := New(defaultShardCount)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(uint16(i), "value")
	}
}

func rawBenchmarkConcurrentMapSet(b *testing.B, shardCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	m, _ := New(shardCount)
	object := "value"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		hash := uint16(rand.Int())
		b.StartTimer()

		m.Set(hash, object)
	}
}

func Benchmark_ConcurrentMap_Set_sc2(b *testing.B) {
	const shardCount = 2

	rawBenchmarkConcurrentMapSet(b, shardCount)
}

func Benchmark_ConcurrentMap_Set_sc32(b *testing.B) {
	const shardCount = 32

	rawBenchmarkConcurrentMapSet(b, shardCount)
}

func Benchmark_ConcurrentMap_Set_sc64(b *testing.B) {
	const shardCount = 64

	rawBenchmarkConcurrentMapSet(b, shardCount)
}

func rawBenchmarkConcurrentMapGet(b *testing.B, width, height uint8, shardCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	m, _ := New(shardCount)

	for y := uint8(0); y < height; y++ {
		for x := uint8(0); x < width; x++ {
			hash := getHash(x, y)

			if hash&1 == 1 {
				m.Set(hash, "object")
			}
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		x := uint8(rand.Intn(int(width)))
		y := uint8(rand.Intn(int(height)))
		hash := getHash(x, y)
		b.StartTimer()

		m.Get(hash)
	}
}

func Benchmark_ConcurrentMap_Get_8x8_sc2(b *testing.B) {
	const (
		width  = 8
		height = 8

		shardCount = 2
	)

	rawBenchmarkConcurrentMapGet(b, width, height, shardCount)
}

func Benchmark_ConcurrentMap_Get_128x128_sc32(b *testing.B) {
	const (
		width  = 128
		height = 128

		shardCount = 32
	)

	rawBenchmarkConcurrentMapGet(b, width, height, shardCount)
}

func Benchmark_ConcurrentMap_Get_255x255_sc64(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64
	)

	rawBenchmarkConcurrentMapGet(b, width, height, shardCount)
}

func rawBenchmarkConcurrentMapMSetMRemove(b *testing.B, width, height uint8, shardCount, dotsCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	size := uint16(width) * uint16(height)
	m, _ := New(shardCount)

	object := "value"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dotIndex := rand.Int()

		dotsToBeRemoved := make([]uint16, 0, dotsCount)
		dotsToBeSet := make(map[uint16]interface{})
		for j := 0; j < dotsCount; j++ {
			index := uint16(dotIndex+j) % size
			hash := getHash(uint8(index%uint16(width)), uint8(index/uint16(width)))
			dotsToBeRemoved = append(dotsToBeRemoved, hash)
			dotsToBeSet[hash] = object
		}

		b.StartTimer()

		m.MSet(dotsToBeSet)
		m.MRemove(dotsToBeRemoved)
	}
}

func Benchmark_ConcurrentMap_MSet_MRemove_64x64_sc16_d12(b *testing.B) {
	const (
		width  = 64
		height = 64

		shardCount = 16

		dotsCount = 12
	)

	rawBenchmarkConcurrentMapMSetMRemove(b, width, height, shardCount, dotsCount)
}

func Benchmark_ConcurrentMap_MSet_MRemove_128x128_sc32_d32(b *testing.B) {
	const (
		width  = 128
		height = 128

		shardCount = 32

		dotsCount = 32
	)

	rawBenchmarkConcurrentMapMSetMRemove(b, width, height, shardCount, dotsCount)
}

func Benchmark_ConcurrentMap_MSet_MRemove_255x255_sc64_d64(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64

		dotsCount = 64
	)

	rawBenchmarkConcurrentMapMSetMRemove(b, width, height, shardCount, dotsCount)
}

func Benchmark_ConcurrentMap_MSet_MRemove_255x255_sc64_d256(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64

		dotsCount = 256
	)

	rawBenchmarkConcurrentMapMSetMRemove(b, width, height, shardCount, dotsCount)
}

func rawBenchmarkConcurrentMapMSetIfAbsentMRemoveCb(b *testing.B, width, height uint8, shardCount, dotsCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	size := uint16(width) * uint16(height)
	m, _ := New(shardCount)

	object := "value"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dotIndex := rand.Int()

		dotsToBeRemoved := make([]uint16, 0, dotsCount)
		dotsToBeSet := make(map[uint16]interface{})

		for j := 0; j < dotsCount; j++ {
			index := uint16(dotIndex+j) % size
			hash := getHash(uint8(index%uint16(width)), uint8(index/uint16(width)))
			dotsToBeRemoved = append(dotsToBeRemoved, hash)
			dotsToBeSet[hash] = object
		}

		b.StartTimer()

		m.MSetIfAbsent(dotsToBeSet)
		m.MRemoveCb(dotsToBeRemoved, func(key uint16, v interface{}, exists bool) bool {
			return exists && v == object
		})
	}
}

func Benchmark_ConcurrentMap_MSetIfAbsent_MRemoveCb_64x64_sc16_d12(b *testing.B) {
	const (
		width  = 64
		height = 64

		shardCount = 16

		dotsCount = 12
	)

	rawBenchmarkConcurrentMapMSetIfAbsentMRemoveCb(b, width, height, dotsCount, shardCount)
}

func Benchmark_ConcurrentMap_MSetIfAbsent_MRemoveCb_128x128_sc32_d32(b *testing.B) {
	const (
		width  = 128
		height = 128

		shardCount = 32

		dotsCount = 32
	)

	rawBenchmarkConcurrentMapMSetIfAbsentMRemoveCb(b, width, height, dotsCount, shardCount)
}

func Benchmark_ConcurrentMap_MSetIfAbsent_MRemoveCb_255x255_sc64_d64(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64

		dotsCount = 64
	)

	rawBenchmarkConcurrentMapMSetIfAbsentMRemoveCb(b, width, height, dotsCount, shardCount)
}

func Benchmark_ConcurrentMap_MSetIfAbsent_MRemoveCb_255x255_sc64_d256(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64

		dotsCount = 256
	)

	rawBenchmarkConcurrentMapMSetIfAbsentMRemoveCb(b, width, height, dotsCount, shardCount)
}

func rawBenchmarkConcurrentMapSetMSetIfAllAbsentRemoveCb(b *testing.B, width, height uint8, shardCount, dotsCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	size := uint16(width) * uint16(height)
	m, _ := New(shardCount)

	object := "value"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dotIndex := rand.Int()

		dotsToBeRemoved := make([]uint16, 0, dotsCount)
		dotsToBeSet := make(map[uint16]interface{})
		for j := 0; j < dotsCount; j++ {
			index := uint16(dotIndex+j) % size
			hash := getHash(uint8(index%uint16(width)), uint8(index/uint16(width)))
			dotsToBeRemoved = append(dotsToBeRemoved, hash)
			dotsToBeSet[hash] = object
		}

		spoiler := dotsToBeRemoved[dotsCount-1]

		b.StartTimer()

		m.Set(spoiler, object)
		m.MSetIfAllAbsent(dotsToBeSet)
		m.MRemoveCb(dotsToBeRemoved, func(key uint16, v interface{}, exists bool) bool {
			return exists && v == object
		})
	}
}

func Benchmark_ConcurrentMap_Set_MSetIfAllAbsent_RemoveCb_64x64_sc16_d12(b *testing.B) {
	const (
		width  = 64
		height = 64

		shardCount = 16

		dotsCount = 12
	)

	rawBenchmarkConcurrentMapSetMSetIfAllAbsentRemoveCb(b, width, height, shardCount, dotsCount)
}

func Benchmark_ConcurrentMap_Set_MSetIfAllAbsent_RemoveCb_128x128_sc32_d32(b *testing.B) {
	const (
		width  = 128
		height = 128

		shardCount = 32

		dotsCount = 32
	)

	rawBenchmarkConcurrentMapSetMSetIfAllAbsentRemoveCb(b, width, height, shardCount, dotsCount)
}

func Benchmark_ConcurrentMap_Set_MSetIfAllAbsent_RemoveCb_255x255_sc64_d64(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64

		dotsCount = 64
	)

	rawBenchmarkConcurrentMapSetMSetIfAllAbsentRemoveCb(b, width, height, shardCount, dotsCount)
}

func Benchmark_ConcurrentMap_Set_MSetIfAllAbsent_RemoveCb_255x255_sc64_d256(b *testing.B) {
	const (
		width  = 255
		height = 255

		shardCount = 64

		dotsCount = 256
	)

	rawBenchmarkConcurrentMapSetMSetIfAllAbsentRemoveCb(b, width, height, shardCount, dotsCount)
}

func BenchmarkSingleInsertPresent(b *testing.B) {
	m, _ := New(defaultShardCount)
	m.Set(keyA, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(keyA, "value")
	}
}

func benchmarkMultiInsertDifferent(b *testing.B, shardCount int) {
	m, _ := New(shardCount)
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(uint16(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_1_Shard(b *testing.B) {
	benchmarkMultiInsertDifferent(b, 1)
}
func BenchmarkMultiInsertDifferent_16_Shard(b *testing.B) {
	benchmarkMultiInsertDifferent(b, 16)
}
func BenchmarkMultiInsertDifferent_32_Shard(b *testing.B) {
	benchmarkMultiInsertDifferent(b, 32)
}
func BenchmarkMultiInsertDifferent_256_Shard(b *testing.B) {
	benchmarkMultiInsertDifferent(b, 256)
}

func BenchmarkMultiInsertSame(b *testing.B) {
	m, _ := New(defaultShardCount)
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	m.Set(keyA, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(keyA, "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSame(b *testing.B) {
	m, _ := New(defaultShardCount)
	finished := make(chan struct{}, b.N)
	get, _ := GetSet(m, finished)
	m.Set(keyA, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get(keyA, "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func benchmarkMultiGetSetDifferent(b *testing.B, shardCount int) {
	m, _ := New(shardCount)
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set(0, "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(uint16(i), "value")
		get(uint16(i+1), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_1_Shard(b *testing.B) {
	benchmarkMultiGetSetDifferent(b, 1)
}
func BenchmarkMultiGetSetDifferent_16_Shard(b *testing.B) {
	benchmarkMultiGetSetDifferent(b, 16)
}
func BenchmarkMultiGetSetDifferent_32_Shard(b *testing.B) {
	benchmarkMultiGetSetDifferent(b, 32)
}
func BenchmarkMultiGetSetDifferent_256_Shard(b *testing.B) {
	benchmarkMultiGetSetDifferent(b, 256)
}

func benchmarkMultiGetSetBlock(b *testing.B, shardCount int) {
	m, _ := New(shardCount)
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(uint16(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(uint16(i%100), "value")
		get(uint16(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_1_Shard(b *testing.B) {
	benchmarkMultiGetSetBlock(b, 1)
}
func BenchmarkMultiGetSetBlock_16_Shard(b *testing.B) {
	benchmarkMultiGetSetBlock(b, 16)
}
func BenchmarkMultiGetSetBlock_32_Shard(b *testing.B) {
	benchmarkMultiGetSetBlock(b, 32)
}
func BenchmarkMultiGetSetBlock_256_Shard(b *testing.B) {
	benchmarkMultiGetSetBlock(b, 256)
}

func GetSet(m *ConcurrentMap, finished chan struct{}) (set func(key uint16, value string), get func(key uint16, value string)) {
	return func(key uint16, value string) {
			for i := 0; i < 10; i++ {
				m.Get(key)
			}
			finished <- struct{}{}
		}, func(key uint16, value string) {
			for i := 0; i < 10; i++ {
				m.Set(key, value)
			}
			finished <- struct{}{}
		}
}

func BenchmarkKeys(b *testing.B) {
	m, _ := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Keys()
	}
}
