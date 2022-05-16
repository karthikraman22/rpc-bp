package database

import (
	"fmt"
	"testing"

	"github.com/karthikraman22/rpc-bp/config"
)

func TestInitDb(t *testing.T) {
	cfg, err := config.NewConfig("../test-conf.yaml", "DB_TEST_")
	if err != nil {
		fmt.Println(err)
	}
	db, err := InitDatabase(cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("db.ConnPool: %v\n", db.ConnPool)
}
