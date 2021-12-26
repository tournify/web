package web

import (
	"fmt"
	"github.com/tournify/web/config"
	"github.com/tournify/web/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

func connectToDatabase(c config.Config) (db *gorm.DB, err error) {
	return connectLoop(c, 0)
}

func connectLoop(c config.Config, count int) (db *gorm.DB, err error) {
	db, err = attemptConnection(c)
	if err != nil {
		if count > 300 {
			return db, fmt.Errorf("could not connect to database after 300 seconds")
		}
		time.Sleep(1 * time.Second)
		return connectLoop(c, count+1)
	}
	return db, err
}

func attemptConnection(c config.Config) (db *gorm.DB, err error) {
	if c.Database == "sqlite" {
		// In-memory sqlite if no database name is specified
		dsn := "file::memory:?cache=shared"
		if c.DatabaseName != "" {
			dsn = fmt.Sprintf("%s.db", c.DatabaseName)
		}
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else if c.Database == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DatabaseUsername, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if c.Database == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", c.DatabaseHost, c.DatabaseUsername, c.DatabasePassword, c.DatabaseName, c.DatabasePort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		return db, fmt.Errorf("no database specified: %s", c.Database)
	}
	return db, err
}

func migrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Session{},
		&models.Tournament{},
		&models.Game{},
		&models.Group{},
		&models.Permission{},
		&models.Player{},
		&models.Post{},
		&models.Role{},
		&models.Score{},
		&models.Subscriber{},
		&models.Team{},
		&models.TournamentOption{},
		&models.TournamentUser{},
		&models.PermissionRole{})

	seed(db)
	return err
}

func seed(db *gorm.DB) {
	adminRole := models.Role{
		Name:  "Admin",
		Label: "admin",
	}
	res := db.Save(&adminRole)
	if res.Error != nil {
		// We expect this to error after it has been created
		log.Println(res.Error)
	}
	userRole := models.Role{
		Name:  "User",
		Label: "user",
	}
	res = db.Save(&userRole)
	if res.Error != nil {
		// We expect this to error after it has been created
		log.Println(res.Error)
	}
}
