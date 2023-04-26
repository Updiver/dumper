package dumper

import "testing"

func TestNew(t *testing.T) {
	dumper := New(
		SetVCSHoster(VCSGithub),
	)

	vcsHoster := dumper.VCSHoster()
	if vcsHoster != VCSGithub {
		t.Errorf("vcsHoster = %v, want %v", vcsHoster, VCSGithub)
	}
}
