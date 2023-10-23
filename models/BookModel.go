package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string `json:"title" form:"title"`
	Author    string `json:"author" form:"author"`
	Publisher string `json:"publisher" form:"publisher"`
}

type BookModel struct {
	db *gorm.DB
}

type IBookModel interface {
	Get() []Book
	Find(key *int) *Book
	Create(user *Book) (bool, *Book)
	Update(user *Book) (bool, *Book)
	Delete(key *int) bool
}

func NewBookModel(db *gorm.DB) IBookModel {
	return &BookModel{
		db: db,
	}
}

func (m *BookModel) Get() []Book {
	books := []Book{}
	if err := m.db.Find(&books).Error; err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return books
}

func (m *BookModel) Find(key *int) *Book {
	book := Book{}
	if err := m.db.First(&book, *key).Error; err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return &book
}

func (m *BookModel) Create(book *Book) (bool, *Book) {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(book).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	}
	return true, book
}

func (m *BookModel) Update(book *Book) (bool, *Book) {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(book).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	}
	return true, book
}

func (m *BookModel) Delete(key *int) bool {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Book{}, *key).Error; err != nil {
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
