package hg

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func IsIt(path string) bool {
	bin, err := exec.LookPath("hg")
	if err != nil {
		logger.Printf("err=", err)
		return false
	}

	args := []string{"status"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return false
	}

	logger.Printf("out=", string(out))
	return cmd.ProcessState != nil && cmd.ProcessState.Success()
}

func List(path string) ([]string, error) {
  tags := make([]string, 0)
	bin, err := exec.LookPath("hg")
	if err != nil {
		logger.Printf("err=", err)
		return tags, err
	}

	args := []string{"tags"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return tags, err
	}

	logger.Printf("out=", string(out))
	for _, v := range strings.Split(string(out), "\n") {
		k := strings.Split(v, " ")
		if len(k) > 0 && k[0] != "tip" {
			tags = append(tags, k[0])
		}
	}
	return tags, nil
}

func IsClean(path string) (bool, error) {
	bin, err := exec.LookPath("hg")
	if err != nil {
		logger.Printf("err=", err)
		return false, nil
	}

	args := []string{"status", "-q"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return false, err
	}

	logger.Printf("out=", string(out))
	return len(string(out)) == 0, nil
}

func CreateTag(path string, tag string) (bool, string, error) {

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=", err)
		return false, "", err
	}

	if contains(tags, tag) {
		return false, "", errors.New("Tag '" + tag + "' already exists")
	}

	bin, err := exec.LookPath("hg")
	if err != nil {
		logger.Printf("err=", err)
		return false, "", nil
	}

	args := []string{"tag", tag}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=", err)
	logger.Printf("out=", string(out))
	return err == nil, string(out), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
