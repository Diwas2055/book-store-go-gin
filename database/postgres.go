package database

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB instance
var DB *gorm.DB
var err error

// User model with UUID primary key
type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email" gorm:"unique"`
	Books []Book    `gorm:"foreignKey:UserID"`
}

// Book model with UUID primary key and UserID as foreign key
type Book struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	UserID uuid.UUID `gorm:"type:uuid" json:"userId"` // FK to User
}

// Hook to set UUIDs before create
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

// Connect to the database and auto-migrate
func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "go_dbtest"
	dbUser := "postgres"
	password := "postgres"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&User{}, &Book{})
	if err != nil {
		log.Fatal("AutoMigrate error:", err)
	}

	fmt.Println("Database connection successful and models migrated with UUID...")
}
