package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"context"
	"fmt"
	"log"

	"github.com/sockleblu/digital_garden_backend/graph/auth"
	"github.com/sockleblu/digital_garden_backend/graph/generated"
	"github.com/sockleblu/digital_garden_backend/graph/helpers"
	"github.com/sockleblu/digital_garden_backend/graph/model"
)

// User is the resolver for the user field.
func (r *articleResolver) User(ctx context.Context, obj *model.Article) (*model.User, error) {
	var user model.User

	r.DB.Where("id = ?", &obj.UserID).First(&user)

	return &user, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *articleResolver) CreatedAt(ctx context.Context, obj *model.Article) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.User, error) {
	var user model.User

	r.DB.Where("username = ?", input.Username).First(&user)

	//log.Fatal("input password is " + input.Password)
	valid, err := helpers.CheckPasswordHash(input.Password, user.Password)
	if !valid || err != nil {
		//return nil, &user.WrongUsernameOrPasswordError{}
		return nil, err
	}

	token, err := helpers.GenerateToken(input.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return &user, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	hash, err := helpers.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: input.Username,
		Password: hash,
		Email:    input.Email,
	}

	err = r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	token, err := helpers.GenerateToken(input.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return &user, nil
}

// ChangeUserPassword is the resolver for the changeUserPassword field.
func (r *mutationResolver) ChangeUserPassword(ctx context.Context, userID int, input model.UserInput) (*model.User, error) {
	user_auth := auth.ForContext(ctx)
	log.Printf("context looks like ", user_auth)
	if user_auth == nil {
		return &model.User{}, fmt.Errorf("Denying")
	}

	user := model.User{
		ID: userID,
	}

	err := r.DB.First(&user).Error
	if err != nil {
		return nil, err
	}

	hash, err := helpers.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hash

	err = r.DB.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUserByID is the resolver for the deleteUserByID field.
func (r *mutationResolver) DeleteUserByID(ctx context.Context, userID int) (bool, error) {
	user_auth := auth.ForContext(ctx)
	log.Printf("context looks like ", user_auth)
	if user_auth == nil {
		return false, fmt.Errorf("Denying")
	}

	user := model.User{
		ID: userID,
	}

	err := r.DB.Delete(&user).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateArticle is the resolver for the createArticle field.
func (r *mutationResolver) CreateArticle(ctx context.Context, input model.ArticleInput) (*model.Article, error) {
	user_auth := auth.ForContext(ctx)
	log.Printf("context looks like ", user_auth)
	if user_auth == nil {
		return &model.Article{}, fmt.Errorf("Denying")
	}

	slug := helpers.CreateSlug(input.Title)

	id, err := helpers.GetUserIdByUsername(r.DB, user.Username)
	if err != nil {
		return &model.Article{}, fmt.Errorf("Username %s was not found", user.Username)
	}

	article := model.Article{
		Slug:    slug,
		Title:   input.Title,
		Tags:    helpers.MapTagsFromInput(input.Tags),
		UserID:  id,
		Content: input.Content,
	}

	//r.DB.Model(&article).Association("Tags").Append(tags)

	creationErr := r.DB.Create(&article).Error
	if creationErr != nil {
		return nil, creationErr
	}

	return &article, nil
}

// UpdateArticle is the resolver for the updateArticle field.
func (r *mutationResolver) UpdateArticle(ctx context.Context, articleID int, input model.ArticleInput) (*model.Article, error) {
	user_auth := auth.ForContext(ctx)
	log.Printf("context looks like ", user_auth)
	if user_auth == nil {
		return &model.Article{}, fmt.Errorf("Denying")
	}

	article := model.Article{
		ID: articleID,
	}

	err := r.DB.First(&article).Error
	if err != nil {
		return nil, err
	}

	article.Title = input.Title
	article.Tags = helpers.MapTagsFromInput(input.Tags)
	article.Content = input.Content

	err = r.DB.Save(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// DeleteArticleByID is the resolver for the deleteArticleByID field.
func (r *mutationResolver) DeleteArticleByID(ctx context.Context, articleID int) (bool, error) {
	user_auth := auth.ForContext(ctx)
	log.Printf("context looks like ", user_auth)
	if user_auth == nil {
		return false, fmt.Errorf("Denying")
	}

	article := model.Article{
		ID: articleID,
	}

	err := r.DB.Delete(&article).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteArticleByTitle is the resolver for the deleteArticleByTitle field.
func (r *mutationResolver) DeleteArticleByTitle(ctx context.Context, title string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input model.TokenInput) (string, error) {
	username, err := helpers.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}

	token, err := helpers.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// AllUsers is the resolver for the allUsers field.
func (r *queryResolver) AllUsers(ctx context.Context) ([]*model.User, error) {
	user_auth := auth.ForContext(ctx)
	log.Printf("context looks like ", user_auth)
	if user_auth == nil {
		return []*model.User{}, fmt.Errorf("Denying")
	}

	var users []*model.User
	r.DB.Find(&users)

	return users, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, userID int) (*model.User, error) {
	var user model.User

	r.DB.Where("id = ?", userID).First(&user)

	return &user, nil
}

// Article is the resolver for the article field.
func (r *queryResolver) Article(ctx context.Context, slug string) (*model.Article, error) {
	var article model.Article
	r.DB.Where("slug = ?", slug).First(&article)

	return &article, nil
}

// AllArticles is the resolver for the allArticles field.
func (r *queryResolver) AllArticles(ctx context.Context) ([]*model.Article, error) {
	var articles []*model.Article
	r.DB.Preload("Tags").Find(&articles)

	return articles, nil
}

// ArticlesByTags is the resolver for the articlesByTags field.
func (r *queryResolver) ArticlesByTags(ctx context.Context, tagsInput []*model.TagInput) ([]*model.Article, error) {
	var tagList []string

	for _, tag := range tagsInput {
		tagList = append(tagList, tag.Tag)
	}

	var tags []*model.Tag

	// Load primary keys into Tag models
	r.DB.Debug().Where("tag IN ?", tagList).Find(&tags)

	var articles []*model.Article

	// Retrieve association using newly loaded tags model
	r.DB.Debug().Model(&tags).Association("Articles").Find(&articles)

	return articles, nil
}

// ArticleByID is the resolver for the articleById field.
func (r *queryResolver) ArticleByID(ctx context.Context, articleID int) (*model.Article, error) {
	var article model.Article
	r.DB.Where("id = ?", articleID).First(&article)

	return &article, nil
}

// Article returns generated.ArticleResolver implementation.
func (r *Resolver) Article() generated.ArticleResolver { return &articleResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type articleResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
