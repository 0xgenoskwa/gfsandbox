package main

import "go.genframe.xyz/internal/bootstrap"

func main() {
	// InitializeServer
	b := bootstrap.InitializeBootstrap()
	b.Run()
}
