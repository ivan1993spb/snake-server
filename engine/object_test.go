package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewObject(t *testing.T) {
	var value = "test123"

	o := NewObject(value)

	require.Equal(t, value, o.value)
}

func Test_Object_Value_ReturnsValue(t *testing.T) {
	var value = "test"

	o := Object{
		value: value,
	}

	require.Equal(t, value, o.Value())
}
