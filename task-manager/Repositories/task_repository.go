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

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository() (*TaskRepository, error) {
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
	collection := db.Collection("tasks")

	return &TaskRepository{
		collection: collection,
	}, nil
}

func (tr *TaskRepository) GetAllTasks() ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := tr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []domain.Task
	for cur.Next(ctx) {
		var task domain.Task
		if err := cur.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tr *TaskRepository) GetTaskByID(id string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, errors.New("invalid id format")
	}

	var task domain.Task
	err = tr.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return domain.Task{}, errors.New("not found")
	}

	return task, err
}

func (tr *TaskRepository) CreateTask(task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task.ID = primitive.NewObjectID()
	_, err := tr.collection.InsertOne(ctx, task)
	return task, err
}

func (tr *TaskRepository) UpdateTask(id string, updated domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, errors.New("invalid id format")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"title":       updated.Title,
			"description": updated.Description,
			"due_date":    updated.DueDate,
			"status":      updated.Status,
		},
	}

	res, err := tr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Task{}, err
	}

	if res.MatchedCount == 0 {
		return domain.Task{}, errors.New("not found")
	}

	return tr.GetTaskByID(id)
}

func (tr *TaskRepository) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	res, err := tr.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}