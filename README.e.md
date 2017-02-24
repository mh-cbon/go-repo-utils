# {{.Name}}

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/godoc" .}}

{{pkgdoc}}
It can list tags, tell if a directory is clean, create tag.

It can speak with `hg` `git` `bzr` `svn`

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# Install

{{template "gh/releases" .}}

#### Glide
{{template "glide/install" .}}

#### Chocolatey
{{template "choco/install" .}}

#### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

# Usage
{{cli "go-repo-utils" "-h"}}

#### Enable debug messages

To enable debug messages, just set `VERBOSE=go-repo-utils` before running the command.

```sh
VERBOSE=go-repo-utils go-repo-utils is-clean
```

# Usage as lib
{{file "main_example.go"}}

# Tests

To run the tests, `sh vagrant/test.sh`, which will do all necessary stuff to run the tests

# See also

- https://github.com/Masterminds/vcs

A way more complete and better api, with a different approach.
