package sql

import (
	"github.com/thebluefowl/suckerfish/db"
)

type Migration func(*db.PGClient) error

var migrations = []Migration{
	func(client *db.PGClient) error {
		if !client.DB.Migrator().HasTable(&User{}) {
			return client.DB.AutoMigrate(&User{})
		}
		return nil
	},
}

func Run(client *db.PGClient) error {
	for _, m := range migrations {
		if err := m(client); err != nil {
			return err
		}
	}
	return nil
}
