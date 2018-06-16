package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Dot_Hash(t *testing.T) {
	require.Equal(t, uint16(0xffaa), Dot{X: 0xff, Y: 0xaa}.Hash())
	require.Equal(t, uint16(0xabcd), Dot{X: 0xab, Y: 0xcd}.Hash())
}
