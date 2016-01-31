package main

import "github.com/nildev/tools/Godeps/_workspace/src/github.com/codeskyblue/go-sh"

func main() {
	sh.Command("less", "less.go").Run()
}
