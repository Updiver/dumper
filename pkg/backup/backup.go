package backup

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	plumbing "gopkg.in/src-d/go-git.v4/plumbing"
)

// Another does something
func Another(url, directory string, branches []string) {
	fmt.Println("Cloning repo: %s to folder: %s", url, directory)

	for _, branchName := range branches {
		git.PlainClone(directory, false, &git.CloneOptions{
			URL:               url,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			SingleBranch:      false,
			ReferenceName:     plumbing.NewBranchReferenceName(branchName),
			// ReferenceName: branchName,
		})
	}

}
