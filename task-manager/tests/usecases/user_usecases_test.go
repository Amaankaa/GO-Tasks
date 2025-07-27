package usecases_test

import (
	"context"
	"errors"
	"testing"

	domain "task-manager/Domain"
	usecases "task-manager/Usecases"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StubRepo simulates UserRepository behaviors for testing
type StubRepo struct {
	OnRegister      func(domain.User) (domain.User, error)
	OnLogin         func(domain.User) (domain.LoginResponse, error)
	OnPromote       func(string) (domain.User, error)
	OnFindByUsername func(string) (domain.User, error)
}

func (r *StubRepo) RegisterUser(u domain.User) (domain.User, error) {
	return r.OnRegister(u)
}
func (r *StubRepo) LoginUser(u domain.User) (domain.LoginResponse, error) {
	return r.OnLogin(u)
}
func (r *StubRepo) PromoteUser(id string) (domain.User, error) {
	return r.OnPromote(id)
}
func (r *StubRepo) GetUserByUsername(username string) (domain.User, error) {
	return r.OnFindByUsername(username)
}

// UserUseCaseSuite is the testing suite for user-related use cases
type UserUseCaseSuite struct {
	suite.Suite
	repo    *StubRepo
	service *usecases.UserUsecase
	ctx     context.Context
}

func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, &UserUseCaseSuite{})
}

func (s *UserUseCaseSuite) SetupTest() {
	s.repo = &StubRepo{}
	s.service = usecases.NewUserUsecase(s.repo)
	s.ctx = context.TODO()
}

func (s *UserUseCaseSuite) TestRegisterUser() {
	s.Run("should register successfully", func() {
		s.SetupTest()
		input := domain.User{Username: "john", Password: "secure123"}
		mocked := input
		mocked.ID = primitive.NewObjectID()
		mocked.Role = "user"

		s.repo.OnRegister = func(u domain.User) (domain.User, error) {
			u.ID = mocked.ID
			u.Role = "user"
			return u, nil
		}

		res, err := s.service.RegisterUser(input)
		s.Require().NoError(err)
		s.Equal(mocked.Username, res.Username)
		s.Equal(mocked.ID, res.ID)
		s.Equal("user", res.Role)
	})

	s.Run("should fail when username is taken", func() {
		s.SetupTest()
		s.repo.OnRegister = func(u domain.User) (domain.User, error) {
			return domain.User{}, errors.New("username already taken")
		}
		_, err := s.service.RegisterUser(domain.User{Username: "john"})
		s.Error(err)
	})
}

func (s *UserUseCaseSuite) TestLoginUser() {
	s.Run("should authenticate with correct credentials", func() {
		s.SetupTest()
		input := domain.User{Username: "jane", Password: "pass123"}
		mockResp := domain.LoginResponse{
			ID:       primitive.NewObjectID(),
			Username: "jane",
			Token:    "valid-token",
		}

		s.repo.OnLogin = func(u domain.User) (domain.LoginResponse, error) {
			return mockResp, nil
		}

		token, err := s.service.LoginUser(input)
		s.Require().NoError(err)
		s.Equal(mockResp, token)
	})

	s.Run("should reject invalid login", func() {
		s.SetupTest()
		s.repo.OnLogin = func(u domain.User) (domain.LoginResponse, error) {
			return domain.LoginResponse{}, errors.New("invalid credentials")
		}
		_, err := s.service.LoginUser(domain.User{Username: "jane", Password: "wrong"})
		s.Error(err)
	})
}

func (s *UserUseCaseSuite) TestPromoteUser() {
	s.Run("should promote user to admin", func() {
		s.SetupTest()
		userID := primitive.NewObjectID().Hex()
		mockUser := domain.User{
			ID:       primitive.NewObjectID(),
			Username: "moderator",
			Role:     "admin",
		}

		s.repo.OnPromote = func(id string) (domain.User, error) {
			return mockUser, nil
		}

		res, err := s.service.PromoteUser(userID)
		s.NoError(err)
		s.Equal("admin", res.Role)
		s.Equal(mockUser.Username, res.Username)
	})

	s.Run("should error if user not found", func() {
		s.SetupTest()
		s.repo.OnPromote = func(id string) (domain.User, error) {
			return domain.User{}, errors.New("not found")
		}
		_, err := s.service.PromoteUser("invalid-id")
		s.Error(err)
	})
}

func (s *UserUseCaseSuite) TestGetUserByUsername() {
	s.Run("should return user by username", func() {
		s.SetupTest()
		uname := "alex"
		expected := domain.User{
			ID:       primitive.NewObjectID(),
			Username: uname,
			Role:     "user",
		}

		s.repo.OnFindByUsername = func(name string) (domain.User, error) {
			s.Equal(uname, name)
			return expected, nil
		}

		u, err := s.service.GetUserByUsername(uname)
		s.NoError(err)
		s.Equal(expected.ID, u.ID)
	})

	s.Run("should fail for unknown user", func() {
		s.SetupTest()
		s.repo.OnFindByUsername = func(name string) (domain.User, error) {
			return domain.User{}, errors.New("user not found")
		}
		_, err := s.service.GetUserByUsername("ghost")
		s.Error(err)
	})
}