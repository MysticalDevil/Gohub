package tests

import (
	"fmt"
	"testing"
	"time"

	"gohub/app/models/category"
	"gohub/app/models/link"
	"gohub/app/models/topic"
	"gohub/app/models/user"
	"gohub/pkg/database"
	"gohub/pkg/jwt"
)

type UserParams struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

func SeedUser(t *testing.T, params UserParams) user.User {
	t.Helper()

	if params.Name == "" {
		params.Name = fmt.Sprintf("user_%d", time.Now().UnixNano())
	}
	if params.Email == "" {
		params.Email = fmt.Sprintf("%s@testing.com", params.Name)
	}
	if params.Phone == "" {
		params.Phone = "00012345678"
	}
	if params.Password == "" {
		params.Password = "password123"
	}

	model := user.User{
		Name:     params.Name,
		Email:    params.Email,
		Phone:    params.Phone,
		Password: params.Password,
	}
	database.DB.Create(&model)
	return model
}

func IssueToken(userModel user.User) string {
	return jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
}

type CategoryParams struct {
	Name        string
	Description string
}

func SeedCategory(t *testing.T, params CategoryParams) category.Category {
	t.Helper()
	if params.Name == "" {
		params.Name = "cat"
	}
	if params.Description == "" {
		params.Description = "desc"
	}

	model := category.Category{
		Name:        params.Name,
		Description: params.Description,
	}
	database.DB.Create(&model)
	return model
}

type TopicParams struct {
	Title string
	Body  string
}

func SeedTopic(t *testing.T, userModel user.User, categoryModel category.Category, params TopicParams) topic.Topic {
	t.Helper()
	if params.Title == "" {
		params.Title = "hello"
	}
	if params.Body == "" {
		params.Body = "this is a body content"
	}

	model := topic.Topic{
		Title:      params.Title,
		Body:       params.Body,
		UserID:     userModel.GetStringID(),
		CategoryID: categoryModel.GetStringID(),
	}
	database.DB.Create(&model)
	return model
}

type LinkParams struct {
	Name string
	URL  string
}

func SeedLink(t *testing.T, params LinkParams) link.Link {
	t.Helper()
	if params.Name == "" {
		params.Name = "example"
	}
	if params.URL == "" {
		params.URL = "https://example.com"
	}

	model := link.Link{
		Name: params.Name,
		URL:  params.URL,
	}
	database.DB.Create(&model)
	return model
}
