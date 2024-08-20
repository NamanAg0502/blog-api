package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/namanag0502/blog-api/pkg/models"
	"github.com/namanag0502/blog-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleHandler struct {
	C *mongo.Collection
}

func NewArticleHandler(c *mongo.Collection) *ArticleHandler {
	return &ArticleHandler{C: c}
}

func (s *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)

	if id == "" || err != nil {
		utils.WriteErrorResponse(w, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	var article models.Article
	filter := bson.D{{Key: "_id", Value: objID}}

	err = s.C.FindOne(r.Context(), filter).Decode(&article)
	if err == mongo.ErrNoDocuments {
		utils.WriteErrorResponse(w, fmt.Errorf("article not found"), http.StatusNotFound)
		return
	} else if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, article, http.StatusOK, "article retrieved successfully", nil)
}

func (s *ArticleHandler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	var articles []models.Article
	cur, err := s.C.Find(r.Context(), bson.D{})
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	defer cur.Close(r.Context())

	for cur.Next(r.Context()) {
		var article models.Article
		err := cur.Decode(&article)
		if err != nil {
			utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}

	utils.WriteJSONResponse(w, articles, http.StatusOK, "articles retrieved successfully", nil)
}

func (s *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req models.CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	article := models.Article{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Author:      req.Author,
		CreatedAt:   time.Now(),
	}

	result, err := s.C.InsertOne(r.Context(), article)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, result.InsertedID, http.StatusCreated, "article created successfully", nil)
}

func (s *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var req models.UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	updateData, err := bson.Marshal(req)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	var update bson.M
	if err := bson.Unmarshal(updateData, &update); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if len(update) == 0 {
		utils.WriteErrorResponse(w, nil, http.StatusBadRequest)
		return
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := s.C.UpdateOne(r.Context(), filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		utils.WriteErrorResponse(w, nil, http.StatusNotFound)
		return
	}

	utils.WriteJSONResponse(w, result, http.StatusOK, "article updated successfully", nil)
}

func (s *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.WriteErrorResponse(w, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	result, err := s.C.DeleteOne(r.Context(), filter)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		utils.WriteErrorResponse(w, nil, http.StatusNotFound)
		return
	}

	utils.WriteJSONResponse(w, result, http.StatusOK, "article deleted successfully", nil)
}
