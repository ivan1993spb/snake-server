package cmap

import "testing"
import "strconv"

const keyA = 0

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
