package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"

	"github.com/mh-cbon/go-repo-utils/commit"
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
	DoCreateTagWithMessage("/home/vagrant/git", t)
	DoFailCreateTag("/home/vagrant/git", t)
	DoFailCreateTagMissTagName("/home/vagrant/git", t)
	DoListTags("/home/vagrant/git", t)
	DoListCommits("/home/vagrant/git", t)
	DoListCommitsBetween("/home/vagrant/git", t)
	DoListCommitsSinceBeginning("/home/vagrant/git", t)
	DoSortCommitsDesc("/home/vagrant/git", t)
	DoTestFirstRevGit("/home/vagrant/git", t)
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
	DoCreateTagWithMessage("/home/vagrant/hg", t)
	DoFailCreateTag("/home/vagrant/hg", t)
	DoFailCreateTagMissTagName("/home/vagrant/hg", t)
	DoListTags("/home/vagrant/hg", t)
	DoListCommits("/home/vagrant/hg", t)
	DoListCommitsBetween("/home/vagrant/hg", t)
	DoListCommitsSinceBeginning("/home/vagrant/hg", t)
	DoSortCommitsDesc("/home/vagrant/hg", t)
	DoTestFirstRevHg("/home/vagrant/hg", t)
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
	DoCreateTag("/home/vagrant/svn_work", t)
	DoCreateTagWithMessage("/home/vagrant/svn_work", t)
	DoFailCreateTag("/home/vagrant/svn_work", t)
	DoFailCreateTagMissTagName("/home/vagrant/svn_work", t)
	DoListTags("/home/vagrant/svn_work", t)
	DoListCommits("/home/vagrant/svn_work", t)
	DoListCommitsBetween("/home/vagrant/svn_work", t)
	DoListCommitsSinceBeginning("/home/vagrant/svn_work", t)
	DoSortCommitsDesc("/home/vagrant/svn_work", t)
	DoTestFirstRevSvn("/home/vagrant/svn_work", t)
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
	DoCreateTagWithMessage("/home/vagrant/bzr", t)
	DoFailCreateTag("/home/vagrant/bzr", t)
	DoFailCreateTagMissTagName("/home/vagrant/bzr", t)
	DoListTags("/home/vagrant/bzr", t)
	DoListCommits("/home/vagrant/bzr", t)
	DoListCommitsBetween("/home/vagrant/bzr", t)
	DoListCommitsSinceBeginning("/home/vagrant/bzr", t)
	DoSortCommitsDesc("/home/vagrant/bzr", t)
	DoTestFirstRevBzr("/home/vagrant/bzr", t)
}

func TestPathArgs(t *testing.T) {
	DoTestFolderUnderVcsWithPath("/home/vagrant/git", t)
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
		t.Errorf("Expected success=true, got success=%t\n", false)
		return ""
	}

	return string(out)
}

func DoTestFolderUnderVcsWithPath(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags", "-p", path}
	out := ExecSuccessCommand(t, cmd, "/home", args)

	expectedOut := "1.0.0\n1.0.2\n1.0.3\n1.0.4\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFolderIsCleanWithPath(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"is-clean", "--path=" + path}
	out := ExecSuccessCommand(t, cmd, "/home", args)

	expectedOut := "yes\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoCreateTagWithMessage(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"create-tag", "1.0.4", "-m", "new tag"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "done\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoListTags(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"list-tags"}
	out := ExecSuccessCommand(t, cmd, path, args)

	expectedOut := "1.0.0\n1.0.2\n1.0.3\n1.0.4\n"
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
	fmt.Printf("%s: %s %s\n", path, "/vagrant/build/go-repo-utils", args)

	err := cmd.Run()
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%t\n", true)
	}
}

func DoFailCreateTagMissTagName(path string, t *testing.T) {
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

func DoListCommits(path string, t *testing.T) {
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

func DoListCommitsBetween(path string, t *testing.T) {
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

func DoListCommitsSinceBeginning(path string, t *testing.T) {
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

func DoSortCommitsDesc(path string, t *testing.T) {
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

func DoTestFirstRevGit(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "d6b486e435f8497b1b873ce8a1e0fafbf82fed0e\n"
	if out != expectedOut {
		// t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
    // git can t be tested, the hash changes at every test session
	}
}

func DoTestFirstRevHg(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "065e4375921ce712e536b95109214b28e8e2c23e\n"
	if out != expectedOut {
		// t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
    // hg can t be tested, the hash changes at every test session
	}
}

func DoTestFirstRevBzr(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "revno:1\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}

func DoTestFirstRevSvn(path string, t *testing.T) {
	cmd := "/vagrant/build/go-repo-utils"
	args := []string{"first-rev"}
	out := ExecSuccessCommand(t, cmd, path, args)
	expectedOut := "1\n"
	if out != expectedOut {
		t.Errorf("Expected out=%q, got out=%q\n", expectedOut, out)
	}
}
