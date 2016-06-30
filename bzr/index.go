// Bazaar implementation of go-reop-utils
package bzr

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mh-cbon/verbose"
	"github.com/mh-cbon/go-repo-utils/commit"
)

var logger = verbose.Auto()

func getCmd(path string, args []string) (*exec.Cmd, error) {
	bin, err := exec.LookPath("bzr")
	if err != nil {
		logger.Printf("err=%s", err)
		return nil, err
	}
	logger.Printf("%s %s (cwd=%s)", bin, args, path)
  cmd := exec.Command(bin, args...)
  cmd.Dir = path
  return cmd, nil
}

// Test if given path is managed by bzr with bzr info
func IsIt(path string) bool {

	args := []string{"info"}
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
		if len(k) > 0 && len(k[0]) > 0 {
			tags = append(tags, k[0])
		}
	}
	return tags, nil
}

// Check uncommited files with bzr status
func IsClean(path string) (bool, error) {

	args := []string{"status"}
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

	args := []string{"tag", tag}
	cmd, err := getCmd(path, args)
	if err != nil {
    return false, "", err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))
	return err == nil, string(out), err
}

// Add given file to bzr on path
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

	args := []string{"commit", "-q", "--strict", "--local", "-m", message}
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
func ListCommitsBetween (path string, since string, to string) ([]commit.Commit, error) {
  ret := make([]commit.Commit, 0)

  if to=="HEAD" {
    to = ""
  }

	args := []string{"log"}
  if len(since)+len(to)>0 {
    if since=="" {
      since = "revno:1"
    } else if (IsTag(path, since)) {
      since = "tag:"+since
    }
    if (to!="" && IsTag(path, to)) {
      to = "tag:"+to
    }
    args = append(args, "-r", since+".."+to)
  }
	cmd, err := getCmd(path, args)
	if err != nil {
    return ret, err
	}

	out, err := cmd.CombinedOutput()
	logger.Printf("err=%s", err)
	logger.Printf("out=%s", string(out))

  ret = ParseBzrLogs(string(out))

	return ret, err
}

func ParseBzrLogs (log string) []commit.Commit {
  ret := make([]commit.Commit, 0)

  splitRe := regexp.MustCompile(`^[-]+$`)
  commitRe := regexp.MustCompile(`^revno:\s+([0-9]+)$`)
  authorRe := regexp.MustCompile(`^committer:\s+([^<]+)\s+<([^>]+)>$`)
  dateRe := regexp.MustCompile(`^timestamp:\s*(.+)$`)
  messageRe := regexp.MustCompile(`message:$`)
  isInMessage := false
  var c *commit.Commit
  for _, line := range strings.Split(log, "\n") {
    line = strings.TrimSpace(line)
    if splitRe.MatchString(line) {
      if c!=nil {
        ret = append(ret, *c)
      }
      c = &commit.Commit{}
      isInMessage = false
    } else if commitRe.MatchString(line) {
      res := commitRe.FindStringSubmatch(line)
      c.Revision = res[1]
    } else if c!=nil && authorRe.MatchString(line) {
      res := authorRe.FindStringSubmatch(line)
      c.Author = strings.TrimSpace(res[1])
      c.Email = res[2]
    } else if c!=nil && dateRe.MatchString(line) {
      res := dateRe.FindStringSubmatch(line)
      c.Date = res[1]
    } else if c!=nil && messageRe.MatchString(line) {
      isInMessage = true
    } else if c!=nil && isInMessage && line != "" {
      if c.Message=="" {
        c.Message = line
      } else {
        c.Message = c.Message+"\n"+line
      }
    }
  }
  if c!=nil && c.Revision!="" {
    ret = append(ret, *c)
  }
  return ret
}

// Get revision of a tag
func GetRevisionTag(path string, tag string) (string, error) {
	ret := ""

	args := []string{"tags"}
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
	for _, v := range strings.Split(string(out), "\n") {
		k := strings.Split(v, " ")
		if len(k) > 0 && len(k[0]) > 0 {
      if k[0]==tag {
        ret = strings.TrimSpace(k[1])
        break
      }
		}
	}
	return ret, nil
}

func IsTag(path string, tag string) bool {
  tags, err := List(path)
  if err!=nil {
    return false
  }

	return contains(tags, tag)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
