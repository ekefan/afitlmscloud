package user

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserReq struct {
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"email"`
	SchId    string `json:"sch_id"`
}

type UpdatedUserReq struct {
	Fullname string `json:"fullname,omitempty"`
	Password string `json:"passowrd,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UserDTO struct {
	ID              int64     `json:"id,omitempty"`
	Fullname        string    `json:"fullname,omitempty"`
	Password        string    `json:"passowrd,omitempty"`
	Email           string    `json:"email,omitempty"`
	PasswordChanged bool      `json:"password_changed"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}
type UserService struct {
	repo repository.UserRespository
}

func NewUserService(userRepo repository.UserRespository) *UserService {
	us := &UserService{
		repo: userRepo,
	}
	return us
}

func (us *UserService) CreateUser(ctx *gin.Context) {
	var req CreateUserReq

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
	}
	// TODO: Hash passwords, CrossOrigin stuff... Validation,
	user, err := us.repo.CreateUser(ctx, db.CreateUserParams{
		Email:          req.Email,
		FullName:       req.Fullname,
		HashedPassword: req.Password,
		SchID:          req.SchId,
	})
	if err != nil {
		slog.Error("Unhandled error here", "detals", err)
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				column := ""
				if pqErr.Constraint == "users_email_key" {
					column = "email"
				}
				if pqErr.Constraint == "users_sch_id_key" {
					column = "sch_id"
				}
				msg := fmt.Sprintf("%v already exists", column)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": msg,
				})
			default:
				slog.Error("Unhandled pq error", "details", pqErr)
			}
			return
		}
		return
	}
	fmt.Println(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "created a new user",
	})
}

func (us *UserService) UpdateUser(ctx *gin.Context) {
	fmt.Println("updating a user")

	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}
	var req UpdatedUserReq
	if err := ctx.BindJSON(&req); err != nil || (req.Email == "" &&
		req.Fullname == "" &&
		req.Password == "") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// TODO: Hash passwords if provided,
	user, err := us.repo.UpdateUser(ctx, db.UpdateUserParams{
		ID:              int64(userId),
		Email:           req.Email,
		FullName:        req.Fullname,
		PasswordChanged: true,
		HashedPassword:  req.Password,
	})
	if err != nil {
		slog.Error("Unhandled error here", "detals", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with such user id",
			})
			return
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				column := ""
				if pqErr.Constraint == "users_email_key" {
					column = "email"
				}
				if pqErr.Constraint == "users_sch_id_key" {
					column = "sch_id"
				}
				msg := fmt.Sprintf("%v already exists", column)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": msg,
				})
			default:
				slog.Error("Unhandled pq error", "details", pqErr)
			}
			return
		}
		return
	}
	fmt.Println(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "creating a new user",
	})
}

func (us *UserService) DeleteUser(ctx *gin.Context) {
	fmt.Println("deleting a user")

	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	res, err := us.repo.DeleteUser(ctx, int64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to delete a user with id: %v", userId),
		})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("user with id %d doesn't exist", userId),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("deleted user with id: %v", userId),
	})
}

func (us *UserService) GetUser(ctx *gin.Context) {
	fmt.Println("getting a user")

	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	user, err := us.repo.GetUserByID(ctx, int64(userId))
	if err != nil {
		slog.Error("Unhandled repo error", "dtails", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("user with id: %v not found", userId),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get a user with id: %v", userId),
		})
		return
	}
	resp := UserDTO{
		ID:              user.ID,
		Fullname:        user.FullName,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt.Time,
		PasswordChanged: user.PasswordChanged,
	}

	ctx.JSON(http.StatusOK, resp)
}
