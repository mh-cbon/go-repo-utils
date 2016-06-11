# go-repo-utils

Go tool to speak with repositories.

It can list tags, tell if a directory is clean, create tag.

It can speak with `hg` `git` `bzr` `svn`

# Install

```sh
glide get github.com/mh-cbon/go-repo-utils
cd $GOPATH/github.com/mh-cbon/go-repo-utils
go install .
```

# Usage

```sh
Go repo utils

Usage:
  go-repo-utils list-tags [-j|--json] [-a|--any] [-r|--reverse] [--path=<path>|-p=<path>]
  go-repo-utils is-clean [-j|--json] [--path=<path>|-p=<path>]
  go-repo-utils create-tag <tag> [-j|--json] [--path=<path>|-p=<path>]
  go-repo-utils -h | --help
  go-repo-utils -v | --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -p=<c> --path=<c>     Path to lookup [default: cwd].
  -j --json             Print JSON encoded data.
  -a --any              List all tags.
  -r --reverse          Reverse tags ordering.

Notes:
  list-tags   List only valid semver tags unless -a|--any options is provided.
  is-clean    Ignores untracked files.
  create-tag  With svn, it always create a new tag folder at /tags/<tag>.
```

#### Examples

```sh
# list tags
go-repo-utils list-tags
# list tags with json response
go-repo-utils list-tags -j --path=/some/where
# check if a directory is clean
go-repo-utis is-clean -p /some/where
# create tag
go-repo-utils create-tag 1.0.3
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

# Tests

To run the tests, `sh vagrant/test.sh`, which will do all necessary stuff to run the tests
