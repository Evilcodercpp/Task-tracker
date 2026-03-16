package userservice

import "gorm.io/gorm"

type UserRepository interface {
	CreateUser(user *User) error
	GetAllUsers() ([]User, error)
	GetUserByID(ID string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(ID string) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepo) GetUserByID(ID string) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepo) DeleteUser(ID string) error {
	return r.db.Delete(&User{}, "id = ?", ID).Error
}
