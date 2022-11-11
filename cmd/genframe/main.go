package main

import (
	"go.genframe.xyz/internal/genframe"
)

func main() {
	// InitializeServer
	g := genframe.InitializeGenframe()
	g.Run()
}
