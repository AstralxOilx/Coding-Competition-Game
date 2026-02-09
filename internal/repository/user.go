package repository

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(user *model.Users) error

	FindById(id string) (*model.Users, error)
	FindByEmail(email string) (*model.Users, error)
	FindByUserName(userName string) (*model.Users, error)
	FindAllUser() ([]model.Users, error)
	FindFriendIDs(userID string) ([]string, error)
	FindUserInfo(userID string) (*model.Users, error)

	UpdateUserInfo(userID string, updatedData *model.Users) (*model.Users, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.Users) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindById(id string) (*model.Users, error) {
	var user model.Users
	// err := r.db.Where("id = ?", id).First(&user).Error
	err := r.db.Preload("Ranks").Where("id = ?", id).First(&user).Error

	return &user, err
}

func (r *userRepo) FindByUserName(userName string) (*model.Users, error) {
	var user model.Users
	err := r.db.Where("user_name = ? ", userName).Find(&user).Error
	return &user, err
}

func (r *userRepo) FindByEmail(email string) (*model.Users, error) {
	var user model.Users
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) FindAllUser() ([]model.Users, error) {

	var users []model.Users
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) FindFriendIDs(userID string) ([]string, error) {
	var friendIDs []string

	// Query หา ID ของคนที่เป็นเพื่อนกับเรา (สถานะต้องเป็น 1 หรือ Accepted)
	err := r.db.Table("friendships").
		Where("(user_id = ? OR friend_id = ?) AND status = 1", userID, userID).
		Select("CASE WHEN user_id = ? THEN friend_id ELSE user_id END", userID).
		Find(&friendIDs).Error

	if err != nil {
		return nil, err
	}

	return friendIDs, nil
}

func (r *userRepo) FindUserInfo(userID string) (*model.Users, error) {
	var user model.Users
	// ดึงเฉพาะฟิลด์ที่จำเป็นเพื่อความเร็ว
	err := r.db.Select("id", "display_name", "avatar_url").Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *userRepo) UpdateUserInfo(userID string, updateData *model.Users) (*model.Users, error) {
	// ใช้ Model จริงในการอัปเดต
	err := r.db.Model(&model.Users{}).
		Where("id = ?", userID).
		Updates(model.Users{
			DisplayName: updateData.DisplayName,
			AvatarURL:   updateData.AvatarURL,
		}).Error

	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลใหม่กลับมา
	var user model.Users
	err = r.db.Preload("Ranks").Where("id = ?", userID).First(&user).Error
	return &user, err
}
