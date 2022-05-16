package database

import (
	"fmt"
	"testing"

	"github.com/karthikraman22/rpc-bp/config"
)

func TestInitDb(t *testing.T) {
	cfg := config.NewConfig("../test-conf.yaml", "DB_TEST_")
	db, err := InitDatabase(cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("db.ConnPool: %v\n", db.ConnPool)
}
