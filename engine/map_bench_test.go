package engine

import "testing"

func Benchmark_Map_Set(b *testing.B) {
	a := MustArea(255, 255)
	m := NewMap(a)
	object := NewObject("value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		index := uint16(i) % a.Size()
		dot := Dot{
			X: uint8(index % uint16(a.Width())),
			Y: uint8(index / uint16(a.Width())),
		}
		m.Set(dot, object)
	}
}

func Benchmark_Map_MSet_MRemove(b *testing.B) {
	const dotsCount = 10
	const dotsPadding = dotsCount

	a := MustArea(255, 255)
	m := NewMap(a)
	object := NewObject("value")

	b.ResetTimer()

	b.StopTimer()

	for i := 0; i < b.N; i += dotsPadding {
		dotsToBeRemoved := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(i+j) % a.Size()
			dot := Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			}
			dotsToBeRemoved = append(dotsToBeRemoved, dot)
		}

		dotsToBeSet := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(i+j+dotsPadding) % a.Size()
			dot := Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			}
			dotsToBeSet = append(dotsToBeSet, dot)
		}

		b.StartTimer()

		m.MRemove(dotsToBeRemoved)
		m.MSet(dotsToBeSet, object)

		b.StopTimer()
	}
}

func Benchmark_Map_MSetIfAbsent_MRemoveObject(b *testing.B) {
	const dotsCount = 10
	const dotsPadding = dotsCount

	a := MustArea(255, 255)
	m := NewMap(a)
	object := NewObject("value")

	b.ResetTimer()

	b.StopTimer()

	for i := 0; i < b.N; i += dotsPadding {
		dotsToBeRemoved := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(i+j) % a.Size()
			dot := Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			}
			dotsToBeRemoved = append(dotsToBeRemoved, dot)
		}

		dotsToBeSet := make([]Dot, 0, dotsCount)
		for j := 0; j < dotsCount; j++ {
			index := uint16(i+j+dotsPadding) % a.Size()
			dot := Dot{
				X: uint8(index % uint16(a.Width())),
				Y: uint8(index / uint16(a.Width())),
			}
			dotsToBeSet = append(dotsToBeSet, dot)
		}

		b.StartTimer()

		m.MRemoveObject(dotsToBeRemoved, object)
		m.MSetIfAbsent(dotsToBeSet, object)

		b.StopTimer()
	}
}
