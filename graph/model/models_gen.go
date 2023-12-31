// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type ArticleInput struct {
	Title   string      `json:"title"`
	Tags    []*TagInput `json:"tags"`
	Content string      `json:"content"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TagInput struct {
	Tag string `json:"tag"`
}

type TokenInput struct {
	Token string `json:"token"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
