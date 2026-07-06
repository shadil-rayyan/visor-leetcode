package main

import "os"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--json" {
		scrapeJsonMain()
		return
	}
	println("usage: go run . --json")
}
