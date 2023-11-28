package helpers

import (
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sockleblu/digital_garden_backend/graph/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = []byte("secret_key")

func MapTagsFromInput(tagsInput []*model.TagInput) []*model.Tag {
	var tags []*model.Tag
	for _, tagInput := range tagsInput {
		tags = append(tags, &model.Tag{
			Tag: tagInput.Tag,
		})
	}

	return tags
}

func CreateSlug(title string) string {
	slug := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(title)), " ", "_")

	return slug
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, err
}

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatal("Error while generating jwt token")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}

func GetUserIdByUsername(db *gorm.DB, username string) (int, error) {
	var user model.User

	db.Where("username = ?", username).First(&user)

	return user.ID, nil
}
