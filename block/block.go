package block

import (
	"github.com/block-api/block-node/utils"
)

type IBlock interface {
	GetName() utils.BlockName
	Actions() map[utils.ActionName]BlockAction
	ActionsNames() []utils.ActionName
	AddAction(name utils.ActionName, action BlockAction)
}

type BlockAction func(payload interface{}) error

type Block struct {
	IBlock
	Name    utils.BlockName
	actions map[utils.ActionName]BlockAction
}

func NewBlock(name utils.BlockName) Block {
	return Block{
		Name:    name,
		actions: make(map[utils.ActionName]BlockAction),
	}
}

func (b *Block) GetName() utils.BlockName {
	return b.Name
}

func (b *Block) Actions() map[utils.ActionName]BlockAction {
	return b.actions
}

func (b *Block) AddAction(name utils.ActionName, action BlockAction) {
	b.actions[name] = action
}

func (b *Block) ActionsNames() []utils.ActionName {
	var actionsNames []utils.ActionName

	for actionName := range b.actions {
		actionsNames = append(actionsNames, actionName)
	}

	return actionsNames
}
