package dumper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPositiveBoolRef(t *testing.T) {
	expected := true
	actual := PositiveBoolRef()
	require.Equal(t, expected, *actual, "expect to have positive bool ref")
}

func TestNegaviteBoolRef(t *testing.T) {
	expected := false
	actual := NegativeBoolRef()
	require.Equal(t, expected, *actual, "expect to have negative bool ref")
}
