package backup

import (
	"fmt"

	"log"

	"os"
	"os/user"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Clone does something
func Clone(url, directory string, creds struct {
	Username string
	Password string
} /* branches []string */) {
	fmt.Printf("Cloning repo: %s to folder: %s\n", url, directory)
	systemUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed getting current system user: %s", err)
		return
	}

	fullDirectoryPath := systemUser.HomeDir + "/cloned_repositories/" + directory

	fmt.Printf("Repo will be cloned here: %s\n", fullDirectoryPath)

	// TODO: implement fetching branches
	// TODO: implement letting user know which repos are > 1Gb
	// and leave links for them to do that manually
	_, err = git.PlainClone(fullDirectoryPath, false, &git.CloneOptions{
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
