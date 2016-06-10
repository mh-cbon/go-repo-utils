package repoutils

import (
	"errors"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/mh-cbon/go-repo-utils/bzr"
	"github.com/mh-cbon/go-repo-utils/git"
	"github.com/mh-cbon/go-repo-utils/hg"
	"github.com/mh-cbon/go-repo-utils/svn"
)

type IsIt func(path string) bool
type ListIt func(path string) ([]string, error)
type IsItClean func(path string) (bool, error)
type DoCreateTag func(path string, tag string) (bool, error)

type isVcsResult struct {
	name  string
	found bool
}

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
		if len(vcsTests) == 4 {
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

func CreateTag(vcs string, path string, tag string) (bool, error) {
	vcsCreateTag := map[string]DoCreateTag{
		"git": git.CreateTag,
		"bzr": bzr.CreateTag,
		"hg":  hg.CreateTag,
		"svn": svn.CreateTag,
	}
	createTag, ok := vcsCreateTag[vcs]
	if ok == false {
		return false, errors.New("Unknown VCS '" + vcs + "'")
	}
	return createTag(path, tag)
}

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

func ReverseTags(tags []string) []string {
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}
	return tags
}
