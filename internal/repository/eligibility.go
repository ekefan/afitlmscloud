package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type EligibilityRepository interface {
	CreateEligibility(ctx context.Context, arg db.CreateEligibilityParams) (db.Eligibility, error)
	DeleteEligibility(ctx context.Context, arg db.DeleteEligibilityParams) (sql.Result, error)
	GetEligibility(ctx context.Context, arg db.GetEligibilityParams) (db.Eligibility, error)
	ListEligibilityForStudent(ctx context.Context, studentID int64) ([]db.Eligibility, error)
	UpdateEligibility(ctx context.Context, arg db.UpdateEligibilityParams) (db.Eligibility, error)
	SetMinEligibility(ctx context.Context, arg db.SetMinEligibilityParams) (db.Eligibility, error)
}

var _ EligibilityRepository = (*eligibilityStore)(nil)

type eligibilityStore struct {
	dbConn *sql.DB
	*db.Queries
}

func NewEligibilityStore(dbConn *sql.DB) EligibilityRepository {
	return &eligibilityStore{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
