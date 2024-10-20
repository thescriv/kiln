package db

import (
	"github.com/kiln-mid/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Client represents a structure containing a database connection using GORM.
type Client struct {
	DB *gorm.DB
}

// CreateClient initializes a new Client with a database connection based on the DSN given in param.
// If the connection is successfull it automatically autoMigrate models.
func CreateClient(DSN string) (Client, error) {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		return Client{}, err
	}

	db.AutoMigrate(&models.Delegations{})

	client := Client{
		DB: db,
	}

	return client, nil
}
