package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentScores struct {
	Id    primitive.ObjectID `json:"_" bson:"_id,omitempty"`
	Score int                `json:"score" bson:"score"`
	Name  string             `json:"name" bson:"name"`
	Rank  int                `json:"rank" bson:"-"`
}

type ScoreResponse struct {
	Results  []*DocumentScores `json:"results"`
	AroundMe []*DocumentScores `json:"aroundMe,omitempty"`
	Meta     Meta              `json:"meta"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
	PerPage     int `json:"per_page"`
	NextPage    int `json:"next_page"`
}

type RequestPostScore struct {
	Name  string `json:"name" binding:"required"`
	Score int    `json:"score" binding:"required,min=0"`
}
