package test_repositories

import (
	"context"
	"testing"

	domain "task-manager/Domain"
	repositories "task-manager/Repositories"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------------------------------------------------------
// Mock Implementations for Dependencies
// -------------------------------------------------------------------

type MockJWTService struct{}

func (f *MockJWTService) GenerateToken(id, role, username string) (string, error) {
	return "dummy-token", nil
}

func (f *MockJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"id":       "dummy-id",
		"username": "dummy-user",
		"role":     "standard",
	}, nil
}

type MockPasswordService struct{}

func (f *MockPasswordService) HashPassword(pw string) (string, error) {
	return "dummy-hash", nil
}

func (f *MockPasswordService) ComparePassword(hashed, plain string) error {
	return nil
}

// -------------------------------------------------------------------
// User Repository Test Suite
// -------------------------------------------------------------------

type AuthRepoTestSuite struct {
	suite.Suite
	db     *mongo.Database
	users  *mongo.Collection
	repo   domain.UserRepository
}

// Launches the test suite
func TestAuthRepoTestSuite(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("MongoDB client not available. Skipping user repo tests.")
	}
	suite.Run(t, new(AuthRepoTestSuite))
}

func (ts *AuthRepoTestSuite) SetupSuite() {
	ts.db = testMongoClient.Database("test_taskdb")
	ts.users = ts.db.Collection("users")
}

func (ts *AuthRepoTestSuite) SetupTest() {
	_, err := ts.users.DeleteMany(context.Background(), bson.D{})
	ts.Require().NoError(err)

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = ts.users.Indexes().CreateOne(context.Background(), index)
	ts.Require().NoError(err)

	jwtService := &MockJWTService{}
	passService := &MockPasswordService{}

	ts.repo = repositories.NewUserRepository(ts.users, jwtService, passService)
}

// -------------------------------------------------------------------
// Test Cases
// -------------------------------------------------------------------

func (ts *AuthRepoTestSuite) TestUserRegistration() {
	ts.Run("Should register new user successfully", func() {
		newUser := &domain.User{
			Username: "sampleuser",
			Password: "plaintext",
		}

		storedUser, err := ts.repo.RegisterUser(*newUser)

		ts.Require().NoError(err)
		ts.NotEmpty(storedUser.ID)
		ts.Equal("sampleuser", storedUser.Username)
	})

	ts.Run("Should reject duplicate usernames", func() {
		first := &domain.User{Username: "dupe", Password: "pw1"}
		_, err := ts.repo.RegisterUser(*first)
		ts.Require().NoError(err)

		second := &domain.User{Username: "dupe", Password: "pw2"}
		_, err = ts.repo.RegisterUser(*second)

		ts.Require().Error(err)
	})
}

func (ts *AuthRepoTestSuite) TestFindUserByUsername() {
	ts.Run("Should locate user by username", func() {
		targetUser := &domain.User{
			ID:       primitive.NewObjectID(),
			Username: "target",
			Password: "hash",
		}
		_, err := ts.users.InsertOne(context.Background(), targetUser)
		ts.Require().NoError(err)

		retrievedUser, err := ts.repo.GetUserByUsername("target")

		ts.Require().NoError(err)
		ts.Equal(targetUser.ID, retrievedUser.ID)
	})

	ts.Run("Should return error for missing user", func() {
		_, err := ts.repo.GetUserByUsername("ghostuser")

		ts.Require().Error(err)
	})
}
