package http

import (
	"database/sql"
	stdhttp "net/http"
	"time"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/config"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/http/handlers"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type Dependencies struct {
	Config config.Config
	DB     *sql.DB
}

func NewRouter(deps Dependencies) stdhttp.Handler {
	mux := stdhttp.NewServeMux()

	healthHandler := handlers.NewHealthHandler(deps.DB)
	userRepo := postgres.NewUserRepository(deps.DB)
	familyMemberRepo := postgres.NewFamilyMemberRepository(deps.DB)
	relationshipRepo := postgres.NewRelationshipRepository(deps.DB)
	storyRepo := postgres.NewStoryRepository(deps.DB)
	wisdomRepo := postgres.NewWisdomExtractRepository(deps.DB)
	userHandler := handlers.NewUserHandler(userRepo)
	familyMemberHandler := handlers.NewFamilyMemberHandler(familyMemberRepo)
	relationshipHandler := handlers.NewRelationshipHandler(relationshipRepo)
	storyHandler := handlers.NewStoryHandler(storyRepo)
	wisdomHandler := handlers.NewWisdomHandler(storyRepo, wisdomRepo)
	familyTreeHandler := handlers.NewFamilyTreeHandler(familyMemberRepo, relationshipRepo)

	authHandler := handlers.NewAuthHandler(userRepo)

	mux.HandleFunc("GET /health", healthHandler.Get)

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)

	protected := AuthMiddleware(deps.Config.JWTSecret)
	aiLimiter := NewRateLimiter(deps.Config.AIRequestsPerMinute, time.Minute)

	mux.Handle("GET /family/tree", protected(stdhttp.HandlerFunc(familyTreeHandler.Tree)))
	mux.Handle("POST /users", protected(stdhttp.HandlerFunc(userHandler.Create)))
	mux.Handle("GET /users", protected(stdhttp.HandlerFunc(userHandler.List)))
	mux.Handle("GET /users/{id}", protected(stdhttp.HandlerFunc(userHandler.Get)))
	mux.Handle("PATCH /users/{id}", protected(stdhttp.HandlerFunc(userHandler.Update)))
	mux.Handle("DELETE /users/{id}", protected(stdhttp.HandlerFunc(userHandler.Delete)))
	mux.Handle("POST /family-members", protected(stdhttp.HandlerFunc(familyMemberHandler.Create)))
	mux.Handle("GET /family-members", protected(stdhttp.HandlerFunc(familyMemberHandler.List)))
	mux.Handle("GET /family-members/{id}", protected(stdhttp.HandlerFunc(familyMemberHandler.Get)))
	mux.Handle("PATCH /family-members/{id}", protected(stdhttp.HandlerFunc(familyMemberHandler.Update)))
	mux.Handle("DELETE /family-members/{id}", protected(stdhttp.HandlerFunc(familyMemberHandler.Delete)))
	mux.Handle("POST /relationships", protected(stdhttp.HandlerFunc(relationshipHandler.Create)))
	mux.Handle("GET /relationships", protected(stdhttp.HandlerFunc(relationshipHandler.List)))
	mux.Handle("GET /relationships/{id}", protected(stdhttp.HandlerFunc(relationshipHandler.Get)))
	mux.Handle("PATCH /relationships/{id}", protected(stdhttp.HandlerFunc(relationshipHandler.Update)))
	mux.Handle("DELETE /relationships/{id}", protected(stdhttp.HandlerFunc(relationshipHandler.Delete)))
	mux.Handle("POST /stories", protected(stdhttp.HandlerFunc(storyHandler.Create)))
	mux.Handle("GET /stories", protected(stdhttp.HandlerFunc(storyHandler.List)))
	mux.Handle("GET /stories/{id}", protected(stdhttp.HandlerFunc(storyHandler.Get)))
	mux.Handle("PATCH /stories/{id}", protected(stdhttp.HandlerFunc(storyHandler.Update)))
	mux.Handle("DELETE /stories/{id}", protected(stdhttp.HandlerFunc(storyHandler.Delete)))
	mux.Handle("POST /stories/{id}/process-wisdom", protected(aiLimiter.Middleware()(stdhttp.HandlerFunc(wisdomHandler.ProcessWisdom))))
	mux.Handle("GET /wisdom-extracts", protected(stdhttp.HandlerFunc(wisdomHandler.List)))

	handler := stdhttp.Handler(mux)
	handler = RequestLoggingMiddleware()(handler)
	handler = CORSMiddleware(deps.Config.AllowedOrigins)(handler)
	return handler
}
