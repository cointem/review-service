package data

import (
	"errors"

	"review-service/internal/conf"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewReviewRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	query *query.Query
	log   *log.Helper
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	query.SetDefault(db)
	return &Data{
		query: query.Q,
		log:   log.NewHelper(logger),
	}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to open database")
	}
	return db, nil
}
