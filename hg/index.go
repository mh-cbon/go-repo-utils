// Hg implementation of go-reop-utils
package hg

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mh-cbon/go-repo-utils/commit"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func getCmd(path string, args []string) (*exec.Cmd, error) {
	bin, err := exec.LookPath("hg")
	if err != nil {
		logger.Printf("err=%s", err)
		return nil, err
	}
	logger.Printf("%s %s (cwd=%s)", bin, args, path)
	cmd := exec.Command(bin, args...)
	cmd.Dir = path
	return cmd, nil
}

// Test if given path is managed by hg with hg status
func IsIt(path string) bool {
	args := []string{"status"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return false
	}

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

	args := []string{"tags"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return tags, err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return tags, err
	}

	logger.Printf("out=%s", string(out))
	for _, v := range strings.Split(string(out), "\n") {
		k := strings.Split(v, " ")
		if len(k) > 0 && k[0] != "tip" {
			tags = append(tags, k[0])
		}
	}
	return tags, nil
}

// Check uncommited files with hg status -q
func IsClean(path string) (bool, error) {

	args := []string{"status", "-q"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return false, nil
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return false, err
	}

	logger.Printf("out=%s", string(out))
	return len(string(out)) == 0, nil
}

// Create given tag on path with the provided message
func CreateTag(path string, tag string, message string) (bool, string, error) {

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return false, "", err
	}

	if contains(tags, tag) {
		return false, "", errors.New("Tag '" + tag + "' already exists")
	}

	args := []string{"tag", tag}
	if len(message) > 0 {
		args = append(args, []string{"-m", message}...)
	}
	cmd, err := getCmd(path, args)
	if err != nil {
		return false, "", err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return err == nil, string(out), nil
}

// Add given file to hg on path
func Add(path string, file string) error {

	args := []string{"add"}
	if len(file) > 0 {
		args = append(args, []string{file}...)
	}
	cmd, err := getCmd(path, args)
	if err != nil {
		return err
	}

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

	args := []string{"commit", "-m", message}
	if len(files) > 0 {
		args = append(args, files...)
	}
	cmd, err := getCmd(path, args)
	if err != nil {
		return err
	}

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

// List commits between two points
func ListCommitsBetween(path string, since string, to string) ([]commit.Commit, error) {
	ret := make([]commit.Commit, 0)

	if to == "HEAD" {
		to = "tip"
	}
	if since == "" {
		since = "0"
	}

	args := []string{"log", "-v"}
	if len(since)+len(to) > 0 {
		args = append(args, "-r", since+".."+to)
	}
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))

	ret = ParseHgLogs(string(out))

	return ret, err
}

func ParseHgLogs(log string) []commit.Commit {
	ret := make([]commit.Commit, 0)

	commitRe := regexp.MustCompile(`^changeset:\s+[0-9]+:([^\s]+)$`)
	authorRe := regexp.MustCompile(`^user:\s+([^<]+)\s+<([^>]+)>$`)
	dateRe := regexp.MustCompile(`^date:\s*(.+)$`)
	messageRe := regexp.MustCompile(`description:$`)
	isInMessage := false
	var c *commit.Commit
	for _, line := range strings.Split(log, "\n") {
		line = strings.TrimSpace(line)
		if commitRe.MatchString(line) {
			if c != nil {
				ret = append(ret, *c)
			}
			c = &commit.Commit{}
			isInMessage = false
			res := commitRe.FindStringSubmatch(line)
			c.Revision = res[1]
		} else if c != nil && authorRe.MatchString(line) {
			res := authorRe.FindStringSubmatch(line)
			c.Author = strings.TrimSpace(res[1])
			c.Email = res[2]
		} else if c != nil && dateRe.MatchString(line) {
			res := dateRe.FindStringSubmatch(line)
			c.Date = res[1]
		} else if c != nil && messageRe.MatchString(line) {
			isInMessage = true
		} else if c != nil && isInMessage && line != "" {
			if c.Message == "" {
				c.Message = line
			} else {
				c.Message = c.Message + "\n" + line
			}
		}
	}
	if c != nil && c.Revision != "" {
		ret = append(ret, *c)
	}
	return ret
}

// Get revision of a tag
func GetRevisionTag(path string, tag string) (string, error) {
	rev := ""

	args := []string{"tags"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return tag, err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return tag, err
	}

	revRe := regexp.MustCompile(`^[0-9]+[:;](.+)$`)
	logger.Printf("out=%s", string(out))
	for _, v := range strings.Split(string(out), "\n") {
		k := strings.Split(v, " ")
		if len(k) > 0 && k[0] == tag {
			if revRe.MatchString(k[1]) {
				res := revRe.FindStringSubmatch(k[1])
				rev = res[1]
				break
			}
		}
	}
	return rev, nil
}
