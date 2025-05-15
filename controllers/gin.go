package controllers

import (
	"net/http"
	"test/go-crud-api/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBook associates the book with a user and stores it in the DB.
func CreateBook(c *gin.Context) {
	var book database.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate getting authenticated user's UUID
	// TODO: Replace with actual user from context/session
	userIDStr := c.GetHeader("X-User-ID") // Example: `X-User-ID` header
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	book.UserID = userID

	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "error creating book",
			"detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": book})
}

// ReadBook fetches a single book by UUID
func ReadBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var book database.Book
	if err := database.DB.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// ReadBooks fetches all books (optionally filter by user)
func ReadBooks(c *gin.Context) {
	var books []database.Book

	// Optional: filter by user
	userIDStr := c.Query("user_id")
	query := database.DB
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	if err := query.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch books"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// UpdateBook updates book details
func UpdateBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var input database.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var book database.Book
	if err := database.DB.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if err := database.DB.Model(&book).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DeleteBook deletes a book by UUID
func DeleteBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var book database.Book
	if err := database.DB.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// CreateUser handles creation of a new user
func CreateUser(c *gin.Context) {
	var user database.User

	// Bind JSON input
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: validate email, name here

	// Check if user with email already exists
	var existing database.User
	if err := database.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
		return
	}

	// Create user
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// GetUser returns user details by UUID, including their books
func GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var user database.User
	// Preload Books to get user's books in the same query
	if err := database.DB.Preload("Books").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser deletes a user and associated books
func DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var user database.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Optional: delete associated books
	if err := database.DB.Where("user_id = ?", id).Delete(&database.Book{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user's books"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
