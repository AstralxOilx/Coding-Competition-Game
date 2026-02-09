package service

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/dto"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
)

type UserService interface {
	AllUsers() ([]model.Users, error)
	Profile(userID string) (*dto.ProfileResponse, error)
	// ✅ 1. เพิ่มเข้า Interface เพื่อให้ Handler เรียกใช้ได้
	UpdateUserInfo(userID string, displayName string, avatarURL string) (*dto.ProfileResponse, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(r repository.UserRepo) UserService {
	return &userService{repo: r}
}

// ... (GetAllUsers และ GetProfile เหมือนเดิม) ...

func (s *userService) AllUsers() ([]model.Users, error) {
	return s.repo.FindAllUser()
}

func (s *userService) Profile(userID string) (*dto.ProfileResponse, error) {
	user, err := s.repo.FindById(userID)
	if err != nil {
		return nil, err
	}

	var rankResponses []dto.UserRankResponse
	for _, r := range user.Ranks {
		winRate := 0.0
		if r.TotalGames > 0 {
			winRate = (float64(r.Win) / float64(r.TotalGames)) * 100
		}
		rankResponses = append(rankResponses, dto.UserRankResponse{
			ModeName:   r.ModeName,
			Rank:       r.Rank,
			RankTier:   r.RankTier,
			RankPoint:  r.RankPoint,
			WinRate:    winRate,
			TotalGames: r.TotalGames,
		})
	}

	return &dto.ProfileResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		PlayerLevel: user.PlayerLevel,
		PlayerExp:   user.PlayerExp,
		LastLogin:   user.LastLogin,
		Ranks:       rankResponses,
	}, nil
}

func (s *userService) UpdateUserInfo(userID string, displayName string, avatarURL string) (*dto.ProfileResponse, error) {
	// 1. สร้างก้อนข้อมูลจาก MODEL (ตัวที่ใช้คุยกับ DB จริงๆ)
	updateData := &model.Users{
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}

	// 2. ส่ง MODEL ไปให้ Repo (ไม่ใช่ส่ง ProfileResponse)
	_, err := s.repo.UpdateUserInfo(userID, updateData)
	if err != nil {
		return nil, err
	}

	// 3. พออัปเดตเสร็จค่อยเรียก GetProfile เพื่อเอาข้อมูลมาแปลงเป็น DTO ส่งกลับ
	return s.Profile(userID)
}
