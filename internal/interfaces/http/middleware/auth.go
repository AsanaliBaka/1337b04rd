package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/infrastructure/api"
	"1337b04rd/pkg/logger"
)

type Middleware struct {
	sessionManager  domain.SessionManager
	rickAndMortyAPI *api.RickAndMortyAPI
	logger          *logger.CustomLogger
}

func NewMiddleware(
	sm domain.SessionManager,
	rickAndMortyAPI *api.RickAndMortyAPI,
	logger *logger.CustomLogger,
) *Middleware {
	return &Middleware{
		sessionManager:  sm,
		rickAndMortyAPI: rickAndMortyAPI,
		logger:          logger,
	}
}

// UserCheck проверяет/создает сессию пользователя
func (m *Middleware) UserCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionID := getSessionIDFromRequest(r)
		if sessionID == "" {
			m.handleNewUser(w, r, next)
			return
		}

		user, err := m.sessionManager.GetSessionId(sessionID, ctx)
		if err != nil {
			m.logger.Info("Session not found, creating new", "session_id", sessionID)
			m.handleNewUser(w, r, next)
			return
		}

		// Обновляем время жизни сессии
		user.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // 1 неделя
		if err := m.sessionManager.CreateSession(ctx, user); err != nil {
			m.logger.Error("Failed to refresh session", "error", err)
		}

		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) handleNewUser(w http.ResponseWriter, r *http.Request, next http.Handler) {
	// Получаем случайного персонажа из Rick and Morty API
	avatarURL, name, err := m.rickAndMortyAPI.GetRandomCharacter(rand.Intn(800))
	if err != nil {
		m.logger.Error("Failed to get random character", "error", err)
		// Используем дефолтные значения если API не доступно
		avatarURL = ""
		name = "Anonymous"
	}

	newUser := &domain.UserRef{
		SessionID: generateNewSessionID(),
		Name:      name,
		AvatarURL: avatarURL,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 1 неделя
	}

	if err := m.sessionManager.CreateSession(r.Context(), newUser); err != nil {
		m.logger.Error("Failed to create session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setSessionIDToResponse(w, newUser.SessionID)
	ctx := context.WithValue(r.Context(), "user", newUser)
	next.ServeHTTP(w, r.WithContext(ctx))
}

// Logging middleware для логирования запросов
func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		m.logger.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr)

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		m.logger.Info("Request completed",
			"status", lrw.statusCode,
			"duration", time.Since(start))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Вспомогательные функции
func getSessionIDFromRequest(r *http.Request) string {
	if cookie, err := r.Cookie("session_id"); err == nil {
		return cookie.Value
	}
	if sessionID := r.Header.Get("X-Session-ID"); sessionID != "" {
		return sessionID
	}
	if sessionID := r.URL.Query().Get("session_id"); sessionID != "" {
		return sessionID
	}
	return ""
}

func setSessionIDToResponse(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	w.Header().Set("X-Session-ID", sessionID)
}

func generateNewSessionID() string {
	// В реальной реализации используйте github.com/google/uuid
	return "generated-unique-id-" + time.Now().Format("20060102150405")
}

// Chain создает цепочку middleware
func (m *Middleware) Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// Recovery middleware для обработки паник
func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("Recovered from panic",
					"error", err,
					"path", r.URL.Path)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
