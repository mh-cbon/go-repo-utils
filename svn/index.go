// Package svn implements go-repo-utils interfaces.
package svn

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
	bin, err := exec.LookPath("svn")
	if err != nil {
		logger.Printf("err=%s", err)
		return nil, err
	}
	logger.Printf("%s %s", bin, args)
	cmd := exec.Command(bin, args...)
	cmd.Dir = path
	return cmd, nil
}

// IsIt Tests if path is managed by SVN using svn list
func IsIt(path string) bool {

	args := []string{"list"}
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

// List svn tags with svn ls ^/tags of given path
func List(path string) ([]string, error) {
	tags := make([]string, 0)

	args := []string{"ls", "^/tags"}
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
		if len(v) > 0 {
			tags = append(tags, v[0:len(v)-1])
		}
	}
	return tags, nil
}

// IsClean Checks uncommited files with svn -q of given path
func IsClean(path string) (bool, error) {

	args := []string{"status", "-q"}
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

// CreateTag Creates given tag at root/tags/[tag] on path with the provided message
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

	args := []string{"copy", root + "/trunk", root + "/tags/" + tag}
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

// CreateTagDir Create an svn tag directory at root/tags/
func CreateTagDir(path string) (string, error) {
	root, err := GetRepositoryRoot(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return "", err
	}

	args := []string{"mkdir", root + "/tags/", "-m", "Create tag folder"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return "", err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return string(out), err
}

// GetRepositoryRoot returns svn root path according to svn info .
func GetRepositoryRoot(path string) (string, error) {

	d, err := GetRepositoryInfo(path)
	if err != nil {
		return "", err
	}

	if x, ok := d["Repository Root"]; ok {
		return x, nil
	}
	return "", nil
}

// GetRepositoryInfo parses svn info to a map.
func GetRepositoryInfo(path string) (map[string]string, error) {

	ret := map[string]string{}

	args := []string{"info", "."}
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("err=%s", err)
		return ret, err
	}

	logger.Printf("out=%s", string(out))
	ID := regexp.MustCompile(`^([^:]+):(.+)`)
	for _, line := range strings.Split(string(out), "\n") {
		if ID.MatchString(line) {
			res := ID.FindAllStringSubmatch(line, -1)
			k, v := res[0][1], res[0][2]
			ret[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return ret, nil
}

// Add given file to svn on path
func Add(path string, file string) error {

	args := []string{"add", file}
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

// ListCommitsBetween List commits between two points
func ListCommitsBetween(path string, since string, to string) ([]commit.Commit, error) {
	ret := make([]commit.Commit, 0)

	tags, err := List(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return ret, err
	}

	if p := pos(tags, since); p > -1 {
		s, err2 := GetRevisionTag(path, since)
		if err2 != nil {
			logger.Printf("err=%s", err2)
			return ret, err2
		}
		if s != "" {
			since = s
		}
	}
	if p := pos(tags, to); p > -1 {
		t, err2 := GetRevisionTag(path, to)
		if err2 != nil {
			logger.Printf("err=%s", err2)
			return ret, err2
		}
		if t != "" {
			to = t
		}
	}
	if since == "" {
		since = "0"
	}

	args := []string{"log"}
	if len(since)+len(to) > 0 {
		args = append(args, "-r", since+":"+to)
	}
	args = append(args, "^/.")
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))

	ret = ParseSvnLog(string(out))

	return ret, err
}

// ParseSvnLog parses an svn log string to a list of commits.
func ParseSvnLog(log string) []commit.Commit {
	ret := make([]commit.Commit, 0)

	splitRe := regexp.MustCompile(`^[-]+$`)
	infoRe := regexp.MustCompile(`^r([0-9]+)\s+\|\s+([^|]+)\|\s+([^\(]+)`)
	var c *commit.Commit
	for _, line := range strings.Split(log, "\n") {
		line = strings.TrimSpace(line)
		if splitRe.MatchString(line) {
			if c != nil {
				ret = append(ret, *c)
			}
			c = &commit.Commit{}
		} else if c != nil && c.Revision == "" && infoRe.MatchString(line) {
			res := infoRe.FindStringSubmatch(line)
			c.Revision = res[1]
			c.Author = strings.TrimSpace(res[2])
			c.Date = strings.TrimSpace(res[3])
		} else if c != nil && line != "" {
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

// GetRevisionTag Get the revision of a tag
func GetRevisionTag(path string, tag string) (string, error) {
	ret := ""

	root, err := GetRepositoryRoot(path)
	if err != nil {
		logger.Printf("err=%s", err)
		return ret, err
	}

	args := []string{"log", root + "/tags/" + tag, "-v", "--stop-on-copy"}
	cmd, err := getCmd(path, args)
	if err != nil {
		return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	re := regexp.MustCompile(`\s+A\s+\/tags\/[^\s]+\s+\(from \/[^:]+:([0-9]+)\)`)
	res := re.FindStringSubmatch(string(out))
	if len(res) > 0 {
		ret = string(res[1])
	}
	return ret, err
}

// GetFirstRevision returns the first revision of the repository.
func GetFirstRevision(path string) (string, error) {
	return "1", nil
}

func pos(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
