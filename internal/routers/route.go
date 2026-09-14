package routers

import (
	"gotta/internal/configs"
	"gotta/internal/db"
	"gotta/internal/handlers"
	"gotta/internal/models"
	"gotta/internal/repository"
	"gotta/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type application struct {
	router *gin.Engine
	db     *sqlx.DB
	cfg    models.Config
}

func NewApp() (*application, error) {
	cfg := configs.Load()

	conn, err := db.Connect(cfg.DbPath)
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	r.Static("/static", "./web/static")

	todoRepo := repository.NewTodoRepository(conn)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandlerHTML(todoService)

	r.GET("/", todoHandler.Index)
	r.POST("/add", todoHandler.Add)
	r.POST("/toggle/:id", todoHandler.Toggle)
	r.DELETE("/delete/:id", todoHandler.Delete)

	return &application{
		router: r,
		db:     conn,
		cfg:    cfg,
	}, nil
}

func (a *application) Run() error {
	log.Printf("Server listening on port: %v\n", a.cfg.ServerPort)
	return a.router.Run(a.cfg.ServerPort)
}

func (a *application) Close() error {
	return a.db.Close()
}
