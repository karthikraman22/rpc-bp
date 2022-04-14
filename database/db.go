package database

import (
	"fmt"

	"achuala.in/rpc-bp/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initializes the database with standard configuration
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", cfg.String("db.driver"),
		cfg.String("db.user"), cfg.String("db.password"), cfg.String("db.host"), cfg.Int("db.port"),
		cfg.String("db.database"), cfg.String("db.sslmode"))

	gormCfg := gorm.Config{Logger: NewGormLogger()}
	return gorm.Open(postgres.New(postgres.Config{DSN: connectionString}), &gormCfg)
}
