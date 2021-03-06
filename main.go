// Package go-repo_utils helps to work with VCS.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt.go"
	"github.com/mh-cbon/go-repo-utils/commit"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

// VERSION contains the last build version.
var VERSION = "0.0.0"

func main() {
	usage := `Go repo utils

Usage:
  go-repo-utils list-tags [-j|--json] [-a|--any] [-r|--reverse] [--path=<path>|-p <path>]
  go-repo-utils list-commits [--path=<path>|-p <path>] [--since=<tag>|-s <tag>] [--until=<tag>|-u <tag>] [-r|--reverse] [--orderbydate]
  go-repo-utils is-clean [-j|--json] [--path=<path>|-p=<path>]
  go-repo-utils create-tag <tag> [-j|--json] [--path=<path>|-p <path>] [-m <message>]
  go-repo-utils first-rev [-j|--json] [--path=<path>|-p <path>]
  go-repo-utils -h | --help
  go-repo-utils -v | --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -p <c> --path=<c>     Path to lookup [default: cwd].
  -s <c> --since=<c>    Since tag, revision, expression.
  -u <c> --until=<c>    To tag, revision, expression.
  -j --json             Print JSON encoded data.
  -a --any              List all tags.
  -r --reverse          Reverse tags ordering.
  -m                    Message for the tag.
  --orderbydate         Order commits by date.

Notes:
  list-tags     List only valid semver tags unless -a|--any options is provided.
  is-clean      Ignores untracked files.
  create-tag    With svn, it always create a new tag folder at /tags/<tag>.
  list-commits  Can receive an expression (hg, bzr), if it does not match a tag name.
                Expression may be automatically adjusted at runtime if it is empty (svn,hg,bzr),
                or matching a tag name.
                HEAD will be normalized given the target vcs (svn,hg,bzr).

Examples
  # list tags
  go-repo-utils list-tags

  # list tags with json response
  go-repo-utils list-tags -j --path=/some/where

  # check if a directory is clean
  go-repo-utis is-clean -p /some/where

  # create tag
  go-repo-utils create-tag 1.0.3 -m "tag message"
`

	arguments, err := docopt.Parse(usage, nil, true, "Go repo utils - "+VERSION, false)

	logger.Println(arguments)
	exitWithError(err)

	cmd := getCommand(arguments)

	path := getPath(arguments)
	if path == "" {
		path, err = os.Getwd()
		exitWithError(err)
	}

	vcs, err := repoutils.WhichVcs(path)
	exitWithError(err)

	if cmd == "list-tags" {
		cmdListTags(arguments, vcs, path)
	} else if cmd == "list-commits" {
		cmdListCommits(arguments, vcs, path)
	} else if cmd == "is-clean" {
		cmdIsClean(arguments, vcs, path)
	} else if cmd == "create-tag" {
		cmdCreateTag(arguments, vcs, path)
	} else if cmd == "first-rev" {
		cmdFirstRev(arguments, vcs, path)
	} else if cmd == "" {
		fmt.Println("Wrong usage: Missing command")
		fmt.Println("")
		fmt.Println(usage)
		os.Exit(1)
	} else {
		log.Println("Unknown command: '" + cmd + "'")
		os.Exit(1)
	}
}

func cmdIsClean(arguments map[string]interface{}, vcs string, path string) {
	isClean, err := repoutils.IsClean(vcs, path)
	exitWithError(err)

	if isJSON(arguments) {
		jsoned, _ := json.Marshal(isClean)
		fmt.Print(string(jsoned))
	} else {
		if isClean {
			fmt.Println("yes")
		} else {
			fmt.Println("no")
		}
	}
}

func cmdListTags(arguments map[string]interface{}, vcs string, path string) {
	tags := make([]string, 0)
	dirtyTags, err := repoutils.List(vcs, path)
	exitWithError(err)

	if isAny(arguments) == false {
		tags = repoutils.FilterSemverTags(dirtyTags)
	} else {
		tags = append(tags, dirtyTags...)
	}

	tags = repoutils.SortSemverTags(tags)

	if isReversed(arguments) {
		tags = repoutils.ReverseTags(tags)
	}

	if isJSON(arguments) {
		jsoned, _ := json.Marshal(tags)
		fmt.Print(string(jsoned))
	} else {
		for _, tag := range tags {
			if len(tag) > 0 {
				fmt.Println(tag)
			}
		}
	}
}

func cmdListCommits(arguments map[string]interface{}, vcs string, path string) {

	since := getSince(arguments)
	until := getUntil(arguments)
	reversed := isReversed(arguments)
	orderbydate := isOrderByDate(arguments)

	if len(until) == 0 {
		until = "HEAD"
	}

	commits, err := repoutils.ListCommitsBetween(vcs, path, since, until)
	exitWithError(err)

	if orderbydate {
		if reversed {
			commit.Commits(commits).OrderByDate("DESC")
		} else {
			commit.Commits(commits).OrderByDate("ASC")
		}
	} else if reversed {
		commit.Commits(commits).Reverse()
	}

	jsoned, err := json.MarshalIndent(commits, "", "  ")
	exitWithError(err)
	fmt.Print(string(jsoned))
}

func cmdCreateTag(arguments map[string]interface{}, vcs string, path string) {

	tag := getTag(arguments)
	if len(tag) == 0 {
		exitWithError(errors.New("Missing tag value"))
	}
	message := getMessage(arguments)
	if len(message) == 0 {
		message = "tag: " + tag
	}

	_, out, err := repoutils.CreateTag(vcs, path, tag, message)
	if err != nil {
		log.Println(out)
		exitWithError(err)
	}

	if isJSON(arguments) {
		jsoned, _ := json.Marshal(true)
		fmt.Print(string(jsoned))
	} else {
		fmt.Println("done")
	}
}

func cmdFirstRev(arguments map[string]interface{}, vcs string, path string) {

	out, err := repoutils.GetFirstRevision(vcs, path)
	if err != nil {
		log.Println(out)
		exitWithError(err)
	}

	if isJSON(arguments) {
		jsoned, _ := json.Marshal(true)
		fmt.Print(string(jsoned))
	} else {
		fmt.Println(string(out))
	}
}

func getCommand(arguments map[string]interface{}) string {
	cmds := []string{
		"list-tags",
		"is-clean",
		"create-tag",
		"list-commits",
		"first-rev",
	}
	for _, cmd := range cmds {
		if p, ok := arguments[cmd]; ok {
			if b, ok := p.(bool); ok && b {
				return cmd
			}
		}
	}
	return ""
}

func getPath(arguments map[string]interface{}) string {
	args := []string{
		"--path",
		"-p",
	}
	for _, arg := range args {
		if p, ok := arguments[arg]; ok {
			if b, ok2 := p.(string); ok2 {
				if b != "cwd" {
					return b
				}
			}
		}
	}
	return ""
}

func getTag(arguments map[string]interface{}) string {
	tag := ""
	if t, ok := arguments["<tag>"].(string); ok {
		tag = t
	} else if t, ok := arguments["--tag"].(string); ok {
		tag = t
	}
	return tag
}

func getSince(arguments map[string]interface{}) string {
	tag := ""
	if t, ok := arguments["--since"].(string); ok {
		tag = t
	}
	return tag
}

func getUntil(arguments map[string]interface{}) string {
	tag := ""
	if t, ok := arguments["--until"].(string); ok {
		tag = t
	}
	return tag
}

func getMessage(arguments map[string]interface{}) string {
	message := ""
	if mess, ok := arguments["-m"].(string); ok {
		message = mess
	}
	return message
}

func isAny(arguments map[string]interface{}) bool {
	any := false
	if isAny, ok := arguments["--any"].(bool); ok {
		any = isAny
	} else {
		if isA, ok := arguments["-a"].(bool); ok {
			any = isA
		}
	}
	return any
}

func isJSON(arguments map[string]interface{}) bool {
	json := false
	if isIt, ok := arguments["--json"].(bool); ok {
		json = isIt
	} else {
		if isJ, ok := arguments["-j"].(bool); ok {
			json = isJ
		}
	}
	return json
}

func isReversed(arguments map[string]interface{}) bool {
	reverse := false
	if isReverse, ok := arguments["--reverse"].(bool); ok {
		reverse = isReverse
	} else {
		if isR, ok := arguments["-r"].(bool); ok {
			reverse = isR
		}
	}
	return reverse
}

func isOrderByDate(arguments map[string]interface{}) bool {
	orderbydate := false
	if isIt, ok := arguments["--orderbydate"].(bool); ok {
		orderbydate = isIt
	}
	return orderbydate
}

func exitWithError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
