package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	// Load config from environment variables
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // fallback default
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	fmt.Println("Port:", port)
	fmt.Println("DSN:", dsn)
	fmt.Println("Redis Addr:", redisAddr)

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

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Go backend!"})
	})

	r.GET("/api/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user", "myusername")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
	})

	r.GET("/api/profile", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	fmt.Println("Backend running on port", port)
	r.Run(":" + port)
}
