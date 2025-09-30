package server

import (
	"blog-service/internal/db/mongo"
	pg "blog-service/internal/db/postgres"
	pb "blog-service/internal/grpc/protobuf"
	"log"
	"net/http"
	"time"

	"blog-service/internal/server/handlers"

	"github.com/rs/cors"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

type Server struct {
	mux            *http.ServeMux
	mongoClient    *mongo.Client
	postgresClient *pg.Client
	authClient     pb.AuthServiceClient
}

func NewServer(mongoClient *mongo.Client, postgresClient *pg.Client, authClient pb.AuthServiceClient) *Server {
	mux := http.NewServeMux()

	s := &Server{
		mux:            mux,
		mongoClient:    mongoClient,
		postgresClient: postgresClient,
		authClient:     authClient,
	}

	s.registerRoutes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	protectedArticleHandler := handlers.AuthMiddleware(
		&handlers.ArticleHandler{
			MongoDB:    s.mongoClient,
			PostgresDB: s.postgresClient,
		},
		s.authClient,
	)

	protectedCommentHandler := handlers.AuthMiddleware(
		&handlers.CommentHandler{
			MongoDB:    s.mongoClient,
			PostgresDB: s.postgresClient,
		},
		s.authClient,
	)

	s.mux.Handle("/article/{id}", protectedArticleHandler)
	s.mux.Handle("/article/{id}/", protectedArticleHandler)

	s.mux.Handle("/article", protectedArticleHandler)
	s.mux.Handle("/article/", protectedArticleHandler)

	s.mux.Handle("/comment", protectedCommentHandler)
	s.mux.Handle("/comment/", protectedCommentHandler)

	s.mux.Handle("/comment/{id}/like", protectedCommentHandler)
	s.mux.Handle("/comment/{id}/like/", protectedCommentHandler)

	s.mux.Handle("/auth/register", &handlers.AuthHandler{AuthClient: s.authClient})
	s.mux.Handle("/auth/register/", &handlers.AuthHandler{AuthClient: s.authClient})

	s.mux.Handle("/auth/login", &handlers.AuthHandler{AuthClient: s.authClient})
	s.mux.Handle("/auth/login/", &handlers.AuthHandler{AuthClient: s.authClient})

	blogHandler := &handlers.BlogHandler{Mongo: s.mongoClient}
	s.mux.Handle("/blog", blogHandler)
	s.mux.Handle("/blog/by-publisher", blogHandler)
}

func (s *Server) Start(addr string) error {
	// --- CORS Configuration ---
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Wrap existing server handler CORS middleware
	handlerWithCORS := c.Handler(s)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      handlerWithCORS,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	log.Printf("Starting blog service on %s", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
