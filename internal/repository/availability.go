package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type AvailabilityRepository interface {
	CreateAvailability(ctx context.Context, arg db.CreateAvailabilityParams) (db.Availability, error)
	GetAvailability(ctx context.Context, arg db.GetAvailabilityParams) (db.Availability, error)
	GetAvailabilityByCourseId(ctx context.Context, courseID int64) (db.Availability, error)
	ListAvailabilityForLecturer(ctx context.Context, lecturerID int64) ([]db.Availability, error)
	UpdateAvailability(ctx context.Context, arg db.UpdateAvailabilityParams) (db.Availability, error)
	DeleteAvailability(ctx context.Context, arg db.DeleteAvailabilityParams) (sql.Result, error)
}

var _ AvailabilityRepository = (*availabilityStore)(nil)

type availabilityStore struct {
	dbConn *sql.DB
	*db.Queries
}

func NewAvailabilityStore(dbConn *sql.DB) AvailabilityRepository {
	return &availabilityStore{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
