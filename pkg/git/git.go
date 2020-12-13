package git

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func GetDiffFiles(repoPath string, revision string) ([]string, error) {
	diff, err := getDiff(repoPath, revision)
	if err != nil {
		return nil, err
	}

	var diffWithFullPath []string

	for _, diffItem := range diff {
		diffWithFullPath = append(diffWithFullPath, path.Join(repoPath, diffItem))
	}

	return diffWithFullPath, nil
}

func GetOldFile(repoPath string, file string, revision string) ([]byte, error) {
	command := exec.Command("git", "show", fmt.Sprintf("%s:%s", revision, file))
	command.Dir = repoPath

	result, err := command.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if ok {
			return nil, fmt.Errorf("%s", string(ee.Stderr))
		}

		return nil, err
	}

	return result, nil
}

func ResolveGitPath(dir string) (string, error) {
	return rootPath(dir)
}

func rootPath(applicationPath string) (string, error) {
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	command.Dir = applicationPath
	resultPath, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(resultPath)), nil
}

func getDiff(path string, first string) ([]string, error) {
	// git diff HEAD^ HEAD
	command := exec.Command("git", "diff", "--name-only", first)
	command.Dir = path
	result, err := command.Output()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(string(result)), "\n")

	return parts, nil
}
