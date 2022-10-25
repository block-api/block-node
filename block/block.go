package block

import "github.com/block-api/block-node/utils"

type IBlock interface {
	GetName() utils.BlockName
	Actions() map[string]BlockAction
}

type BlockAction func(payload interface{}) error

type Block struct {
	IBlock
	Name    utils.BlockName
	Actions map[string]BlockAction
}

func (b *Block) GetName() utils.BlockName {
	return b.Name
}
