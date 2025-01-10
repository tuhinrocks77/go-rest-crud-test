package main

import (
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBConnection() (*gorm.DB, error) {
	// detect if running as test
	isTest := strings.HasSuffix(os.Args[0], "go-rest-crud-test.test")

	// TODO: Replace this quick & dirty "mocking" with proper mocking using repository
	// TODO: pick db conencdtion details from .env
	dbFileName := "local_dummy.db"
	if isTest {
		dbFileName = "test_.db"
	}

	// connect
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return db, err
	}

	// run migrations on connect
	err = db.AutoMigrate(&Task{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
