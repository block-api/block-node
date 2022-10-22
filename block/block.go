package block

type IBlock interface {
	// LoadFromFile(filePath string) error
	// GetConfig() *Config
	GetName() string
	Actions() map[string]BlockAction
}

type BlockAction func(payload interface{}) error

type Block struct {
	IBlock
	Name    string
	Actions map[string]BlockAction
}

func (b *Block) GetName() string {
	return b.Name
}
