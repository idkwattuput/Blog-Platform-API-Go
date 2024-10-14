package api

import (
	"database/sql"
	"github.com/idkwattuput/blogging-platform-api-go/services/post"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	postStore := post.NewStore(s.db)
	postHandler := post.NewHandler(postStore)
	postHandler.RegisterRoutes(router)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	server := http.Server{
		Addr:    s.addr,
		Handler: v1,
	}

	log.Printf("Server has started  %s", s.addr)

	return server.ListenAndServe()
}
