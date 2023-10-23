package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	UserID  uint   `json:"user_id" form:"user_id"`
}

type BlogModel struct {
	db *gorm.DB
}

type IBlogModel interface {
	Get() []Blog
	Find(key *int) *Blog
	Create(blog *Blog) (bool, *Blog)
	Update(blog *Blog) (bool, *Blog)
	Delete(key *int) bool
}

func NewBlogModel(db *gorm.DB) IBlogModel {
	return &BlogModel{
		db: db,
	}
}

func (m *BlogModel) Get() []Blog {
	blogs := []Blog{}
	if err := m.db.Find(&blogs).Error; err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return blogs
}

func (m *BlogModel) Find(key *int) *Blog {
	blog := Blog{}
	if err := m.db.First(&blog, *key).Error; err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return &blog
}

func (m *BlogModel) Create(blog *Blog) (bool, *Blog) {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&blog).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	}
	return true, blog
}

func (m *BlogModel) Update(blog *Blog) (bool, *Blog) {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&blog).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	}
	return true, blog
}

func (m *BlogModel) Delete(key *int) bool {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Blog{}, *key).Error; err != nil {
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
