// Git implementation of go-reop-utils
package git

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
	bin, err := exec.LookPath("git")
	if err != nil {
		logger.Printf("err=%s", err)
		return nil, err
	}
	logger.Printf("%s %s (cwd=%s)", bin, args, path)
	cmd := exec.Command(bin, args...)
	cmd.Dir = path
	return cmd, nil
}

// Test if given path is managed by git with git info
func IsIt(path string) bool {
	args := []string{"rev-parse"}
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

	args := []string{"tag"}
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
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) > 0 {
			tags = append(tags, line)
		}
	}
	return tags, nil
}

// Check uncommited files with git status --porcelain --untracked-files=no
func IsClean(path string) (bool, error) {

	args := []string{"status", "--porcelain", "--untracked-files=no"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return false, err
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

	args := []string{"tag", "-a", tag}
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
	return err == nil, string(out), err
}

// Add given file to git on path
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

// List commits between two points
func ListCommitsBetween(path string, since string, to string) ([]commit.Commit, error) {
	ret := make([]commit.Commit, 0)

	args := []string{"log"}
	if len(since)+len(to) > 0 {
		revset := ""
		if since != "" {
			revset += since + ".."
		}
		revset += to
		args = append(args, revset)
	}
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))

	ret = ParseGitLog(string(out))

	return ret, err
}

func ParseGitLog(logs string) []commit.Commit {
	ret := make([]commit.Commit, 0)

	commitRe := regexp.MustCompile(`^commit\s+(.+)$`)
	authorRe := regexp.MustCompile(`^Author:\s+([^<]+)\s+<([^>]+)>$`)
	dateRe := regexp.MustCompile(`^Date:\s*(.+)$`)
	messageRe := regexp.MustCompile(`^\s+(.+)$`)
	var c *commit.Commit
	for _, line := range strings.Split(logs, "\n") {
		if commitRe.MatchString(line) {
			if c != nil {
				ret = append(ret, *c)
			}
			c = &commit.Commit{}
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
			res := messageRe.FindStringSubmatch(line)
			if c.Message == "" {
				c.Message = strings.TrimSpace(res[1])
			} else {
				c.Message = c.Message + "\n" + strings.TrimSpace(res[1])
			}
		}
	}
	if c != nil && c.Revision != "" {
		ret = append(ret, *c)
	}
	return ret
}

func GetRevisionTag(path string, tag string) (string, error) {
	ret := ""

	args := []string{"log", "-n", "1", tag}
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))

	return strings.TrimSpace(string(out)), err
}

func GetFirstRevision(path string) (string, error) {
	ret := ""

	args := []string{"rev-list", "--max-parents=0", "HEAD"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))

  // when a merge has occured, it will return multiple hash,
  // take the last one only 
  sout := strings.Split(strings.TrimSpace(string(out)), "\n")
	return sout[len(sout)-1], err
}
