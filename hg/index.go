package hg

import (
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
	bin, err := exec.LookPath("hg")
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
	ret := make([]string, 0)
	for _, v := range strings.Split(string(out), "\n") {
		k := strings.Split(v, " ")
		if len(k) > 0 && k[0] != "tip" {
			ret = append(ret, k[0])
		}
	}
	return ret, nil
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
	return len(string(out))==0, nil
}

func CreateTag(path string, tag string) (bool, error) {
	bin, err := exec.LookPath("hg")
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
