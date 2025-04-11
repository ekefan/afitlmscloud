package user

import (
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"email"`
	SchId    string `json:"sch_id"`
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
	fmt.Println("creating a user")
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
	}
	fmt.Println(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "creating a new user",
	})
}

func (us *UserService) UpdateUser(ctx *gin.Context) {
	fmt.Println("updating a user")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "updating a new user",
	})
}

func (us *UserService) DeleteUser(ctx *gin.Context) {
	fmt.Println("deleting a user")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "deleting a new user",
	})
}

func (us *UserService) GetUser(ctx *gin.Context) {
	fmt.Println("getting a user")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "geting a new user",
	})
}
