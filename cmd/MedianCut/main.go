package main

import (
	_ "image/gif"
	_ "image/jpeg"

	mediancut "github.com/joyme123/WebAssembly-MedianCut"
)

func main() {

	var api mediancut.WebAPI

	api.Init()
}
