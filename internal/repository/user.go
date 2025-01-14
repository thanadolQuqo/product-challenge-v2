package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"

	"gorm.io/gorm"
	"product-challenge/internal/models"
	"product-challenge/pkg/config"
)

type UserRepository interface {
	UserRegister(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error)
	UserLogin(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error)
}

type userRepository struct {
	db  *gorm.DB
	cfg config.Config
}

func NewUserRepository(db *gorm.DB, cfg *config.Config) UserRepository {
	return &userRepository{
		db:  db,
		cfg: *cfg,
	}
}

func (r *userRepository) UserRegister(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error) {
	var (
		count int64
	)
	// 1. check if username already exist or not
	qerr := r.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error
	if qerr != nil {
		if !errors.Is(qerr, gorm.ErrRecordNotFound) {
			return nil, qerr
		}

	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	expireTime := time.Now().Add(time.Hour * 1).Unix()
	// 2. generate jwt token & set expire time
	token, err := r.GenerateJWT(*req, expireTime)
	if err != nil {
		fmt.Println("error generating token : ", err)
		return nil, err
	}
	// 3. insert to DB
	user := models.User{
		Username:  req.Username,
		Password:  req.Password,
		Token:     token,
		ExpiresAt: expireTime,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 4. return token to user
	UserCred := models.UserAuthResponse{
		Username: req.Username,
		Token:    token,
	}

	return &UserCred, nil
}

func (r *userRepository) UserLogin(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error) {
	var (
		userData models.User
		token    string
	)

	// 1. check if cred match
	qerr := r.db.Model(&models.User{}).Where("username = ?", req.Username).First(&userData).Error
	if qerr != nil {
		if errors.Is(qerr, gorm.ErrRecordNotFound) {
			return nil, errors.New("username or password is incorrect")
		}
		return nil, qerr
	}

	// Check password. 1st arg will be encrypt password. 2nd will be the input from request
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password)); err != nil {
		fmt.Println("error : ", err)
		return nil, errors.New("password is incorrect")
	}

	// 2. check if time.now() is > token expires field. if it is , generate new token, update token and expire date
	current := time.Now().Unix()

	if userData.ExpiresAt > current { // if current token not expire, use it
		token = userData.Token
	} else { // generate new token
		newToken, err := r.GenerateJWT(*req, current)
		if err != nil {
			return nil, err
		}
		// update token
		token = newToken
		userData.Token = newToken
		userData.ExpiresAt = current
		if err := r.db.Where("username = ?", req.Username).
			Updates(&userData).Error; err != nil {
			return nil, err
		}

	}
	// 3. return token
	userCred := models.UserAuthResponse{
		Username: userData.Username,
		Token:    token,
	}

	return &userCred, nil
}

func (r *userRepository) GenerateJWT(user models.UserAuthRequest, expTime int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      expTime,
	})

	return token.SignedString([]byte(r.cfg.JwtSecret))
}
