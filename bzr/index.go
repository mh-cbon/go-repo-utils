// Bazaar implementation of go-reop-utils
package bzr

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

// Test if given path is managed by bzr with bzr status
func IsIt(path string) bool {
	bin, err := exec.LookPath("bzr")
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

// List tags on given path
func List(path string) ([]string, error) {
	tags := make([]string, 0)
	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=%s", err)
		return tags, err
	}

	args := []string{"tags"}
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
		k := strings.Split(v, " ")
		if len(k) > 0 && len(k[0]) > 0 {
			tags = append(tags, k[0])
		}
	}
	return tags, nil
}

// Check uncommited files with bzr status
func IsClean(path string) (bool, error) {
	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=%s", err)
		return false, nil
	}

	args := []string{"status"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return false, err
	}

	logger.Printf("out=%s", string(out))
	verb, _ := regexp.Compile("^(added|unknown|removed|modified):$")
	changes := make([]string, 0)
	catch := false
	for _, v := range strings.Split(string(out), "\n") {
		if verb.MatchString(v) {
			catch = strings.Index(v, "unknown") == -1
		} else if catch {
			k := strings.Split(v, " ")
			changes = append(changes, k[len(k)-1])
		}
	}

	return len(changes) == 0, nil
}

// Create given tag on path with the provided message
func CreateTag(path string, tag string, message string) (bool, string, error) {

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", err
	}

	if len(message) > 0 {
		logger.Println("Unused message: " + message)
	}

	if contains(tags, tag) {
		return false, "", errors.New("Tag '" + tag + "' already exists")
	}

	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", nil
	}

	args := []string{"tag", tag}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return err == nil, string(out), err
}

// Add given file to bzr on path
func Add(path string, file string) error {

	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=%s", err)
		return err
	}

	args := []string{"add"}
	if len(file) > 0 {
		args = append(args, []string{file}...)
	}
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

	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=%s", err)
		return err
	}

	args := []string{"commit", "-q", "--strict", "--local", "-m", message}
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
