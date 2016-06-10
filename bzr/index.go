package bzr

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func IsIt(path string) bool {
	bin, err := exec.LookPath("bzr")
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
	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=", err)
		return make([]string, 0), err
	}

	args := []string{"tags"}
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
		k := strings.Split(v, " ")
		if len(k) > 0 {
			ret[i] = k[0]
		}
	}
	return ret, nil
}

func IsClean(path string) (bool, error) {
	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=", err)
		return false, nil
	}

	args := []string{"status"}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return false, err
	}

	logger.Printf("out=", string(out))
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

func CreateTag(path string, tag string) (bool, error) {

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=", err)
		return false, err
	}

	if contains(tags, tag) {
		return false, errors.New("Tag '" + tag + "' already exists")
	}

	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=", err)
		return false, nil
	}

	args := []string{"tag", tag}
	cmd := exec.Command(bin, args...)
	cmd.Dir = path

	logger.Printf("%s %s (cwd=%s)", bin, args, path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=", err)
		return false, err
	}

	logger.Printf("out=", string(out))
	return true, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
