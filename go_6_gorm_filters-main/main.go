package main

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("before!")
	if user.Name == "" {
		return errors.New("name can not be blank")
	}
	return nil
}

func (user *User) AfterCreate(tx *gorm.DB) (err error) {
	// sendWelcomeEmail(user)
	fmt.Println("after!")
	fmt.Printf("User ID: %d, email: %s \n", user.ID, user.Email)
	return nil
}

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the model to create the "users" table
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to auto-migrate")
	}

	// Create a new user
	user := User{Name: "John", Email: "john@codeheim.com"}
	result := db.Create(&user)
	if result.Error != nil {
		panic("failed to create user")
	}
}
