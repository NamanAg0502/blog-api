package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/namanag0502/blog-api/pkg/db"
	"github.com/namanag0502/blog-api/pkg/handlers"
)

type Router struct {
	db *db.Handler
}

func NewRouter(h *db.Handler) *Router {
	return &Router{db: h}
}

func (h *Router) Router() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	articleCollection := h.db.Client.Database("blog-api").Collection("articles")
	articleHandler := handlers.NewArticleHandler(articleCollection)

	mux.Route("/api/v1/articles", func(r chi.Router) {
		r.Get("/", articleHandler.GetAllArticles)
		r.Get("/{id}", articleHandler.GetArticle)
		r.Post("/", articleHandler.CreateArticle)
		r.Put("/{id}", articleHandler.UpdateArticle)
		r.Delete("/{id}", articleHandler.DeleteArticle)
	})
	return mux
}
