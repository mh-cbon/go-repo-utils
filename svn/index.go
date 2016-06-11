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

func CreateTag(path string, tag string) (bool, string, error) {

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=", err)
		return false, "", err
	}

	if contains(tags, tag) {
		return false, "", errors.New("Tag '" + tag + "' already exists")
	}

	root, err := GetRepositoryRoot(path)
	if err != nil {
		logger.Printf("err=", err)
		return false, "", err
	}

	CreateTagDir(path)

	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=", err)
		return false, "", nil
	}

	args := []string{"copy", root + "/trunk", root + "/tags/" + tag, "-m", "tag: " + tag}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=", err)
	logger.Printf("out=", string(out))
	return err == nil, string(out), err
}

func CreateTagDir(path string) (string, error) {
	root, err := GetRepositoryRoot(path)
	if err != nil {
		logger.Printf("err=", err)
		return "", err
	}

	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=", err)
		return "", nil
	}

	args := []string{"mkdir", root + "/tags/", "-m", "Create tag folder"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=", err)
	logger.Printf("out=", string(out))
	return string(out), err
}

func GetRepositoryRoot(path string) (string, error) {
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=", err)
		return "", err
	}

	args := []string{"info", "."}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return "", err
	}

	logger.Printf("out=", string(out))
	p := "Repository Root:"
	root := ""
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Index(line, p) == 0 {
			root = line[len(p)+1:]
		}
	}
	return root, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
