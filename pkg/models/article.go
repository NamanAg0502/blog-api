package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Content     string             `json:"content" bson:"content"`
	Author      string             `json:"author" bson:"author"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type GetArticleRequest struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}

type CreateArticleRequest struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Content     string `json:"content" bson:"content"`
	Author      string `json:"author" bson:"author"`
}

type UpdateArticleRequest struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Content     string `json:"content" bson:"content"`
	Author      string `json:"author" bson:"author"`
}

type DeleteArticleRequest struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}
