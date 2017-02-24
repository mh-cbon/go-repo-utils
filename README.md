# go-repo-utils

[![travis Status](https://travis-ci.org/mh-cbon/go-repo-utils.svg?branch=master)](https://travis-ci.org/mh-cbon/go-repo-utils)[![appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/go-repo-utils?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/go-repo-utils)
[![GoDoc](https://godoc.org/github.com/mh-cbon/go-repo-utils?status.svg)](http://godoc.org/github.com/mh-cbon/go-repo-utils)


Package go-repo_utils helps to work with VCS.

It can list tags, tell if a directory is clean, create tag.

It can speak with `hg` `git` `bzr` `svn`

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# Install

Check the [release page](https://github.com/mh-cbon/go-repo-utils/releases)!

#### Glide

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/go-repo-utils
cd $GOPATH/src/github.com/mh-cbon/go-repo-utils
git clone https://github.com/mh-cbon/go-repo-utils.git .
glide install
go install
```


#### Chocolatey
```sh
choco install go-repo-utils
```

#### linux rpm/deb repository
```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/go-repo-utils sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/go-repo-utils sh -xe
```

#### linux rpm/deb standalone package
```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-repo-utils sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-repo-utils sh -xe
```

# Usage

__$ go-repo-utils -h__
```sh
Go repo utils

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
```

#### Enable debug messages

To enable debug messages, just set `VERBOSE=go-repo-utils` before running the command.

```sh
VERBOSE=go-repo-utils go-repo-utils is-clean
```

# Usage as lib

__> main_example.go__
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mh-cbon/go-repo-utils/repoutils"
)

// ExampleMain demonstrate go-repo-utils api
func ExampleMain() {

	path := "path/to/folder"

	vcs, err := repoutils.WhichVcs(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	tags, _ := repoutils.List(vcs, path)
	fmt.Println(tags)

	isClean, _ := repoutils.IsClean(vcs, path)
	fmt.Println(isClean)

	ok, _, _ := repoutils.CreateTag(vcs, path, "1.0.3", "the new tag")
	fmt.Println(ok)
}
```

# Tests

To run the tests, `sh vagrant/test.sh`, which will do all necessary stuff to run the tests

# See also

- https://github.com/Masterminds/vcs

A way more complete and better api, with a different approach.
