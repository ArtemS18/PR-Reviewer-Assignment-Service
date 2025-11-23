package postgres

import (
	"log"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/service/pull_request"
	"reviewer-api/internal/app/service/team"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgers(dsn string, autoMigrate bool) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if autoMigrate {
		err = db.AutoMigrate(
			&ds.PullRequest{},
			&ds.User{},
			&ds.Reviewer{},
			&ds.Team{},
		)
		if err != nil {
			panic("cant migrate db")
		}
		log.Println("success migrate db")
	}
	return &Postgres{db}, nil
}

func (p *Postgres) WithPRTransaction(fn pull_request.TxFunc) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Postgres{db: tx}
		return fn(txRepo)
	})
}
func (p *Postgres) WithTeamTransaction(fn team.TxFunc) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Postgres{db: tx}
		return fn(txRepo)
	})
}
