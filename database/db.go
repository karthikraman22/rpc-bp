package database

import (
	"fmt"
	"time"

	"github.com/karthikraman22/rpc-bp/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initializes the database with standard configuration
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", cfg.String("db.driver"),
		cfg.String("db.user"), cfg.String("db.password"), cfg.String("db.host"), cfg.Int("db.port"),
		cfg.String("db.database"), cfg.String("db.sslmode"))

	gormCfg := gorm.Config{Logger: NewGormLogger()}
	if db, err := gorm.Open(postgres.New(postgres.Config{DSN: connectionString}), &gormCfg); err != nil {
		return nil, err
	} else {
		sqlDB, err := db.DB()
		if err == nil {
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDB.SetMaxIdleConns(2)
			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDB.SetMaxOpenConns(cfg.Int("db.poolSize"))
			sqlDB.SetConnMaxLifetime(time.Hour)
		}
		return db, err
	}
}
