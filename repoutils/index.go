// proxy implemntation to a vcs implementation
package repoutils

import (
	"errors"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/mh-cbon/go-repo-utils/bzr"
	"github.com/mh-cbon/go-repo-utils/commit"
	"github.com/mh-cbon/go-repo-utils/git"
	"github.com/mh-cbon/go-repo-utils/hg"
	"github.com/mh-cbon/go-repo-utils/svn"
)

// Func type declarations
type IsIt func(path string) bool
type ListIt func(path string) ([]string, error)
type IsItClean func(path string) (bool, error)
type DoCreateTag func(path string, tag string, message string) (bool, string, error)
type DoAdd func(path string, file string) error
type DoCommit func(path string, message string, files []string) error
type DoListCommitsBetween func(path string, since string, to string) ([]commit.Commit, error)

type isVcsResult struct {
	name  string
	found bool
}

// Determine the kind of VCS of given path
func WhichVcs(path string) (string, error) {
	vcsTester := map[string]IsIt{
		"git": git.IsIt,
		"bzr": bzr.IsIt,
		"hg":  hg.IsIt,
		"svn": svn.IsIt,
	}
	vcsTests := map[string]bool{}

	out := make(chan isVcsResult, len(vcsTester))
	for vcs, isIt := range vcsTester {
		go func(vcs string, isIt IsIt) {
			out <- isVcsResult{name: vcs, found: isIt(path)}
		}(vcs, isIt)
	}
	for res := range out {
		vcsTests[res.name] = res.found
		if len(vcsTests) == len(vcsTester) {
			close(out)
		}
	}

	vcsFound := ""
	howMuchFound := 0
	for vcs, found := range vcsTests {
		if found {
			howMuchFound++
			vcsFound += vcs
		}
	}

	if howMuchFound == 0 {
		return "", errors.New("No vcs project found at '" + path + "'")
	}

	if howMuchFound > 1 {
		return "", errors.New("Multiple vcs project found at '" + path + "'. ?? => '" + vcsFound + "'")
	}

	return vcsFound, nil
}

// List tags on given path according to given vcs
func List(vcs string, path string) ([]string, error) {
	vcsLister := map[string]ListIt{
		"git": git.List,
		"bzr": bzr.List,
		"hg":  hg.List,
		"svn": svn.List,
	}
	lister, ok := vcsLister[vcs]
	if ok == false {
		return make([]string, 0), errors.New("Unknown VCS '" + vcs + "'")
	}
	return lister(path)
}

// Ensure given path does not contain uncommited files
func IsClean(vcs string, path string) (bool, error) {
	vcsIsClean := map[string]IsItClean{
		"git": git.IsClean,
		"bzr": bzr.IsClean,
		"hg":  hg.IsClean,
		"svn": svn.IsClean,
	}
	isItClean, ok := vcsIsClean[vcs]
	if ok == false {
		return false, errors.New("Unknown VCS '" + vcs + "'")
	}
	return isItClean(path)
}

// Create tag on given path
func CreateTag(vcs string, path string, tag string, message string) (bool, string, error) {
	vcsCreateTag := map[string]DoCreateTag{
		"git": git.CreateTag,
		"bzr": bzr.CreateTag,
		"hg":  hg.CreateTag,
		"svn": svn.CreateTag,
	}
	createTag, ok := vcsCreateTag[vcs]
	if ok == false {
		return false, "", errors.New("Unknown VCS '" + vcs + "'")
	}
	return createTag(path, tag, message)
}

// Add a file
func Add(vcs string, path string, file string) error {
	vcsAdd := map[string]DoAdd{
		"git": git.Add,
		"bzr": bzr.Add,
		"hg":  hg.Add,
		"svn": svn.Add,
	}
	doAdd, ok := vcsAdd[vcs]
	if ok == false {
		return errors.New("Unknown VCS '" + vcs + "'")
	}
	return doAdd(path, file)
}

// Commit files on path with message
func Commit(vcs string, path string, message string, files []string) error {
	vcsCommit := map[string]DoCommit{
		"git": git.Commit,
		"bzr": bzr.Commit,
		"hg":  hg.Commit,
		"svn": svn.Commit,
	}
	doCommit, ok := vcsCommit[vcs]
	if ok == false {
		return errors.New("Unknown VCS '" + vcs + "'")
	}
	return doCommit(path, message, files)
}

// Filter out invalid semver tags
func FilterSemverTags(dirtyTags []string) []string {
	tags := make([]string, 0)
	for _, tag := range dirtyTags {
		_, err := semver.NewVersion(tag)
		if err == nil {
			tags = append(tags, tag)
		}
	}
	return tags
}

// Sort given list of semver tags, invalid semver tags are appended to the end
func SortSemverTags(unsortedTags []string) []string {
	dirtyTags := make([]string, 0)
	vs := make([]*semver.Version, 0)
	for _, r := range unsortedTags {
		v, err := semver.NewVersion(r)
		if err != nil {
			dirtyTags = append(dirtyTags, r)
		} else {
			vs = append(vs, v)
		}
	}
	sort.Sort(semver.Collection(vs))
	sortedTags := make([]string, 0)
	for _, t := range vs {
		sortedTags = append(sortedTags, t.String())
	}
	sortedTags = append(sortedTags, dirtyTags...)
	return sortedTags
}

// Reverse given list of tags
func ReverseTags(tags []string) []string {
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}
	return tags
}

// List commits since given tag
func ListCommitsBetween(vcs string, path string, since string, to string) ([]commit.Commit, error) {
	ret := make([]commit.Commit, 0)
	listCommits := map[string]DoListCommitsBetween{
		"git": git.ListCommitsBetween,
		"bzr": bzr.ListCommitsBetween,
		"hg":  hg.ListCommitsBetween,
		"svn": svn.ListCommitsBetween,
	}
	doListCommits, ok := listCommits[vcs]
	if ok == false {
		return ret, errors.New("Unknown VCS '" + vcs + "'")
	}
	return doListCommits(path, since, to)
}
