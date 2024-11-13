package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Name   string
	Email  string
	Orders []Order
}

type Order struct {
	gorm.Model
	UserId      int64
	OrderTime   time.Time
	PaymentMode string // Card or Cash
	Price       int
	User        User
}

var DB *gorm.DB

func connectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	database, err := gorm.Open(mysql.Open("codeheim:tmp_pwd@tcp(127.0.0.1:3306)/gorm_scopes?charset=utf8&parseTime=true"), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Failed to connect to databse!")
	}

	DB = database
}

func dbMigrate() {
	DB.AutoMigrate(&User{}, &Order{})
}

func CardOrders(db *gorm.DB) *gorm.DB {
	return db.Where("payment_mode = ?", "card")
}

func PriceGreaterThan30(db *gorm.DB) *gorm.DB {
	return db.Where("price > ?", 30)
}

func UsersFromDomain(domain string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email like ?", "%"+domain)
	}
}

func main() {
	connectDatabase()
	dbMigrate()

	var orders []Order
	DB.Scopes(CardOrders, PriceGreaterThan30).Find(&orders)

	fmt.Println("orders:")
	for _, order := range orders {
		fmt.Printf("Price: %d, Payment Type: %s\n", order.Price, order.PaymentMode)
	}

	var users []User
	DB.Scopes(UsersFromDomain("example.com")).Preload("Orders", CardOrders).Find(&users)

	fmt.Printf("Users: \n")
	for _, user := range users {
		fmt.Printf("User email: %s\n", user.Email)
	}

	fmt.Printf("Orders from a user (%s): \n", users[0].Email)
	for _, order := range users[0].Orders {
		fmt.Printf("Price: %d, Payment Type: %s\n", order.Price, order.PaymentMode)
	}
}
