package main

import (
	"log"
	"os"

	"github.com/codenotary/cas/pkg/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = doc.GenMarkdownTree(cmd.Root(), pwd+"/docs/cmd")
	if err != nil {
		log.Fatal(err)
	}
}
