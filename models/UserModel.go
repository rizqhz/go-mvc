package models

import (
	"github.com/rizghz/api/routes/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Token    string `json:"token" form:"token"`
	Blogs    []Blog `json:"blogs"`
}

type UserModel struct {
	db *gorm.DB
}

type IUserModel interface {
	Get() []User
	Find(key *int) *User
	Create(user *User) (bool, *User)
	Update(user *User) (bool, *User)
	Delete(key *int) bool
	Check(user *User) (*User, error)
}

func NewUserModel(db *gorm.DB) IUserModel {
	return &UserModel{
		db: db,
	}
}

func (m *UserModel) Get() []User {
	users := []User{}
	if err := m.db.Find(&users).Error; err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return users
}

func (m *UserModel) Find(key *int) *User {
	user := User{}
	if err := m.db.First(&user, key).Error; err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return &user
}

func (m *UserModel) Create(user *User) (bool, *User) {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	}
	return true, user
}

func (m *UserModel) Update(user *User) (bool, *User) {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	}
	return true, user
}

func (m *UserModel) Delete(key *int) bool {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&User{}, key).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func (m *UserModel) Check(user *User) (*User, error) {
	result := m.db.Where("email = ? AND password = ?", user.Email, user.Password).First(user)
	if err := result.Error; err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	var err error
	user.Token, err = middleware.CreateToken(int(user.ID))
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if err := m.db.Save(user).Error; err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return user, nil
}
