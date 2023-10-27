package main

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"strconv"
)

type config struct {
	Address     string `env:"SERVER" envDefault:"127.0.0.1"`
	Port        int    `env:"PORT" envDefault:"8080"`
	SqlUser     string `env:"POSTGRES_USER" envDefault:"defaultUser"`
	SqlPassword string `env:"POSTGRES_PASSWORD" envDefault:"somePassword"`
	SqlDB       string `env:"POSTGRES_DB" envDefault:"postgres"`
	SslMode     string `env:"SSLMODE" envDefault:"disabled"`
	DbName      string `env:"DBNAME" envDefault:"postgres"`
	DbHost      string `env:"PGHOST" envDefault:"0.0.0.0"`
}

var (
	db  *sql.DB
	cfg config
	//go:embed db/migrations/*.sql
	embedMigrations embed.FS
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	//Reading configs
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}
	//Connecting to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DbHost, 5432, cfg.SqlUser, cfg.SqlPassword, cfg.DbName)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	defer db.Close()
	//upping migrations
	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("%+v\n", err)
	}
	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Printf("%+v\n", err)
	}
	log.Println("Server almost created!")
	//Launching server
	r := gin.Default()
	r.GET("/list", getListOfTasks)
	r.DELETE("/:id", deleteTaskById)
	r.PUT("/:id", updateTaskById)
	r.POST("", createTask)
	fmt.Println("Server was created successfully!")
	if err := r.Run(":" + strconv.Itoa(cfg.Port)); err != nil {
		log.Printf("%+v\n", err)
	}
}
