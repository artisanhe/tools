package mysql

import (
	"fmt"
	"os"
	"time"

	"github.com/artisanhe/gorm"
)

func connectMysql(conn string, poolSize int, maxLifetime time.Duration) (*gorm.DB, error) {
	database, err := gorm.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	database.DB().SetMaxOpenConns(poolSize)
	database.DB().SetMaxIdleConns(poolSize / 2)
	database.DB().SetConnMaxLifetime(maxLifetime)
	database.SingularTable(true)
	database.LogMode(true)

	err = database.DB().Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect mysql failed[%s]\n", err)
		return nil, err
	}

	return &database, nil
}
