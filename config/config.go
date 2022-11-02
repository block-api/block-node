package config

import (
	"github.com/block-api/block-node/common"
	"github.com/block-api/block-node/common/types"
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
	Debug             bool        `yaml:"debug"`
	Transporter       Transporter `yaml:"transporter"`
	Database          Database    `yaml:"database"`
}

// Transporter configuration options
type Transporter struct {
	Timeout           uint                 `yaml:"timeout"`
	HeartbeatInterval uint                 `yaml:"heartbeatInterval"`
	DeliveryMethod    types.DeliveryMethod `yaml:"deliveryMethod"`
	Redis             *RedisTransporter    `yaml:"redis"`
}

// RedisTransporter config options
type RedisTransporter struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       uint   `yaml:"db"`
}

// Database config options
type Database struct {
	LevelDB map[string]DatabaseLevelDB `yaml:"leveldb"`
	SQLite  map[string]DatabaseSQLite  `yaml:"sqlite"`
	CouchDB DatabaseCouchDB            `yaml:"couch_db"`
}

// DatabaseLevelDB leveldb config options
type DatabaseLevelDB struct {
	DbPath string `yaml:"path"`
}

// DatabaseSQLite sqlite config option
type DatabaseSQLite struct {
	DbPath             string `yaml:"path"`
	Options            string `yaml:"options"`
	MaxOpenConnections int    `yaml:"maxOpenConnections"`
}

// DatabaseCouchDB couchDB config options
type DatabaseCouchDB struct {
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
