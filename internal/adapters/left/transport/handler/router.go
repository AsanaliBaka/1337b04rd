package handler

import "net/http"

func 		SetupRoutes(postHandlers PostHandlers) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /catalog", postHandlers.GetPosts)
	router.HandleFunc("GET /create-post", postHandlers.CreatePost)
	router.HandleFunc("GET /post/{id}", postHandlers.GetPost)
	router.HandleFunc("POST /post/{id}/comment", postHandlers.AddComment)

	return router
}
