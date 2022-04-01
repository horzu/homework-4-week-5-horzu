package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string `json:"title,omitempty"`
	Page      int    `json:"page,omitempty"`
	Stock     int    `json:"stock,omitempty"`
	Price     string `json:"price,omitempty"`
	StockCode string `json:"stockCode,omitempty"`
	ISBN      string `json:"ISBN,omitempty"`
	AuthorID  uint   `json:"AuthorID,omitempty"`
	// Authors	Author	`json:"Authors,omitempty" gorm:"foreignkey:id;references:AuthorID"`
}

type Books struct {
	gorm.Model
	Title     string `json:"title,omitempty"`
	Page      int    `json:"page,omitempty"`
	Stock     int    `json:"stock,omitempty"`
	Price     string `json:"price,omitempty"`
	StockCode string `json:"stockCode,omitempty"`
	ISBN      string `json:"ISBN,omitempty"`
	AuthorID  uint   `json:"AuthorID,omitempty"`
	Authors	Author	`json:"Authors,omitempty" gorm:"foreignkey:id;references:AuthorID"`
}

func (b *Book) toString() string {
	return fmt.Sprintf("ID : %d, Title : %s, Page : %d, Stock : %d, Price : %s, StockCode : %s, ISBN : %s, AuthorID : %d, CreatedAt : %s",
		b.ID, b.Title, b.Page, b.Stock, b.Price, b.StockCode, b.ISBN, b.AuthorID, b.CreatedAt.Format("2006-01-02 15:04:05"))
}

func (b *Book) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Book (%s) deleting...\n", b.Title)
	return nil
}

func (b *Book) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Book (%s) deleted...\n", b.Title)
	return nil
}
