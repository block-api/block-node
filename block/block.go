package block

type IBlock interface {

	// LoadFromFile(filePath string) error
	// GetConfig() *Config
	GetName() string
}

type Block struct {
	IBlock
	Name string
}

func (b *Block) GetName() string {
	return b.Name
}
