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
