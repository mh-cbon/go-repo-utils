package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt.go"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func main() {
	usage := `Go repo utils

Usage:
  go-repo-utils list-tags [-j|--json] [-a|--any] [-r|--reverse] [--path=<path>|-p <path>]
  go-repo-utils is-clean [-j|--json] [--path=<path>|-p=<path>]
  go-repo-utils create-tag <tag> [-j|--json] [--path=<path>|-p <path>] [-m <message>]
  go-repo-utils -h | --help
  go-repo-utils -v | --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -p <c> --path=<c>     Path to lookup [default: cwd].
  -j --json             Print JSON encoded data.
  -a --any              List all tags.
  -r --reverse          Reverse tags ordering.
  -m                    Message for the tag.

Notes:
  list-tags   List only valid semver tags unless -a|--any options is provided.
  is-clean    Ignores untracked files.
  create-tag  With svn, it always create a new tag folder at /tags/<tag>.

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

	arguments, err := docopt.Parse(usage, nil, true, "Go repo utils", false)

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
	} else if cmd == "is-clean" {
		cmdIsClean(arguments, vcs, path)
	} else if cmd == "create-tag" {
		cmdCreateTag(arguments, vcs, path)
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

	if isJson(arguments) {
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

	if isJson(arguments) {
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

	if isJson(arguments) {
		jsoned, _ := json.Marshal(true)
		fmt.Print(string(jsoned))
	} else {
		fmt.Println("done")
	}
}

func getCommand(arguments map[string]interface{}) string {
	p, ok := arguments["list-tags"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "list-tags"
		}
	}
	p, ok = arguments["is-clean"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "is-clean"
		}
	}
	p, ok = arguments["create-tag"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "create-tag"
		}
	}
	return ""
}

func getPath(arguments map[string]interface{}) string {
	p, ok := arguments["--path"]
	if ok {
		if str, ok := p.(string); ok {
			if str != "cwd" {
				return str
			}
		}
	}
	p, ok = arguments["-p"]
	if ok {
		if str, ok := p.(string); ok {
			if str != "cwd" {
				return str
			}
		}
	}
	return ""
}

func getTag(arguments map[string]interface{}) string {
	tag := ""
	if t, ok := arguments["<tag>"].(string); ok {
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

func isJson(arguments map[string]interface{}) bool {
	json := false
	if isJson, ok := arguments["--json"].(bool); ok {
		json = isJson
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

func exitWithError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
