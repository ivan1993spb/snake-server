package engine

import (
	"math/rand"
	"testing"
	"time"
)

func rawBenchmarkMapSet(b *testing.B, width, height uint8) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	a := MustArea(width, height)
	m := NewMap(a)
	object := NewObject("value")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dot := a.NewRandomDot(0, 0)
		b.StartTimer()

		m.Set(dot, object)
	}
}

func Benchmark_Map_Set_8x8(b *testing.B) {
	const (
		width  = 8
		height = 8
	)

	rawBenchmarkMapSet(b, width, height)
}

func Benchmark_Map_Set_128x128(b *testing.B) {
	const (
		width  = 128
		height = 128
	)

	rawBenchmarkMapSet(b, width, height)
}

func Benchmark_Map_Set_255x255(b *testing.B) {
	const (
		width  = 255
		height = 255
	)

	rawBenchmarkMapSet(b, width, height)
}

func rawBenchmarkMapGet(b *testing.B, width, height uint8) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	a := MustArea(width, height)
	m := NewMap(a)

	for y := uint8(0); y < a.Height(); y++ {
		for x := uint8(0); x < a.Width(); x++ {
			var dot = Dot{
				X: x,
				Y: y,
			}
			if dot.Hash()&1 == 1 {
				m.Set(dot, NewObject("value"))
			}
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dot := a.NewRandomDot(0, 0)
		b.StartTimer()

		m.Get(dot)
	}
}

func Benchmark_Map_Get_8x8(b *testing.B) {
	const (
		width  = 8
		height = 8
	)

	rawBenchmarkMapGet(b, width, height)
}

func Benchmark_Map_Get_128x128(b *testing.B) {
	const (
		width  = 128
		height = 128
	)

	rawBenchmarkMapGet(b, width, height)
}

func Benchmark_Map_Get_255x255(b *testing.B) {
	const (
		width  = 255
		height = 255
	)

	rawBenchmarkMapGet(b, width, height)
}

func rawBenchmarkMapMSetMRemove(b *testing.B, width, height uint8, dotsCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	a := MustArea(width, height)
	m := NewMap(a)
	object := NewObject("value")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dotIndex := rand.Int()

		dots := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(dotIndex+j) % a.Size()
			dots = append(dots, Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			})
		}

		b.StartTimer()

		m.MSet(dots, object)
		m.MRemove(dots)
	}
}

func Benchmark_Map_MSet_MRemove_64x64_d12(b *testing.B) {
	const (
		width  = 64
		height = 64

		dotsCount = 12
	)

	rawBenchmarkMapMSetMRemove(b, width, height, dotsCount)
}

func Benchmark_Map_MSet_MRemove_128x128_d32(b *testing.B) {
	const (
		width  = 128
		height = 128

		dotsCount = 32
	)

	rawBenchmarkMapMSetMRemove(b, width, height, dotsCount)
}

func Benchmark_Map_MSet_MRemove_255x255_d64(b *testing.B) {
	const (
		width  = 255
		height = 255

		dotsCount = 64
	)

	rawBenchmarkMapMSetMRemove(b, width, height, dotsCount)
}

func Benchmark_Map_MSet_MRemove_255x255_d256(b *testing.B) {
	const (
		width  = 255
		height = 255

		dotsCount = 256
	)

	rawBenchmarkMapMSetMRemove(b, width, height, dotsCount)
}

func rawBenchmarkMapMSetIfAbsentMRemoveObject(b *testing.B, width, height uint8, dotsCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	a := MustArea(width, height)
	m := NewMap(a)
	object := NewObject("value")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dotIndex := rand.Int()

		dots := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(dotIndex+j) % a.Size()
			dots = append(dots, Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			})
		}

		b.StartTimer()

		m.MSetIfAbsent(dots, object)
		m.MRemoveObject(dots, object)
	}
}

func Benchmark_Map_MSetIfAbsent_MRemoveObject_64x64_d12(b *testing.B) {
	const (
		width  = 64
		height = 64

		dotsCount = 12
	)

	rawBenchmarkMapMSetIfAbsentMRemoveObject(b, width, height, dotsCount)
}

func Benchmark_Map_MSetIfAbsent_MRemoveObject_128x128_d32(b *testing.B) {
	const (
		width  = 128
		height = 128

		dotsCount = 32
	)

	rawBenchmarkMapMSetIfAbsentMRemoveObject(b, width, height, dotsCount)
}

func Benchmark_Map_MSetIfAbsent_MRemoveObject_255x255_d64(b *testing.B) {
	const (
		width  = 255
		height = 255

		dotsCount = 64
	)

	rawBenchmarkMapMSetIfAbsentMRemoveObject(b, width, height, dotsCount)
}

func Benchmark_Map_MSetIfAbsent_MRemoveObject_255x255_d256(b *testing.B) {
	const (
		width  = 255
		height = 255

		dotsCount = 256
	)

	rawBenchmarkMapMSetIfAbsentMRemoveObject(b, width, height, dotsCount)
}

func rawBenchmarkMapSetMSetIfAllAbsentRemoveObject(b *testing.B, width, height uint8, dotsCount int) {
	b.ReportAllocs()

	rand.Seed(time.Now().UTC().UnixNano())

	a := MustArea(width, height)
	m := NewMap(a)
	object := NewObject("value")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dotIndex := rand.Int()

		dots := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(dotIndex+j) % a.Size()
			dots = append(dots, Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			})
		}

		dotSpoiler := dots[dotsCount-1]

		b.StartTimer()

		m.Set(dotSpoiler, object)
		m.MSetIfAllAbsent(dots, object)
		m.RemoveObject(dotSpoiler, object)
	}
}

func Benchmark_Map_Set_MSetIfAllAbsent_RemoveObject_64x64_d12(b *testing.B) {
	const (
		width  = 64
		height = 64

		dotsCount = 12
	)

	rawBenchmarkMapSetMSetIfAllAbsentRemoveObject(b, width, height, dotsCount)
}

func Benchmark_Map_Set_MSetIfAllAbsent_RemoveObject_128x128_d32(b *testing.B) {
	const (
		width  = 128
		height = 128

		dotsCount = 32
	)

	rawBenchmarkMapSetMSetIfAllAbsentRemoveObject(b, width, height, dotsCount)
}

func Benchmark_Map_Set_MSetIfAllAbsent_RemoveObject_255x255_d64(b *testing.B) {
	const (
		width  = 255
		height = 255

		dotsCount = 64
	)

	rawBenchmarkMapSetMSetIfAllAbsentRemoveObject(b, width, height, dotsCount)
}

func Benchmark_Map_Set_MSetIfAllAbsent_RemoveObject_255x255_d256(b *testing.B) {
	const (
		width  = 255
		height = 255

		dotsCount = 256
	)

	rawBenchmarkMapSetMSetIfAllAbsentRemoveObject(b, width, height, dotsCount)
}
