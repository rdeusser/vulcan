package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rdeusser/stacktrace"
)

// Init initializes a new git repo, adds all files, and creates the first
// commit.
func Init(baseDir, branch, remote string) error {
	if baseDir == "" {
		return fmt.Errorf("baseDir cannot be empty")
	}

	if branch == "" {
		return fmt.Errorf("branch cannot be empty")
	}

	if remote == "" {
		return fmt.Errorf("remote cannot be empty")
	}

	// Create an empty git repo.
	repo, err := git.PlainInit(baseDir, false)
	if err != nil {
		return stacktrace.Propagate(err, "creating empty git repository")
	}

	// Setup a worktree so we can add files.
	worktree, err := repo.Worktree()
	if err != nil {
		return stacktrace.Propagate(err, "creating worktree")
	}

	addOptions := &git.AddOptions{
		All: true,
	}

	// Add all files.
	if err := worktree.AddWithOptions(addOptions); err != nil {
		return stacktrace.Propagate(err, "adding objects to worktree")
	}

	commitOptions := &git.CommitOptions{
		All: true,
	}

	// Commit everything.
	if _, err := worktree.Commit("initial commit", commitOptions); err != nil {
		return stacktrace.Propagate(err, "commiting objects to worktree")
	}

	headRef, err := repo.Head()
	if err != nil {
		return stacktrace.Propagate(err, "getting HEAD reference")
	}

	branchRef := plumbing.NewBranchReferenceName(branch)

	// Create and checkout the new branch.
	checkoutOptions := &git.CheckoutOptions{
		Branch: branchRef,
		Create: true,
	}

	if err := worktree.Checkout(checkoutOptions); err != nil {
		return stacktrace.Propagate(err, "checking out %s", branchRef)
	}

	// Remove the reference to the old master branch.
	if err := repo.Storer.RemoveReference(headRef.Name()); err != nil {
		return stacktrace.Propagate(err, "removing refs")
	}

	return nil
}
