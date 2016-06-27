// Svn implementation of go-repo-utils
package svn

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

// Test if path is managed by SVN using svn list
func IsIt(path string) bool {
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return false
	}

	args := []string{"list"}
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

// List svn tags with svn ls ^/tags of given path
func List(path string) ([]string, error) {
	tags := make([]string, 0)
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return tags, err
	}

	args := []string{"ls", "^/tags"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return tags, err
	}

	logger.Printf("out=%s", string(out))
	for _, v := range strings.Split(string(out), "\n") {
		if len(v) > 0 {
			tags = append(tags, v[0:len(v)-1])
		}
	}
	return tags, nil
}

// Check uncommited files with svn -q of given path
func IsClean(path string) (bool, error) {
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return false, err
	}

	args := []string{"status", "-q"}
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

// Create given tag at root/tags/[tag] on path with the provided message
func CreateTag(path string, tag string, message string) (bool, string, error) {

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", err
	}

	if contains(tags, tag) {
		return false, "", errors.New("Tag '" + tag + "' already exists")
	}

	root, err := GetRepositoryRoot(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", err
	}

	CreateTagDir(path)

	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", nil
	}

	args := []string{"copy", root + "/trunk", root + "/tags/" + tag}
	if len(message) > 0 {
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

// Create root/tags directory
func CreateTagDir(path string) (string, error) {
	root, err := GetRepositoryRoot(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return "", err
	}

	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return "", nil
	}

	args := []string{"mkdir", root + "/tags/", "-m", "Create tag folder"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return string(out), err
}

// Get svn root path using svn info .
func GetRepositoryRoot(path string) (string, error) {
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return "", err
	}

	args := []string{"info", "."}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return "", err
	}

	logger.Printf("out=%s", string(out))
	p := "Repository Root:"
	root := ""
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Index(line, p) == 0 {
			root = line[len(p)+1:]
		}
	}
	return root, nil
}

// Add given file to svn on path
func Add(path string, file string) error {

	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return err
	}

	args := []string{"add", file}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return err
}

// Commit given files with message on path
func Commit(path string, message string, files []string) error {

	if len(message) == 0 {
		return errors.New("Message is required")
	}

	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return err
	}

	args := []string{"commit", "-m", message}
	if len(files) > 0 {
		args = append(args, files...)
	}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
