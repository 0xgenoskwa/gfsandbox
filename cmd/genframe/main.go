package main

import (
	"fmt"

	"go.genframe.xyz/internal/genframe"
)

func main() {
	// InitializeServer
	g := genframe.InitializeGenframe()
	fmt.Println("genframe run 1")
	g.Run()
	fmt.Println("genframe run 2")
}
