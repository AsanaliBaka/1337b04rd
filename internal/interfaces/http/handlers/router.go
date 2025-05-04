package handlers

import (
	"net/http"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/infrastructure/api"
	"1337b04rd/internal/interfaces/http/middleware"
	"1337b04rd/pkg/logger"
)

type Handler struct {
	postHandler    *PostHandler
	sessionManager domain.SessionManager
	avatarAPI      *api.RickAndMortyAPI
	logger         *logger.CustomLogger
}

func NewHandler(
	postService domain.PostServer,
	sessionManager domain.SessionManager,
	avatarAPI *api.RickAndMortyAPI,
	logger *logger.CustomLogger,
) *Handler {
	return &Handler{
		postHandler:    NewPostHandler(postService, logger),
		sessionManager: sessionManager,
		avatarAPI:      avatarAPI,
		logger:         logger,
	}
}

func (h *Handler) SetupRoutes() http.Handler {
	// Создаем основной роутер
	router := http.NewServeMux()
	mw := middleware.NewMiddleware(h.sessionManager, h.avatarAPI, h.logger)

	// 1. Обработка статических файлов (должен быть первым)
	fs := http.FileServer(http.Dir("web/static"))
	router.Handle("/static/*", http.StripPrefix("/static", fs))

	// 2. API роуты (группируем под префиксом /api)
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("GET /posts", mw.Chain(
		mw.Logging,
		mw.UserCheck,
	)(http.HandlerFunc(h.postHandler.GetAllPosts)).ServeHTTP)

	apiRouter.HandleFunc("POST /posts", mw.Chain(
		mw.Logging,
		mw.UserCheck,
	)(http.HandlerFunc(h.postHandler.CreatePostHandler)).ServeHTTP)

	apiRouter.HandleFunc("GET /posts/{id}", mw.Chain(
		mw.Logging,
	)(http.HandlerFunc(h.postHandler.GetById)).ServeHTTP)

	apiRouter.HandleFunc("POST /comments", mw.Chain(
		mw.Logging,
		mw.UserCheck,
	)(http.HandlerFunc(h.postHandler.CreateComment)).ServeHTTP)

	router.Handle("/api/", http.StripPrefix("/api", apiRouter))

	// 3. HTML роуты (должны быть последними)
	router.HandleFunc("GET /{$}", h.serveIndex) // Явно указываем завершающий слэш
	router.HandleFunc("GET /posts/{id}/{$}", h.servePostPage)

	// Применяем глобальные middleware
	return mw.Chain(
		mw.Logging,
		mw.Recovery,
	)(router)
}

func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	// Serve HTML template for index page
	// Implement template rendering
}

func (h *Handler) servePostPage(w http.ResponseWriter, r *http.Request) {
	// Serve HTML template for post page
	// Implement template rendering
}
