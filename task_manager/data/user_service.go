package data

import (
	"context"
	"errors"
	"task_manager/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// For development only. In production, use a secure secret management approach.
var jwtSecret = []byte("your_dev_secret_key")

func RegisterUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Username == "" {
		return models.User{}, errors.New("username is a required field")
	}

	var existing models.User
	err := userCol.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return models.User{}, errors.New("username already taken")
	}
	if err != mongo.ErrNoDocuments {
		return models.User{}, err
	}

	// Check if this is the first user (make admin if so)
	userCount, err := userCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return models.User{}, err
	}
	if userCount == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)

	user.ID = primitive.NewObjectID()

	_, err = userCol.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	user.Password = ""
	return user, nil
}

type LoginResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Token    string             `json:"token"`
}

func LoginUser(user models.User) (LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Username == "" {
		return LoginResponse{}, errors.New("username is a required field")
	}

	var existingUser models.User
	err := userCol.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err != nil {
		return LoginResponse{}, errors.New("invalid username or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return LoginResponse{}, errors.New("invalid username or password")
	}

	// Generate JWT with role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":      existingUser.ID.Hex(),
		"username": existingUser.Username,
		"role":     existingUser.Role,
	})

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:       existingUser.ID,
		Username: existingUser.Username,
		Token:    jwtToken,
	}, nil
}

func PromoteUser(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid user ID")
	}
	update := bson.M{"$set": bson.M{"role": "admin"}}
	res, err := userCol.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return models.User{}, err
	}
	if res.MatchedCount == 0 {
		return models.User{}, errors.New("user not found")
	}
	var updatedUser models.User
	err = userCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)
	if err != nil {
		return models.User{}, err
	}
	updatedUser.Password = ""
	return updatedUser, nil
}

func UserCol() *mongo.Collection {
	return userCol
}
