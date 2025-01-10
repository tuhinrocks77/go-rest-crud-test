package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	SetupRouter().Run("localhost:8080")
}
