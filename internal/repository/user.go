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
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
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
