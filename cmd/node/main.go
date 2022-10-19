package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/block-api/block-node/block"
	"github.com/block-api/block-node/log"
)

func main() {
	log.Default("block-node example")

	options := block.BlockNodeOptions{
		Name:    "block-node-service",
		Version: 1,
	}

	blockNode := block.NewBlockNode(&options)

	// authBlock := AuthBlock{
	// 	block.Block{
	// 		Name: "auth-block",
	// 	},
	// }

	// blockNode.AddBlock(&authBlock)
	blockNode.Start()

	var osSignal chan os.Signal = make(chan os.Signal)

	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-osSignal:
			log.Warning("shutting down, please wait")
			os.Exit(0)
		}
	}
}

type AuthBlock struct {
	block.Block
}

func (ab *AuthBlock) GetName() string {
	return ab.Block.Name
}

// func (ab *AuthBlock) GetName() string {
// 	return ab.Block.Name
// }
