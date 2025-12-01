package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	// defineUserRoutes(r)
	defineEventRoutes(r)
	defineSchedulerRoutes(r)
	defineContactRoutes(r)
	defineConsolidatedRoutes(r)
	return r
}

// func defineUserRoutes(r *gin.Engine) {
// 	r.POST("/api/users", func(c *gin.Context) {
// 		fmt.Print("Not implemented")
// 		c.JSON(501, gin.H{"message": "not implemented"})
// 	})
// 	r.GET("/api/users/:id", func(c *gin.Context) {
// 		fmt.Print("Not implemented")
// 		c.JSON(501, gin.H{"message": "not implemented"})
// 	})
// 	r.GET("/api/users/", func(c *gin.Context) {
// 		fmt.Print("Not implemented")
// 		c.JSON(501, gin.H{"message": "not implemented"})
// 	})
// 	r.PUT("/api/users/:id", func(c *gin.Context) {
// 		fmt.Print("Not implemented")
// 		c.JSON(501, gin.H{"message": "not implemented"})
// 	})
// }

func defineEventRoutes(r *gin.Engine) {
	r.POST("/api/events", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/events/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/events/", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.PUT("/api/events/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
}

func defineSchedulerRoutes(r *gin.Engine) {
	r.POST("/api/schedulers", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/schedulers/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/schedulers/", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.PUT("/api/schedulers/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
}

func defineContactRoutes(r *gin.Engine) {
	r.POST("/api/contacts", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/contacts/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/contacts/", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.PUT("/api/contacts/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
}

func defineConsolidatedRoutes(r *gin.Engine) {
	r.POST("/api/consolidateds", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/consolidateds/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.GET("/api/consolidateds/", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
	r.PUT("/api/consolidateds/:id", func(c *gin.Context) {
		fmt.Print("Not implemented")
		c.JSON(501, gin.H{"message": "not implemented"})
	})
}
