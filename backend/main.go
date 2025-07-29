package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	cache *redis.Client
)

func main() {
	// Load config
	port := "8080"
	dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
	redisAddr := "localhost:6379"

	// Connect Postgres
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Connect Redis
	cache = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if err := cache.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	// API routes
	r := gin.Default()
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Go backend!"})
	})

	r.GET("/api/login", func(c *gin.Context) {
		username := "myusername"
		c.JSON(http.StatusOK, gin.H{
			"message": "User logged in",
			"user":    username,
		})
	})

	// ดึง user จาก query parameter (เช่น /api/profile?user=myusername)
	r.GET("/api/profile", func(c *gin.Context) {
		user := c.Query("user") // รับค่าผ่าน query string
		if user == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	fmt.Println("Backend running on port", port)
	r.Run(":" + port)
}
