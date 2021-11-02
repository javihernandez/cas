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
	log.Printf("Generating Markdown pages")
	err = doc.GenMarkdownTree(cmd.Root(), pwd+"/docs/cmd")
	if err != nil {
		log.Fatal(err)
	}

	header := &doc.GenManHeader{
		Title: "CAS",
		Section: "1",
	}
	
	log.Printf("Generating man pages")
	err = doc.GenManTree(cmd.Root(), header, pwd+"/docs/man")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done")
}
