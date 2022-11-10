package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/block-api/block-node/block"
	"github.com/block-api/block-node/log"
	"github.com/urfave/cli/v2"
)

func StartNode(cCtx *cli.Context) error {
	node, err := block.NewNode()
	if err != nil {
		return err
	}

	node.Start()

	var osSignal chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-osSignal:
			log.Warning("shutting down, please wait")
			node := block.GetNode()

			if node != nil {
				block.GetNode().Stop()
			}

			os.Exit(0)
		}
	}
}
