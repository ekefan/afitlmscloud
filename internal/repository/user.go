package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type UserRespository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id int64) (sql.Result, error)
	GetUserByID(ctx context.Context, id int64) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	UpdateUserEmail(ctx context.Context, arg db.UpdateUserEmailParams) (db.User, error)
	UpdateUserPassword(ctx context.Context, arg db.UpdateUserPasswordParams) (db.User, error)
	EnrollUser(ctx context.Context, arg db.EnrollUserParams) (db.User, error)
}

var _ UserRespository = (*UserStore)(nil)

type UserStore struct {
	dbConn *sql.DB
	*db.Queries
}

func NewUserStore(dbConn *sql.DB) UserRespository {
	return &UserStore{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
