package dumper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	ErrRepositoryIsNonBare = errors.New("repository is already non-bare")
	ErrRepositoryIsBare    = errors.New("repository is already bare")
)

type RepositoryType int

const (
	RepositoryTypeBare RepositoryType = iota
	RepositoryTypeNonBare
)

type Converter interface {
	Convert(destination string, repoType RepositoryType) error
}

// Convert converts git repository from bare to non-bare and vice versa
func Convert(destination string, convertToRepoType RepositoryType) error {
	repository, err := Repository(destination)
	if err != nil {
		return fmt.Errorf("convert repository: %w", err)
	}
	config, err := repository.Config()
	if err != nil {
		return fmt.Errorf("convert repository: %w", err)
	}

	switch convertToRepoType {
	case RepositoryTypeBare:
		if ok := config.Core.IsBare; ok {
			return ErrRepositoryIsBare
		}

		err = convertToBare(destination)
		if err != nil {
			return fmt.Errorf("convert repository: %w", err)
		}

		return nil
	case RepositoryTypeNonBare:
		if ok := config.Core.IsBare; !ok {
			return ErrRepositoryIsNonBare
		}
		return convertToNonBare(repository, destination)
	default:
		return errors.New("unknown repository type")
	}
}

func convertToNonBare(repository *git.Repository, destination string) error {
	rConfig, err := repository.Config()
	if err != nil {
		return err
	}

	rConfig.Core.IsBare = false
	repository.SetConfig(rConfig)

	err = os.Mkdir(filepath.Join(destination, ".git"), 0755)
	if err != nil {
		return fmt.Errorf("init repo: %w", err)
	}

	// move all content of repository into .git
	// we need this because after checkout files from
	// branch will be placed into root directory
	err = moveFolderContent(destination, filepath.Join(destination, ".git"))
	if err != nil {
		return fmt.Errorf("move bare repo content into .git: %w", err)
	}

	// once git-related files are moved into .git
	// it won't have any knowledge about branches
	// so we need to read .git folder again as current
	// repository instance does nothing about new fs structure
	repository, err = git.PlainOpen(destination)
	if err != nil {
		return fmt.Errorf("open repository: %w", err)
	}

	refs, err := repository.References()
	if err != nil {
		return err
	}
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			branchName := ref.Name().Short()
			branchRef := plumbing.NewBranchReferenceName(branchName)
			worktree, err := repository.Worktree()
			if err != nil {
				return fmt.Errorf("create worktree %s: %s", branchName, err)
			}

			err = worktree.Checkout(&git.CheckoutOptions{
				Branch: branchRef,
				Create: false,
				Force:  true,
			})
			if err != nil {
				return fmt.Errorf("checkout branch %s: %s", branchName, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("convert repository to non bare: %w", err)
	}

	return nil
}

// INFO: TO BE DONE
func convertToBare(destination string) error {
	return nil
}

func moveFolderContent(source, destination string) error {
	files, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}

		err := os.Rename(filepath.Join(source, file.Name()), filepath.Join(destination, file.Name()))
		if err != nil {
			return fmt.Errorf(
				"move file/dir[%s] -> to [%s]: %w",
				file.Name(),
				filepath.Join(destination, file.Name()),
				err,
			)
		}
	}

	return nil
}
