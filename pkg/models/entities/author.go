package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name  string `json:"Name"`
	Books []Book `json:"Books,omitempty" gorm:"foreignKey:AuthorID;references:id"`
}

func (a *Author) toString() string {
	return fmt.Sprintf("ID : %d, Name : %s, CreatedAt : %s",
		a.ID, a.Name, a.CreatedAt.Format("2006-01-02 15:04:05"))
}

func (a *Author) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Author (%s) deleting...\n", a.Name)
	return nil
}

func (a *Author) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Author (%s) deleted...\n", a.Name)
	return nil
}
