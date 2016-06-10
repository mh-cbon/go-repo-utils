package main

import (
	"os/exec"
	"fmt"
	"testing"
)

func TestGit(t *testing.T) {
	DoTestFolderUnderVcs("/home/vagrant/git", t)
	DoTestFolderUnderVcsAsJson("/home/vagrant/git", t)
	DoTestFolderUnderVcsAny("/home/vagrant/git", t)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/git", t)
	DoTestFolderIsClean("/home/vagrant/git", t)
	DoTestFolderIsCleanJson("/home/vagrant/git", t)
	DoTestFolderIsDirty("/home/vagrant/git_dirty", t)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/git_untracked", t)
	DoCreateTag("/home/vagrant/git", t)
	DoFailCreateTag("/home/vagrant/git", t)
	DoFailCreateTagMissTagName("/home/vagrant/git", t)
}

func TestHg(t *testing.T) {
	DoTestFolderUnderVcs("/home/vagrant/hg", t)
	DoTestFolderUnderVcsAsJson("/home/vagrant/hg", t)
	DoTestFolderUnderVcsAny("/home/vagrant/hg", t)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/hg", t)
	DoTestFolderIsClean("/home/vagrant/hg", t)
	DoTestFolderIsCleanJson("/home/vagrant/hg", t)
	DoTestFolderIsDirty("/home/vagrant/hg_dirty", t)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/hg_untracked", t)
	DoCreateTag("/home/vagrant/hg", t)
	DoFailCreateTag("/home/vagrant/hg", t)
	DoFailCreateTagMissTagName("/home/vagrant/hg", t)
}

func TestSvn(t *testing.T) {
	DoTestFolderUnderVcs("/home/vagrant/svn_work", t)
	DoTestFolderUnderVcsAsJson("/home/vagrant/svn_work", t)
	DoTestFolderUnderVcsAny("/home/vagrant/svn_work", t)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/svn_work", t)
	DoTestFolderIsClean("/home/vagrant/svn_work", t)
	DoTestFolderIsCleanJson("/home/vagrant/svn_work", t)
	DoTestFolderIsDirty("/home/vagrant/svn_dirty_work", t)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/svn_untracked_work", t)
	DoFailCreateTag("/home/vagrant/svn_work", t)
}

func TestBzr(t *testing.T) {
	DoTestFolderUnderVcs("/home/vagrant/bzr", t)
	DoTestFolderUnderVcsAsJson("/home/vagrant/bzr", t)
	DoTestFolderUnderVcsAny("/home/vagrant/bzr", t)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/bzr", t)
	DoTestFolderIsClean("/home/vagrant/bzr", t)
	DoTestFolderIsCleanJson("/home/vagrant/bzr", t)
	DoTestFolderIsDirty("/home/vagrant/bzr_dirty", t)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/bzr_untracked", t)
	DoCreateTag("/home/vagrant/bzr", t)
	DoFailCreateTag("/home/vagrant/bzr", t)
	DoFailCreateTagMissTagName("/home/vagrant/bzr", t)
}

func TestPathArgs(t *testing.T) {
  DoTestFolderUnderVcsWithPath("/home/vagrant/git", t)
	DoTestFolderIsCleanWithPath("/home/vagrant/bzr", t)
}

func TestFolderNotUnderVcs(t *testing.T) {
	args := []string{"list-tags"}
	cmd := exec.Command("/vagrant/build/go-repo-utils", args...)
	cmd.Dir = "/home/vagrant"

	err := cmd.Run()
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%q\n", true)
	}
}

func ExecSuccessCommand(t *testing.T, cmd string, cwd string, args []string) string {
  fmt.Printf("%s: %s %s\n", cwd, cmd, args)
	execCmd := exec.Command(cmd, args...)
	execCmd.Dir = cwd

	out, err := execCmd.CombinedOutput()
	if err != nil {
    fmt.Println(string(out))
		t.Errorf("Expected err=nil, got err=%s\n", err)
    return ""
	}
	if execCmd.ProcessState != nil && execCmd.ProcessState.Success() == false {
		t.Errorf("Expected success=true, got success=%q\n", false)
    return ""
	}

  return string(out)
}

func DoTestFolderUnderVcsWithPath(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-p", path}
  out := ExecSuccessCommand(t, cmd, "/home", args)

	expectedOut := "1.0.0\n1.0.2\n1.0.3\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsCleanWithPath(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
  args := []string{"is-clean", "--path="+path}
  out := ExecSuccessCommand(t, cmd, "/home", args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoCreateTag(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
	args := []string{"create-tag", "1.0.3"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "done\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoFailCreateTag(path string, t *testing.T) {
  args := []string{"create-tag", "1.0.3"}
	cmd := exec.Command("/vagrant/build/go-repo-utils", args...)
	cmd.Dir = path

	err := cmd.Run()
	if err==nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%q\n", true)
	}
}

func DoFailCreateTagMissTagName(path string, t *testing.T) {
  args := []string{"create-tag"}
	cmd := exec.Command("/vagrant/build/go-repo-utils", args...)
	cmd.Dir = path

	err := cmd.Run()
	if err==nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%q\n", true)
	}
}

func DoTestFolderUnderVcs(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "1.0.0\n1.0.2\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderUnderVcsAsJson(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-j"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "[\"1.0.0\",\"1.0.2\"]"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderUnderVcsAny(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-a"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "1.0.0\n1.0.2\nnotsemvertag\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderUnderVcsAnyReversed(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-a", "-r"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "notsemvertag\n1.0.2\n1.0.0\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsClean(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
  args := []string{"is-clean"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsCleanJson(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
  args := []string{"is-clean", "-j"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "true"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsDirty(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
  args := []string{"is-clean"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "no\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}
func DoTestFolderIsCleanEvenWithUntrackedFiles(path string, t *testing.T) {
  cmd := "/vagrant/build/go-repo-utils"
  args := []string{"is-clean"}
  out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}
