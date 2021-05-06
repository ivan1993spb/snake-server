package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewContainer(t *testing.T) {
	var object = "test123"

	o := NewContainer(object)

	require.Equal(t, object, o.object)
}

func Test_Container_GetObject_ReturnsTheRightObject(t *testing.T) {
	var object = "test"

	o := Container{
		object: object,
	}

	require.Equal(t, object, o.GetObject())
}
