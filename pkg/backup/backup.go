package backup

import (
	"fmt"

	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Clone does something
func Clone(url, directory string, creds struct {
	Username string
	Password string
} /* branches []string */) {
	fmt.Printf("Cloning repo: %s to folder: %s\n", url, directory)
	// systemUser, err := user.Current()
	// if err != nil {
	// 	log.Fatalf("Failed getting current system user: %s", err)
	// 	return
	// }

	fullDirectoryPath := directory

	fmt.Printf("Repo will be cloned here: %s\n", fullDirectoryPath)

	// TODO: implement fetching branches
	// TODO: implement letting user know which repos are > 1Gb
	// and leave links for them to do that manually
	_, err := git.PlainClone(fullDirectoryPath, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		SingleBranch:      false,
		Progress:          os.Stderr,
		Auth: &http.BasicAuth{
			Username: creds.Username,
			Password: creds.Password,
		},
		// ReferenceName:     plumbing.NewBranchReferenceName(branchName),
	})

	if err != nil {
		fmt.Printf("Failed pulling repository: %s\n", err)
	} else {
		fmt.Println("Repo cloned OK")
	}

}
