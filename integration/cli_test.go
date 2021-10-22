package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

/**
Main test, meant to compile the binary. If test fails everything else will be ignored
**/
func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("Error changing directory, aborting tests")
		os.Exit(1)
	}
	build := exec.Command("make")
	err = build.Run()
	if err != nil {
		fmt.Printf("Error compiling main binary, aborting tests")
		os.Exit(1)
	}

}
