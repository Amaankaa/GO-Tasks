package usecases_test

import (
	"errors"
	"testing"

	domain "task-manager/Domain"
	usecases "task-manager/Usecases"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -----------------------------------------------------------
// Fake implementation of domain.TaskRepository for testing
// -----------------------------------------------------------

type StubTaskRepo struct {
	OnCreate  func(domain.Task) (domain.Task, error)
	OnFind    func(string) (domain.Task, error)
	OnFetch   func() ([]domain.Task, error)
	OnUpdate  func(string, domain.Task) (domain.Task, error)
	OnRemove  func(string) error
}

func (s *StubTaskRepo) CreateTask(t domain.Task) (domain.Task, error) {
	if s.OnCreate != nil {
		return s.OnCreate(t)
	}
	return domain.Task{}, errors.New("CreateTask not implemented")
}

func (s *StubTaskRepo) GetTaskByID(id string) (domain.Task, error) {
	if s.OnFind != nil {
		return s.OnFind(id)
	}
	return domain.Task{}, errors.New("GetTaskByID not implemented")
}

func (s *StubTaskRepo) GetAllTasks() ([]domain.Task, error) {
	if s.OnFetch != nil {
		return s.OnFetch()
	}
	return nil, errors.New("GetAllTasks not implemented")
}

func (s *StubTaskRepo) UpdateTask(id string, t domain.Task) (domain.Task, error) {
	if s.OnUpdate != nil {
		return s.OnUpdate(id, t)
	}
	return domain.Task{}, errors.New("UpdateTask not implemented")
}

func (s *StubTaskRepo) DeleteTask(id string) error {
	if s.OnRemove != nil {
		return s.OnRemove(id)
	}
	return errors.New("DeleteTask not implemented")
}

// -----------------------------------------------------------
// Task Use Case Test Suite
// -----------------------------------------------------------

type TaskUseCaseSuite struct {
	suite.Suite
	mockStore *StubTaskRepo
	handler   *usecases.TaskUsecase
}

func TestTaskUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseSuite))
}

func (ts *TaskUseCaseSuite) SetupTest() {
	ts.mockStore = &StubTaskRepo{}
	ts.handler = usecases.NewTaskUsecase(ts.mockStore)
}

func (ts *TaskUseCaseSuite) TestCreateTask() {
	ts.Run("Success", func() {
		ts.SetupTest()

		incoming := domain.Task{
			Title:       "Prepare report",
			Description: "End of quarter summary",
			Status:      "pending",
		}
		result := incoming
		result.ID = primitive.NewObjectID()

		ts.mockStore.OnCreate = func(t domain.Task) (domain.Task, error) {
			t.ID = result.ID
			return t, nil
		}

		out, err := ts.handler.CreateTask(incoming)

		ts.Require().NoError(err)
		ts.Require().NotNil(out)
		ts.Equal(result.Title, out.Title)
		ts.Equal(result.ID, out.ID)
	})
}

func (ts *TaskUseCaseSuite) TestGetTaskByID() {
	ts.Run("Success", func() {
		ts.SetupTest()
		objID := primitive.NewObjectID()
		expected := domain.Task{ID: objID, Title: "Mock Task"}

		ts.mockStore.OnFind = func(id string) (domain.Task, error) {
			ts.Equal(objID.Hex(), id)
			return expected, nil
		}

		found, err := ts.handler.GetTaskByID(objID.Hex())

		ts.Require().NoError(err)
		ts.Require().NotNil(found)
		ts.Equal(expected.ID, found.ID)
	})
}

func (ts *TaskUseCaseSuite) TestGetAllTasks() {
	ts.Run("Success", func() {
		ts.SetupTest()
		mocked := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "One"},
			{ID: primitive.NewObjectID(), Title: "Two"},
		}

		ts.mockStore.OnFetch = func() ([]domain.Task, error) {
			return mocked, nil
		}

		results, err := ts.handler.GetAllTasks()

		ts.Require().NoError(err)
		ts.Len(results, 2)
		ts.Equal(mocked, results)
	})
}

func (ts *TaskUseCaseSuite) TestUpdateTask() {
	ts.Run("Success", func() {
		ts.SetupTest()
		oid := primitive.NewObjectID()
		oidStr := oid.Hex()

		updates := domain.Task{
			Title:       "Edited",
			Description: "Changes made",
			Status:      "done",
		}
		final := updates
		final.ID = oid

		ts.mockStore.OnUpdate = func(id string, in domain.Task) (domain.Task, error) {
			parsed, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return domain.Task{}, err
			}
			in.ID = parsed
			return in, nil
		}

		out, err := ts.handler.UpdateTask(oidStr, updates)

		ts.Require().NoError(err)
		ts.Require().NotNil(out)
		ts.Equal(final.Title, out.Title)
		ts.Equal(final.ID, out.ID)
	})
}

func (ts *TaskUseCaseSuite) TestDeleteTask() {
	ts.Run("Success", func() {
		ts.SetupTest()
		toRemove := primitive.NewObjectID()

		ts.mockStore.OnRemove = func(id string) error {
			ts.Equal(toRemove.Hex(), id)
			return nil
		}

		err := ts.handler.DeleteTask(toRemove.Hex())

		ts.Require().NoError(err)
	})
}