package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pgClient *PGClient

type PGClient struct {
	DB *gorm.DB
}

type PGConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

func (c *PGConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		c.Host,
		c.User,
		c.Password,
		c.DBName,
		c.Port,
	)
}

func InitPGClient(config *PGConfig) error {
	dsn := config.GetDSN()
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	pgClient = &PGClient{}
	pgClient.DB = conn
	return nil
}

func GetPGClient(config *PGConfig) (*PGClient, error) {
	if pgClient != nil {
		return pgClient, nil
	}
	if err := InitPGClient(config); err != nil {
		return nil, err
	}
	return pgClient, nil
}
