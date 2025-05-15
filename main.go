package main

import (
	"fmt"
	"log"

	"test/go-crud-api/controllers"
	"test/go-crud-api/database"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting application ...")
	database.DatabaseConnection()

	r := gin.Default()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	r.GET("/books/:id", controllers.ReadBook)
	r.GET("/books", controllers.ReadBooks)
	r.POST("/books", controllers.CreateBook)
	r.PUT("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)
	r.POST("/users", controllers.CreateUser)
	r.GET("/users/:id", controllers.GetUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Server is Running....")
	})

	// r.SetTrustedProxies([]string{"127.0.0.1"}) //  developing locally,
	r.SetTrustedProxies(nil) // for development only
	r.Run(":5000")
}
