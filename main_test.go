package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/mh-cbon/go-repo-utils/commit"
)

func init() {
	t := &TestingStub{}
	mustFileExists(t, "/home/vagrant")
}

type TestingStub struct{}

func (t *TestingStub) Errorf(s string, a ...interface{}) {
	log.Fatalf(s+"\n", a...)
}

type TestingExiter struct{ t *testing.T }

func (t *TestingExiter) Errorf(s string, a ...interface{}) {
	panic(
		fmt.Errorf(s, a...),
	)
}

type Errorer interface {
	Errorf(string, ...interface{})
}

func TestGit(t *testing.T) {
	tt := &TestingExiter{t}
	DoTestFolderUnderVcs("/home/vagrant/git", tt)
	DoTestFolderUnderVcsAsJSON("/home/vagrant/git", tt)
	DoTestFolderUnderVcsAny("/home/vagrant/git", tt)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/git", tt)
	DoTestFolderIsClean("/home/vagrant/git", tt)
	DoTestFolderIsCleanJSON("/home/vagrant/git", tt)
	DoTestFolderIsDirty("/home/vagrant/git_dirty", tt)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/git_untracked", tt)
	DoCreateTag("/home/vagrant/git", tt)
	DoCreateTagWithMessage("/home/vagrant/git", tt)
	DoFailCreateTag("/home/vagrant/git", tt)
	DoFailCreateTagMissTagName("/home/vagrant/git", tt)
	DoListTags("/home/vagrant/git", tt)
	DoListCommits("/home/vagrant/git", tt)
	DoListCommitsBetween("/home/vagrant/git", tt)
	DoListCommitsSinceBeginning("/home/vagrant/git", tt)
	DoSortCommitsDesc("/home/vagrant/git", tt)
	DoTestFirstRevGit("/home/vagrant/git", tt)
}

func TestHg(t *testing.T) {
	tt := &TestingExiter{t}
	DoTestFolderUnderVcs("/home/vagrant/hg", tt)
	DoTestFolderUnderVcsAsJSON("/home/vagrant/hg", tt)
	DoTestFolderUnderVcsAny("/home/vagrant/hg", tt)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/hg", tt)
	DoTestFolderIsClean("/home/vagrant/hg", tt)
	DoTestFolderIsCleanJSON("/home/vagrant/hg", tt)
	DoTestFolderIsDirty("/home/vagrant/hg_dirty", tt)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/hg_untracked", tt)
	DoCreateTag("/home/vagrant/hg", tt)
	DoCreateTagWithMessage("/home/vagrant/hg", tt)
	DoFailCreateTag("/home/vagrant/hg", tt)
	DoFailCreateTagMissTagName("/home/vagrant/hg", tt)
	DoListTags("/home/vagrant/hg", tt)
	DoListCommits("/home/vagrant/hg", tt)
	DoListCommitsBetween("/home/vagrant/hg", tt)
	DoListCommitsSinceBeginning("/home/vagrant/hg", tt)
	DoSortCommitsDesc("/home/vagrant/hg", tt)
	DoTestFirstRevHg("/home/vagrant/hg", tt)
}

func TestSvn(t *testing.T) {
	tt := &TestingExiter{t}
	DoTestFolderUnderVcs("/home/vagrant/svn_work", tt)
	DoTestFolderUnderVcsAsJSON("/home/vagrant/svn_work", tt)
	DoTestFolderUnderVcsAny("/home/vagrant/svn_work", tt)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/svn_work", tt)
	DoTestFolderIsClean("/home/vagrant/svn_work", tt)
	DoTestFolderIsCleanJSON("/home/vagrant/svn_work", tt)
	DoTestFolderIsDirty("/home/vagrant/svn_dirty_work", tt)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/svn_untracked_work", tt)
	DoCreateTag("/home/vagrant/svn_work", tt)
	DoCreateTagWithMessage("/home/vagrant/svn_work", tt)
	DoFailCreateTag("/home/vagrant/svn_work", tt)
	DoFailCreateTagMissTagName("/home/vagrant/svn_work", tt)
	DoListTags("/home/vagrant/svn_work", tt)
	DoListCommits("/home/vagrant/svn_work", tt)
	DoListCommitsBetween("/home/vagrant/svn_work", tt)
	DoListCommitsSinceBeginning("/home/vagrant/svn_work", tt)
	DoSortCommitsDesc("/home/vagrant/svn_work", tt)
	DoTestFirstRevSvn("/home/vagrant/svn_work", tt)
}

func TestBzr(t *testing.T) {
	tt := &TestingExiter{t}
	DoTestFolderUnderVcs("/home/vagrant/bzr", tt)
	DoTestFolderUnderVcsAsJSON("/home/vagrant/bzr", tt)
	DoTestFolderUnderVcsAny("/home/vagrant/bzr", tt)
	DoTestFolderUnderVcsAnyReversed("/home/vagrant/bzr", tt)
	DoTestFolderIsClean("/home/vagrant/bzr", tt)
	DoTestFolderIsCleanJSON("/home/vagrant/bzr", tt)
	DoTestFolderIsDirty("/home/vagrant/bzr_dirty", tt)
	DoTestFolderIsCleanEvenWithUntrackedFiles("/home/vagrant/bzr_untracked", tt)
	DoCreateTag("/home/vagrant/bzr", tt)
	DoCreateTagWithMessage("/home/vagrant/bzr", tt)
	DoFailCreateTag("/home/vagrant/bzr", tt)
	DoFailCreateTagMissTagName("/home/vagrant/bzr", tt)
	DoListTags("/home/vagrant/bzr", tt)
	DoListCommits("/home/vagrant/bzr", tt)
	DoListCommitsBetween("/home/vagrant/bzr", tt)
	DoListCommitsSinceBeginning("/home/vagrant/bzr", tt)
	DoSortCommitsDesc("/home/vagrant/bzr", tt)
	DoTestFirstRevBzr("/home/vagrant/bzr", tt)
}

func TestPathArgs(t *testing.T) {
	tt := &TestingExiter{t}
	DoTestFolderUnderVcsWithPath("/home/vagrant/git", tt)
	DoTestFolderIsCleanWithPath("/home/vagrant/bzr", t)
}

func TestFolderNotUnderVcs(t *testing.T) {
	args := []string{"list-tags"}
	cmd := exec.Command("/vagrant/build/go-repo-utils", args...)
	cmd.Dir = "/home/vagrant"
	fmt.Printf("%s: %s %s\n", cmd.Dir, "/vagrant/build/go-repo-utils", args)

	err := cmd.Run()
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%t\n", true)
	}
}

func ExecSuccessCommand(t Errorer, cmd string, cwd string, args []string) string {
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
		t.Errorf("Expected success=true, got success=%t\n", false)
		return ""
	}

	return string(out)
}

func DoTestFolderUnderVcsWithPath(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-p", path}
	out := ExecSuccessCommand(t, cmd, "/home", args)

	expectedOut := "1.0.0\n1.0.2\n1.0.3\n1.0.4\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsCleanWithPath(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"is-clean", "--path=" + path}
	out := ExecSuccessCommand(t, cmd, "/home", args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoCreateTagWithMessage(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"create-tag", "1.0.4", "-m", "new tag"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "done\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoListTags(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "1.0.0\n1.0.2\n1.0.3\n1.0.4\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoCreateTag(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"create-tag", "1.0.3"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "done\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoFailCreateTag(path string, t Errorer) {
	args := []string{"create-tag", "1.0.3"}
	cmd := exec.Command("/vagrant/build/go-repo-utils", args...)
	cmd.Dir = path
	fmt.Printf("%s: %s %s\n", path, "/vagrant/build/go-repo-utils", args)

	err := cmd.Run()
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%t\n", true)
	}
}

func DoFailCreateTagMissTagName(path string, t Errorer) {
	args := []string{"create-tag"}
	cmd := exec.Command("/vagrant/build/go-repo-utils", args...)
	cmd.Dir = path
	fmt.Printf("%s: %s %s\n", path, "/vagrant/build/go-repo-utils", args)

	err := cmd.Run()
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%t\n", true)
	}
}

func DoTestFolderUnderVcs(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "1.0.0\n1.0.2\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderUnderVcsAsJSON(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-j"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "[\"1.0.0\",\"1.0.2\"]"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderUnderVcsAny(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-a"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "1.0.0\n1.0.2\nnotsemvertag\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderUnderVcsAnyReversed(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-a", "-r"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "notsemvertag\n1.0.2\n1.0.0\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsClean(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"is-clean"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsCleanJSON(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"is-clean", "-j"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "true"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsDirty(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"is-clean"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "no\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}
func DoTestFolderIsCleanEvenWithUntrackedFiles(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"is-clean"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoListCommits(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-commits", "--since", "notsemvertag"}
	out := ExecSuccessCommand(t, cmd, path, args)

	fmt.Println(string(out))

	var commits []commit.Commit
	err := json.Unmarshal([]byte(out), &commits)
	if err != nil {
		t.Errorf("Expected err=nil, got err=%q\n", err)
	}

	if len(commits) == 0 {
		t.Errorf("Expected to have commits")
	}

	found := false
	message := "tomate 1.0.2"
	for _, c := range commits {
		if c.Message == message {
			found = true
		}
	}

	if found == false {
		t.Errorf("Expected commits to contain an entry with message=%q, but it was not found\n", message)
	}
}

func DoListCommitsBetween(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-commits", "--since", "v1.0.2", "--until", "v1.0.0"}
	out := ExecSuccessCommand(t, cmd, path, args)

	fmt.Println(string(out))

	var commits []commit.Commit
	err := json.Unmarshal([]byte(out), &commits)
	if err != nil {
		t.Errorf("Expected err=nil, got err=%q\n", err)
	}

	if len(commits) == 0 {
		t.Errorf("Expected to have commits")
	}

	found := false
	message := "tomate 1.0.0"
	for _, c := range commits {
		if c.Message == message {
			found = true
		}
	}

	if found == false {
		t.Errorf("Expected commits to contain an entry with message=%q, but it was not found\n", message)
	}
}

func DoListCommitsSinceBeginning(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-commits", "--until", "v1.0.0"}
	out := ExecSuccessCommand(t, cmd, path, args)

	var commits []commit.Commit
	err := json.Unmarshal([]byte(out), &commits)
	if err != nil {
		t.Errorf("Expected err=nil, got err=%q\n", err)
		fmt.Println(string(out))
	}

	if len(commits) == 0 {
		t.Errorf("Expected to have commits")
	}

	found := false
	message := "tomate notsemvertag"
	for _, c := range commits {
		if c.Message == message {
			found = true
		}
	}

	if found == false {
		t.Errorf("Expected commits to contain an entry with message=%q, but it was not found\n", message)
	}
}

func DoSortCommitsDesc(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-commits"}
	out := ExecSuccessCommand(t, cmd, path, args)

	var commits commit.Commits
	err := json.Unmarshal([]byte(out), &commits)
	if err != nil {
		t.Errorf("Expected err=nil, got err=%q\n", err)
		fmt.Println(string(out))
	}
	if len(commits) == 0 {
		t.Errorf("Expected to have commits")
	}

	commits.OrderByDate("DESC")

	first := commits[0].GetDate()
	last := commits[len(commits)-1].GetDate()

	if first.After(*last) == false {
		t.Errorf("Expected commits to be ordered DESC, they are not\n")
	}
}

func DoTestFirstRevGit(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "d6b486e435f8497b1b873ce8a1e0fafbf82fed0e\n"
	if out != expectedOut {
		// t.Errorf("Expected=%q, got out=%q\n", expectedOut, out)
		// git can t be tested, the hash changes at every test session
	}
}

func DoTestFirstRevHg(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "065e4375921ce712e536b95109214b28e8e2c23e\n"
	if out != expectedOut {
		// t.Errorf("Expected=%q, got out=%q\n", expectedOut, out)
		// hg can t be tested, the hash changes at every test session
	}
}

func DoTestFirstRevBzr(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "revno:1\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFirstRevSvn(path string, t Errorer) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "1\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}
func mustFileExists(t Errorer, p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Errorf("file mut exists %q", p)
		return false
	}
	return true
}
