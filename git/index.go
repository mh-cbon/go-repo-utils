package git

import (
	"os/exec"
	"strings"

	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func IsIt(path string) bool {
	bin, err := exec.LookPath("git")
	if err != nil {
		logger.Printf("err=%s", err)
		return false
	}

	args := []string{"status"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return false
	}

	logger.Printf("out=%s", string(out))
	return cmd.ProcessState != nil && cmd.ProcessState.Success()
}

func List(path string) ([]string, error) {
	tags := make([]string, 0)
	bin, err := exec.LookPath("git")
	if err != nil {
		logger.Printf("err=%s", err)
		return tags, err
	}

	args := []string{"tag"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return tags, err
	}

	logger.Printf("out=%s", string(out))
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) > 0 {
			tags = append(tags, line)
		}
	}
	return tags, nil
}

func IsClean(path string) (bool, error) {
	bin, err := exec.LookPath("git")
	if err != nil {
		logger.Printf("err=%s", err)
		return false, nil
	}

	args := []string{"status", "--porcelain", "--untracked-files=no"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return false, err
	}

	logger.Printf("out=%s", string(out))
	return len(string(out)) == 0, nil
}

func CreateTag(path string, tag string, message string) (bool, string, error) {

	bin, err := exec.LookPath("git")
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", nil
	}

	args := []string{"tag", tag}
  if len(message)>0 {
    args = append(args, []string{"-m", message}...)
  }
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return err == nil, string(out), err
}
