package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/kiln-mid/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DelegationsRepository interface {
	CreateMany(ctx context.Context, Delegations *[]models.Delegations) (int64, error)
	FindMostRecent(ctx context.Context) (*models.Delegations, error)
	FindAndOrderByTimestamp(ctx context.Context, limit int, offset int) (*[]models.Delegations, error)
	FindFromYear(ctx context.Context, year int, limit int, offset int) (*[]models.Delegations, error)
	FindAvailableYear(ctx context.Context) (*[]int, error)
}

// NewDelegationsAdapter returns an implementation of the DelegationsRepository using GORM for database interactions.
func NewDelegationsAdapter(db *gorm.DB) DelegationsRepository {
	return &DelegationsAdapter{DB: db}
}

// DelegationsAdapter provides a GORM-based implementation of DelegationsRepository.
type DelegationsAdapter struct {
	DB *gorm.DB
}

// CreateMany inserts multiple Delegations records into the database, no error is returned if their is a conflict based on UNIQUE key
// return the number of inserted rows.
func (r *DelegationsAdapter) CreateMany(ctx context.Context, d *[]models.Delegations) (int64, error) {
	res := r.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id_tezos", Table: "delegations"}},
		DoNothing: true,
	}).Create(&d)

	if res.Error != nil {
		return 0, fmt.Errorf("gorm error: %s", res.Error)
	}

	return res.RowsAffected, nil
}

// FindAvailableYear return a slice of available year which delegations can be searched.
func (r *DelegationsAdapter) FindAvailableYear(ctx context.Context) (*[]int, error) {
	var years []int

	res := r.DB.Model(models.Delegations{}).
		Select("DISTINCT YEAR(timestamp) AS year").
		Order("year").
		Pluck("year", &years)

	if res.Error != nil {
		return nil, fmt.Errorf("gorm error: %s", res.Error)
	}

	return &years, nil
}

// FindFromYear fetch and return with a limit and an offset all delegations who can be found for a given year.
func (r *DelegationsAdapter) FindFromYear(ctx context.Context, year int, limit int, offset int) (*[]models.Delegations, error) {
	var d []models.Delegations

	res := r.DB.Limit(limit).
		Offset(offset).Where("YEAR(timestamp) = ?", year).Order("timestamp desc").Find(&d)
	if res.Error != nil {
		return nil, fmt.Errorf("gorm error: %s", res.Error)
	}

	return &d, nil
}

// FindAndOrderByTimestamp fetch and return with a limit and an offset all delegations ordered by timestamp.
func (r *DelegationsAdapter) FindAndOrderByTimestamp(ctx context.Context, limit int, offset int) (*[]models.Delegations, error) {
	var d []models.Delegations

	res := r.DB.Limit(limit).
		Offset(offset).Order("timestamp desc").Find(&d)
	if res.Error != nil {
		return nil, fmt.Errorf("gorm error: %s", res.Error)
	}

	return &d, nil
}

// FindMostRecent fetch and return the most recent delegations.
func (r *DelegationsAdapter) FindMostRecent(ctx context.Context) (*models.Delegations, error) {
	var d models.Delegations

	res := r.DB.Order("timestamp desc").First(&d)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &d, nil
	}

	if res.Error != nil {
		return &d, fmt.Errorf("gorm error: %s", res.Error)
	}

	return &d, nil
}
