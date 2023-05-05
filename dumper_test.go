package dumper

import (
	"testing"
)

// dummy test
func TestNew(t *testing.T) {
	dumper := New()

	if dumper == nil {
		t.Errorf("expect to have proper dump instance")
	}
}
