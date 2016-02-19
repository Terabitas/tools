package main

import "github.com/nildev/project/Godeps/_workspace/src/github.com/codeskyblue/go-sh"

func main() {
	sh.Command("less", "less.go").Run()
}
