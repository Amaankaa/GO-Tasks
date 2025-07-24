package repositories

import (
	"context"
	"errors"
	"os"
	"task-manager/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection  *mongo.Collection
	jwtService  domain.JWTService
	passwordService domain.PasswordService
}

func NewUserRepository(jwtService domain.JWTService, passwordService domain.PasswordService) (*UserRepository, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database("taskdb")
	collection := db.Collection("users")

	return &UserRepository{
		collection:      collection,
		jwtService:      jwtService,
		passwordService: passwordService,
	}, nil
}

func (ur *UserRepository) RegisterUser(user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Username == "" || user.Password == "" {
		return domain.User{}, errors.New("fields cannot be empty")
	}

	var existing domain.User
	err := ur.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return domain.User{}, errors.New("username already taken")
	}
	if err != mongo.ErrNoDocuments {
		return domain.User{}, err
	}

	// Check if this is the first user (make admin if so)
	userCount, err := ur.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return domain.User{}, err
	}
	if userCount == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	// Hash the password before storing
	hashedPassword, err := ur.passwordService.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword

	user.ID = primitive.NewObjectID()

	_, err = ur.collection.InsertOne(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = ""
	return user, nil
}

func (ur *UserRepository) LoginUser(user domain.User) (domain.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Username == "" {
		return domain.LoginResponse{}, errors.New("username is a required field")
	}

	var existingUser domain.User
	err := ur.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err != nil {
		return domain.LoginResponse{}, errors.New("invalid username or password")
	}

	if err = ur.passwordService.ComparePassword(existingUser.Password, user.Password); err != nil {
		return domain.LoginResponse{}, errors.New("invalid username or password")
	}

	// Generate JWT with role
	jwtToken, err := ur.jwtService.GenerateToken(existingUser.ID.Hex(), existingUser.Username, existingUser.Role)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		ID:       existingUser.ID,
		Username: existingUser.Username,
		Token:    jwtToken,
	}, nil
}

func (ur *UserRepository) PromoteUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, errors.New("invalid user ID")
	}

	update := bson.M{"$set": bson.M{"role": "admin"}}
	res, err := ur.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return domain.User{}, err
	}

	if res.MatchedCount == 0 {
		return domain.User{}, errors.New("user not found")
	}

	var updatedUser domain.User
	err = ur.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)
	if err != nil {
		return domain.User{}, err
	}

	updatedUser.Password = ""
	return updatedUser, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := ur.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return domain.User{}, errors.New("user not found")
	}

	return user, err
}