# go-repo-utils

Go tool to speak with repositories.

It can list tags, tell if a directory is clean, create tag.

It can speak with `hg` `git` `bzr` `svn`

# Install

You can grab a pre-built binary file in the [releases page](https://github.com/mh-cbon/go-repo-utils/releases)

```sh
mkdir -p $GOPATH/github.com/mh-cbon
cd $GOPATH/github.com/mh-cbon
git clone https://github.com/mh-cbon/go-repo-utils.git
cd go-repo-utils
glide install
go install
```

# Usage

```sh
Go repo utils

Usage:
  go-repo-utils list-tags [-j|--json] [-a|--any] [-r|--reverse] [--path=<path>|-p <path>]
  go-repo-utils is-clean [-j|--json] [--path=<path>|-p=<path>]
  go-repo-utils list-commits [--path=<path>|-p <path>] [--since=<tag>|-s <tag>] [--until=<tag>|-u <tag>] [-r|--reverse] [--orderbydate]
  go-repo-utils create-tag <tag> [-j|--json] [--path=<path>|-p <path>] [-m <message>]
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

```go
package main

import (
  "fmt"

  "github.com/mh-cbon/go-repo-utils/repoutils"
)

func main() {

  path := "path/to/folder"

  vcs, err := repoutils.WhichVcs(path)
  if err!=nil {
    log.Println(err)
    os.Exit(1)
  }

  tags := make([]string, 0)
  tags, _ := repoutils.List(vcs, path)
  fmt.Println(tags)

  isClean, _ := repoutils.IsClean(vcs, path)
  fmt.Println(isClean)

  ok, _, _ := repoutils.CreateTag(vcs, path, "1.0.3")
  fmt.Println(ok)
}

```

# Changelog

- 0.0.7 - Add list-commits command,
fixed bugs in vcs detection when the tree contained uncommitted files,
code refactoring and updated tests
- 0.0.6 - Use annotated tags for git (`git tag -a ...`)
- 0.0.5 - Release scripts
- 0.0.4 - Minor improvements
- 0.0.2 - Add '-m' argument to set a tag message
- 0.0.1 - initial release

# Tests

To run the tests, `sh vagrant/test.sh`, which will do all necessary stuff to run the tests

# See also

- https://github.com/Masterminds/vcs

A way more complete and better api, with a different approach.
