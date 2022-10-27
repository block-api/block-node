package block

import (
	"github.com/block-api/block-node/common/types"
)

type IBlock interface {
	GetName() types.BlockName
	Actions() map[types.ActionName]BlockAction
	ActionsNames() []types.ActionName
	AddAction(name types.ActionName, action BlockAction)
}

type BlockAction func(payload []byte) (any, error)

type Block struct {
	IBlock
	Name      types.BlockName
	blockNode *BlockNode
	actions   map[types.ActionName]BlockAction
}

func NewBlock(bn *BlockNode, name types.BlockName) Block {
	return Block{
		blockNode: bn,
		Name:      name,
		actions:   make(map[types.ActionName]BlockAction),
	}
}

func (b *Block) BlockNode() *BlockNode {
	return b.blockNode
}

func (b *Block) GetName() types.BlockName {
	return b.Name
}

func (b *Block) Actions() map[types.ActionName]BlockAction {
	return b.actions
}

func (b *Block) AddAction(name types.ActionName, action BlockAction) {
	b.actions[name] = action
}

func (b *Block) ActionsNames() []types.ActionName {
	var actionsNames []types.ActionName

	for actionName := range b.actions {
		actionsNames = append(actionsNames, actionName)
	}

	return actionsNames
}
