package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Cozzytree/nait/internal/auth"
	"github.com/Cozzytree/nait/internal/database"
	_ "github.com/lib/pq"
)

type my_server struct {
	db *database.Queries
}

func InitServer() *http.Server {
	auth.InitAuth()

	db_url := os.Getenv("DATABASE_URL")
	db_conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	s := my_server{
		db: database.New(db_conn),
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	return &http.Server{
		Addr:    port,
		Handler: s.registerRoutes(),
	}
}
