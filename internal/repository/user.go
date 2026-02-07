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
	err := r.db.Where("id = ?", id).First(&user).Error

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
