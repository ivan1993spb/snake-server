package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Dot_Hash(t *testing.T) {
	require.Equal(t, uint16(0xffaa), Dot{X: 0xff, Y: 0xaa}.Hash())
	require.Equal(t, uint16(0xabcd), Dot{X: 0xab, Y: 0xcd}.Hash())
	require.Equal(t, uint16(0x0), Dot{X: 0x0, Y: 0x0}.Hash())
}

func Test_HashToDot(t *testing.T) {
	require.Equal(t, Dot{X: 0xff, Y: 0xaa}, HashToDot(0xffaa))
	require.Equal(t, Dot{X: 0xab, Y: 0xcd}, HashToDot(0xabcd))
	require.Equal(t, Dot{X: 0x0, Y: 0x0}, HashToDot(0x0))
}

func Test_Dot_Equals(t *testing.T) {
	require.True(t, Dot{0, 0}.Equals(Dot{0, 0}))
	require.True(t, Dot{1, 1}.Equals(Dot{1, 1}))
	require.True(t, Dot{5, 5}.Equals(Dot{5, 5}))
	require.True(t, Dot{0xff, 0xff}.Equals(Dot{0xff, 0xff}))

	require.False(t, Dot{0xff, 0xff}.Equals(Dot{0xff, 0x0}))
	require.False(t, Dot{0, 0}.Equals(Dot{0, 1}))
	require.False(t, Dot{255, 0}.Equals(Dot{0, 255}))
	require.False(t, Dot{0, 0xff}.Equals(Dot{0xff, 0}))
}
