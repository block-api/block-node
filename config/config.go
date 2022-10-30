package config

import (
	"github.com/block-api/block-node/common"
	"github.com/block-api/block-node/log"
)

// ApplicationConfig interface
type ApplicationConfig interface {
	LoadFromFile(filePath string) error
	GetConfig() *Config
}

// Config file definition
type Config struct {
	ApplicationConfig `yaml:"-"`
	Debug             bool           `yaml:"debug"`
	Transporter       Transporter    `yaml:"transporter"`
	Database          ConfigDatabase `yaml:"database"`
}

// Transporter configuration options
type Transporter struct {
	Timeout           uint              `yaml:"timeout"`
	HeartbeatInterval uint              `yaml:"heartbeatInterval"`
	Redis             *RedisTransporter `yaml:"redis"`
}

// RedisTransporter config options
type RedisTransporter struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       uint   `yaml:"db"`
}

type ConfigDatabase struct {
	LevelDB map[string]ConfigLevelDB `yaml:"leveldb"`
	CouchDB ConfigDatabaseCouchDB    `yaml:"couch_db"`
}

type ConfigLevelDB struct {
	DbPath string `yaml:"path"`
}

type ConfigDatabaseCouchDB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
}

// NewConfig creates new empty Config struct
func NewConfig() Config {
	return Config{}
}

// GetConfig returns Config struct options
func (c *Config) GetConfig() *Config {
	return c
}

// LoadFromFile load config yaml file from path
func (c *Config) LoadFromFile(filePath string) error {
	log.New("load config from file: "+filePath, log.TypeDebug)

	file, err := common.OpenFile(filePath, common.YML)
	if err != nil {
		return err
	}

	defer file.Close()

	errUnmarshal := file.Parse(&c)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	return nil
}
