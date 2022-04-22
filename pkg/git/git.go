package git

import (
	"os"
	"os/exec"
)

const InitialCommit = "initial commit"

type Repo struct {
	dir string
}

func NewRepo(dir, branch string) (Repo, error) {
	r := Repo{dir: dir}

	return r, r.run("init", "-q", "-b", branch)
}

func (r Repo) AddAll() error {
	return r.run("add", "-A")
}

func (r Repo) Commit(msg string) error {
	return r.run("commit", "-q", "-m", msg)
}

func (r Repo) run(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
