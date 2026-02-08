package fly_test

import (
	"context"
	"fly-go/config"
	"fly-go/database"
	"fly-go/fly"
	"os"
	"testing"

	log "fly-go/logger"

	"go.mongodb.org/mongo-driver/bson"
)

func setupManagerTest() (*fly.TaskManager, error) {
	// Set environment variable for config
	err := os.Setenv("FLY_CONFIG", "D:\\code\\web-app\\fly-go\\config.yaml")
	if err != nil {
		return nil, err
	}

	// Initialize logger
	logger := log.DefaultLogger()

	// Initialize database connection with config
	db, err := database.NewMongoDB(config.DatabaseConfig{
		MongoHost:     "120.48.130.105",
		MongoPort:     8717,
		MongoDatabase: "fly_test",
		MongoUsername: "root",
		MongoPassword: "example",
	})
	if err != nil {
		return nil, err
	}

	// Create TaskManager
	tm := fly.NewTaskManager(db, logger)

	return tm, nil
}

// TestTaskManager_DumpDefaultTask_Success tests successful dump of default task
func TestTaskManager_DumpDefaultTask_Success(t *testing.T) {
	t.Run("successfully dump default task to database", func(t *testing.T) {
		tm, err := setupManagerTest()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}

		// Clean up any existing default task before testing
		ctx := context.Background()
		collection, err := tm.DB.Collection(tm.Collection)
		if err != nil {
			t.Fatal(err)
		}
		_, err = collection.DeleteMany(ctx, bson.M{"id": "000"})
		if err != nil {
			t.Fatal(err)
		}

		// Call DumpDefaultTask
		err = tm.DumpDefaultTask()
		if err != nil {
			t.Errorf("DumpDefaultTask() failed: %v", err)
		}

		// Verify the task was inserted into the database
		var result bson.M
		err = collection.FindOne(ctx, bson.M{"id": "000"}).Decode(&result)
		if err != nil {
			t.Errorf("Failed to find dumped task in database: %v", err)
		}

		// Verify task fields
		if result["id"] != "000" {
			t.Errorf("Expected id '000', got %v", result["id"])
		}
		if result["name"] != "default_task" {
			t.Errorf("Expected name 'default_task', got %v", result["name"])
		}
		if result["status"] != fly.StatusUnknown {
			t.Errorf("Expected status 'unknown', got %v", result["status"])
		}
		if result["msg"] != "This is a default task" {
			t.Errorf("Expected msg 'This is a default task', got %v", result["msg"])
		}

		// Verify trigger fields
		trigger, ok := result["trigger"].(bson.M)
		if !ok {
			t.Errorf("Expected trigger to be a bson.M, got %T", result["trigger"])
		}

		if trigger["period"] != int64(1) {
			t.Errorf("Expected trigger.period 1, got %v", trigger["period"])
		}
		if trigger["start_time"] != "00:00" {
			t.Errorf("Expected trigger.start_time '00:00', got %v", trigger["start_time"])
		}
		if trigger["end_time"] != "23:59" {
			t.Errorf("Expected trigger.end_time '23:59', got %v", trigger["end_time"])
		}
		if int(trigger["type"].(int32)) != fly.Interval {
			t.Errorf("Expected trigger.type 'interval', got %v", trigger["type"])
		}
		if trigger["enabled"] != true {
			t.Errorf("Expected trigger.enabled true, got %v", trigger["enabled"])
		}

		// Clean up
		_, err = collection.DeleteMany(ctx, bson.M{"id": "000"})
		if err != nil {
			t.Errorf("Failed to clean up test data: %v", err)
		}
	})

	t.Run("verify trigger configuration details", func(t *testing.T) {
		tm, err := setupManagerTest()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}

		ctx := context.Background()
		collection, err := tm.DB.Collection(tm.Collection)
		if err != nil {
			t.Fatal(err)
		}
		_, err = collection.DeleteMany(ctx, bson.M{"id": "000"})
		if err != nil {
			t.Fatal(err)
		}

		// Call DumpDefaultTask
		err = tm.DumpDefaultTask()
		if err != nil {
			t.Errorf("DumpDefaultTask() failed: %v", err)
		}

		// Verify trigger date ranges
		var result bson.M
		err = collection.FindOne(ctx, bson.M{"id": "000"}).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		trigger := result["trigger"].(bson.M)

		if trigger["start_at"] != "2026-01-01" {
			t.Errorf("Expected trigger.start_at '2026-01-01', got %v", trigger["start_at"])
		}
		if trigger["end_at"] != "2099-01-01" {
			t.Errorf("Expected trigger.end_at '2099-01-01', got %v", trigger["end_at"])
		}

		// Clean up
		_, err = collection.DeleteMany(ctx, bson.M{"id": "000"})
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestTaskManager_LoadTask_Success tests successful loading of tasks from database
func TestTaskManager_LoadTask_Success(t *testing.T) {
	t.Run("successfully load tasks from database", func(t *testing.T) {
		tm, err := setupManagerTest()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}

		ctx := context.Background()
		collection, err := tm.DB.Collection(tm.Collection)
		if err != nil {
			t.Fatal(err)
		}

		// Clean up existing tasks
		_, err = collection.DeleteMany(ctx, bson.M{"name": "default_task"})
		if err != nil {
			t.Fatal(err)
		}

		// Insert test task
		err = tm.DumpDefaultTask()
		if err != nil {
			t.Fatal(err)
		}

		// Load tasks
		err = tm.LoadTask()
		if err != nil {
			t.Errorf("LoadTask() failed: %v", err)
		}

		// Verify task was loaded
		if tm.Count != 1 {
			t.Errorf("Expected 1 task loaded, got %d", tm.Count)
		}

		if len(tm.Names) != 1 {
			t.Errorf("Expected 1 task name, got %d", len(tm.Names))
		}

		if tm.Names[0] != "default_task" {
			t.Errorf("Expected task name 'default_task', got %s", tm.Names[0])
		}

		// Verify task in map
		key := "id.000.default_task"
		runner, exists := tm.TM[key]
		if !exists {
			t.Errorf("Task not found in TM map with key %s", key)
		}

		if runner.ID != "000" {
			t.Errorf("Expected task ID '000', got %s", runner.ID)
		}

		if runner.Name != "default_task" {
			t.Errorf("Expected task name 'default_task', got %s", runner.Name)
		}

		// Verify DB and Logger are set
		if runner.DB == nil {
			t.Error("Expected DB to be set on loaded task")
		}

		if runner.Logger == nil {
			t.Error("Expected Logger to be set on loaded task")
		}

		// Clean up
		_, err = collection.DeleteMany(ctx, bson.M{"id": "000"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("load multiple tasks from database", func(t *testing.T) {
		tm, err := setupManagerTest()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}

		ctx := context.Background()
		collection, err := tm.DB.Collection(tm.Collection)
		if err != nil {
			t.Fatal(err)
		}

		// Clean up existing tasks
		_, err = collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			t.Fatal(err)
		}

		// Insert multiple test tasks
		tasks := []interface{}{
			bson.M{
				"id":      "test001",
				"name":    "default_task",
				"status":  "idle",
				"trigger": bson.M{"type": fly.Interval, "period": 1, "enabled": true},
			},
			bson.M{
				"id":      "test002",
				"name":    "default_task",
				"status":  "running",
				"trigger": bson.M{"type": fly.Interval, "period": 2, "enabled": true},
			},
			bson.M{
				"id":      "test003",
				"name":    "default_task",
				"status":  "success",
				"trigger": bson.M{"type": fly.Interval, "period": 3, "enabled": true},
			},
		}
		_, err = collection.InsertMany(ctx, tasks)
		if err != nil {
			t.Fatal(err)
		}

		// Load tasks
		err = tm.LoadTask()
		if err != nil {
			t.Errorf("LoadTask() failed: %v", err)
		}

		// Verify all tasks were loaded
		if tm.Count != 3 {
			t.Errorf("Expected 3 tasks loaded, got %d", tm.Count)
		}

		if len(tm.Names) != 3 {
			t.Errorf("Expected 3 task names, got %d", len(tm.Names))
		}

		// Clean up
		_, err = collection.DeleteMany(ctx, bson.M{"name": "default_task"})
		if err != nil {
			t.Fatal(err)
		}
	})
}
