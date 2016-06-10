package svn

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func IsIt(path string) bool {
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=", err)
		return false
	}

	args := []string{"list"}
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
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=", err)
		return make([]string, 0), err
	}

	args := []string{"ls", "^/tags"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return make([]string, 0), err
	}

	logger.Printf("out=", string(out))
	ret := strings.Split(string(out), "\n")
	for i, v := range ret {
		if len(v) > 0 {
			ret[i] = v[0 : len(v)-1]
		}
	}
	return ret, nil
}

func IsClean(path string) (bool, error) {
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=", err)
		return false, err
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

func CreateTag(path string, tag string) (bool, error) {
	return true, errors.New("Cannot create tag with svn")
}
