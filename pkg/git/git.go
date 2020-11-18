package git

import (
	"os/exec"
	"path"
	"strings"
)

func GetDiffFiles(repoPath string) ([]string, error) {
	diff, err := getDiff(repoPath, "HEAD^", "HEAD")
	if err != nil {
		return nil, err
	}

	var diffWithFullPath []string

	for _, diffItem := range diff {
		diffWithFullPath = append(diffWithFullPath, path.Join(repoPath, diffItem))
	}

	return diffWithFullPath, nil
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

func getDiff(path string, first, second string) ([]string, error) {
	// git diff HEAD^ HEAD
	command := exec.Command("git", "diff", "--name-only", first, second)
	command.Dir = path
	result, err := command.Output()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(string(result)), "\n")

	return parts, nil
}
