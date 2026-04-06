package main

import (
	"context"
	"log"
	"os"

	"agnos_assessment/api"
	"agnos_assessment/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("database unreachable:", err)
	}
	log.Println("connected to database")

	r := gin.Default()

	r.POST("/staff/create", api.CreateStaff(pool))
	r.POST("/staff/login", api.Login(pool))

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/patients/search", api.SearchPatient(pool))
	}

	r.Run(":8081")
}
