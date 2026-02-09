package service

import (
	"context"
	"errors"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/config"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/dto"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository/cache/redis"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Signup(req dto.SignupRequest) (string, error)
	Signin(req dto.SigninRequest) (*dto.SigninResponse, error)
	RefreshToken(refreshToken string) (map[string]string, error)
}

type authService struct {
	userRepo repository.UserRepo
}

func NewAuthService(repo repository.UserRepo) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Signup(req dto.SignupRequest) (string, error) {
	hashedPassword, _ := util.HashPassword(req.Password)
	user := model.Users{
		ID:          util.GenerateID(14),
		DisplayName: req.DisplayName,
		UserName:    req.UserName,
		Email:       req.Email,
		Password:    hashedPassword,
	}
	if err := s.userRepo.CreateUser(&user); err != nil {
		return "", errors.New("exists")
	}
	return user.ID, nil
}

func (s *authService) Signin(req dto.SigninRequest) (*dto.SigninResponse, error) {
	// สร้าง context พื้นฐานสำหรับใช้ในฟังก์ชันนี้
	ctx := context.Background()

	user, err := s.userRepo.FindByUserName(req.UserName)
	if err != nil || !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("unauthorized")
	}

	// ✅ แก้จาก nil เป็น ctx
	oldToken, err := redis.GetUserSession(ctx, user.ID)
	if err == nil && oldToken != "" {
		_ = database.RDB.Del(ctx, "session:"+user.ID).Err()
		return nil, errors.New("conflict")
	}

	accessToken, _ := util.GenerateAccessToken(user.ID, user.UserRole)
	refreshToken, _ := util.GenerateRefreshToken(user.ID)

	// ✅ แก้จาก nil เป็น ctx
	err = redis.SetUserSession(ctx, user.ID, refreshToken, config.AppConfig.SessionDuration)
	if err != nil {
		return nil, err
	}

	return &dto.SigninResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User: dto.UserInfo{
			DisplayName: user.DisplayName,
			Role:        user.UserRole, // อย่าลืมเช็คเรื่อง int/string ตามที่คุยกันก่อนหน้านี้ครับ
		},
	}, nil
}

func (s *authService) RefreshToken(tokenStr string) (map[string]string, error) {
	token, err := util.ValidateRefreshToken(tokenStr)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid_token")
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID, _ := claims["user_id"].(string)
	user, _ := s.userRepo.FindById(userID)

	newAccess, _ := util.GenerateAccessToken(user.ID, user.UserRole)
	newRefresh, _ := util.GenerateRefreshToken(user.ID)

	return map[string]string{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	}, nil
}
