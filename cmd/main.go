package main

import (
	"context"
	"net/http"
	"os"

	"github.com/fykyby/chat-app-backend/database"
	"github.com/fykyby/chat-app-backend/handler"
	"github.com/fykyby/chat-app-backend/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Database connection
	ctx := context.Background()
	dbConnection, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close(ctx)
	db := database.New(dbConnection)

	// JWT
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CLIENT_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Content-Length", "Accept-Encoding", "Set-Cookie", "Origin"},
		AllowCredentials: true,
	}))

	apiHandler := handler.ApiHandler{
		DB:        db,
		TokenAuth: tokenAuth,
	}
	r.Route("/api", apiHandler.Handler)

	wsHandler := ws.WebSocketHandler{
		DB:        db,
		TokenAuth: tokenAuth,
	}
	r.Route("/ws", wsHandler.Handler)

	http.ListenAndServe(":3001", r)
}
