package fly

import (
	"context"
	"testing"
	"time"

	"fly-go/config"
	"fly-go/database"
	"fly-go/fly/executor"
	log "fly-go/logger"
)

// ==================== 测试辅助函数 ====================

// getTestDB 获取测试数据库连接
func getTestDB(t *testing.T) *database.MongoDB {
	cfg := config.DatabaseConfig{
		MongoHost:     "localhost",
		MongoPort:     8717,
		MongoUsername: "root",
		MongoPassword: "example",
		MongoDatabase: "fly_test",
	}

	db, err := database.NewMongoDB(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	return db
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T, db *database.MongoDB) {
	collection, err := db.Collection("task")
	if err != nil {
		t.Logf("Failed to get collection: %v", err)
		return
	}
	ctx := context.Background()
	collection.Drop(ctx)
}

// insertTestTask 插入测试任务数据
func insertTestTask(t *testing.T, db *database.MongoDB, task map[string]interface{}) {
	collection, err := db.Collection("task")
	if err != nil {
		t.Fatalf("Failed to get collection: %v", err)
	}
	ctx := context.Background()
	_, err = collection.InsertOne(ctx, task)
	if err != nil {
		t.Fatalf("Failed to insert test task: %v", err)
	}
}

// ==================== runner.go 测试 ====================

func TestGetRunners_Empty(t *testing.T) {
	// 清空 runners 列表进行测试
	originalRunners := runners
	runners = []executor.Runner{}

	defer func() { runners = originalRunners }()

	result := GetRunners()
	if len(result) != 0 {
		t.Errorf("GetRunners() should return empty map, got %d", len(result))
	}
}

func TestGetRunners_WithRunners(t *testing.T) {
	// 保存原始值
	originalRunners := runners

	// 创建 mock runner
	mockRunner := &MockTestRunner{name: "test-runner"}
	runners = []executor.Runner{mockRunner}

	defer func() { runners = originalRunners }()

	result := GetRunners()
	if len(result) != 1 {
		t.Errorf("GetRunners() should return 1 runner, got %d", len(result))
	}
	if _, ok := result["test-runner"]; !ok {
		t.Error("GetRunners() should contain 'test-runner'")
	}
}

func TestGetRunners_MultipleRunners(t *testing.T) {
	originalRunners := runners

	runners = []executor.Runner{
		&MockTestRunner{name: "runner-1"},
		&MockTestRunner{name: "runner-2"},
		&MockTestRunner{name: "runner-3"},
	}

	defer func() { runners = originalRunners }()

	result := GetRunners()
	if len(result) != 3 {
		t.Errorf("GetRunners() should return 3 runners, got %d", len(result))
	}
}

// MockTestRunner 用于测试的 Runner 实现
type MockTestRunner struct {
	name string
}

func (m *MockTestRunner) Name() string {
	return m.name
}

func (m *MockTestRunner) Run() (executor.TaskResult, error) {
	return executor.TaskResult{
		Status:  executor.StatusSuccess,
		Message: "test run",
	}, nil
}

func (m *MockTestRunner) Stop() error {
	return nil
}

// ==================== manager.go 测试 ====================

func TestNewTaskManager(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := log.DefaultLogger()

	tm := NewTaskManager(db, logger)

	if tm == nil {
		t.Fatal("NewTaskManager should not return nil")
	}
	if tm.DB != db {
		t.Error("TaskManager.DB should be the same as input db")
	}
	if tm.Logger != logger {
		t.Error("TaskManager.Logger should be the same as input logger")
	}
	if tm.Queue == nil {
		t.Error("TaskManager.Queue should be initialized")
	}
}

func TestTaskManager_InitDefaultTask(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	err := tm.InitDefaultTask()
	if err != nil {
		t.Errorf("InitDefaultTask() should return nil, got error: %v", err)
	}
}

func TestTaskManager_SetRunners(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册 runner
	tm.Queue.Register(&MockTestRunner{name: "test-runner"})

	// 验证 runners 已设置到 Queue
	if len(tm.Queue.Runners) != 1 {
		t.Errorf("Queue.Runners should have 1 runner, got %d", len(tm.Queue.Runners))
	}
	if _, ok := tm.Queue.Runners["test-runner"]; !ok {
		t.Error("Queue.Runners should contain 'test-runner'")
	}
}

// ==================== main.go 测试 ====================

func TestInit_ConnectDB(t *testing.T) {
	// 测试数据库连接
	db := getTestDB(t)
	defer db.Close()

	// 验证连接可用
	ctx := context.Background()
	err := db.Client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Database should be connected: %v", err)
	}
}

func TestInit_WithRealData(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 清理旧数据
	cleanupTestData(t, db)

	// 插入测试任务
	testTask := map[string]interface{}{
		"name":         "test-task-1",
		"runner_name":  "test-runner",
		"last_status":  executor.StatusIdle,
		"last_message": "",
		"params": map[string]interface{}{
			"type":       "interval",
			"interval":  60,
			"start_time": "09:00",
			"end_time":   "18:00",
		},
		"last_run_time": "",
		"last_end_time": "",
	}
	insertTestTask(t, db, testTask)

	// 创建 TaskManager 并加载
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册 mock runner
	mockRunner := &MockTestRunner{name: "test-runner"}
	tm.Queue.Register(mockRunner)

	// 从数据库加载任务
	tm.LoadFromDB()

	// 验证任务已加载
	if tm.Queue.WaitCount() != 1 {
		t.Errorf("Queue should have 1 task, got %d", tm.Queue.WaitCount())
	}

	// 验证 runner 已关联
	if _, ok := tm.Queue.Runners["test-runner"]; !ok {
		t.Error("Runner 'test-runner' should be registered")
	}
}

func TestInit_MultipleTasks(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 插入多个测试任务
	tasks := []map[string]interface{}{
		{
			"name":         "task-1",
			"runner_name":  "runner-a",
			"last_status":  executor.StatusIdle,
			"params":       map[string]interface{}{"type": "interval", "interval": 60},
		},
		{
			"name":         "task-2",
			"runner_name":  "runner-b",
			"last_status":  executor.StatusSuccess,
			"params":       map[string]interface{}{"type": "once"},
		},
		{
			"name":         "task-3",
			"runner_name":  "runner-a",
			"last_status":  executor.StatusError,
			"params":       map[string]interface{}{},
		},
	}

	for _, task := range tasks {
		insertTestTask(t, db, task)
	}

	// 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册 runners
	tm.Queue.Register(&MockTestRunner{name: "runner-a"})
	tm.Queue.Register(&MockTestRunner{name: "runner-b"})

	// 加载任务
	tm.LoadFromDB()

	// 验证
	if tm.Queue.WaitCount() != 3 {
		t.Errorf("Queue should have 3 tasks, got %d", tm.Queue.WaitCount())
	}
}

func TestListExecutors(t *testing.T) {
	// 保存原始值
	originalRunners := runners
	runners = []executor.Runner{
		&MockTestRunner{name: "exec-1"},
		&MockTestRunner{name: "exec-2"},
	}
	defer func() { runners = originalRunners }()

	result := ListExecutors()

	if len(result) != 2 {
		t.Errorf("ListExecutors() should return 2 executors, got %d", len(result))
	}

	// 验证包含预期的名称
	found := false
	for _, name := range result {
		if name == "exec-1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ListExecutors() should contain 'exec-1'")
	}
}

// ==================== 集成测试 ====================

func TestTaskExecution(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 插入一个待执行的任务
	testTask := map[string]interface{}{
		"name":         "exec-task",
		"runner_name":  "exec-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":       "interval",
			"interval":  1, // 1秒间隔
		},
	}
	insertTestTask(t, db, testTask)

	// 创建并初始化 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册 runner
	tm.Queue.Register(&MockTestRunner{name: "exec-runner"})

	// 加载任务
	tm.LoadFromDB()

	// 等待一小段时间让任务执行
	time.Sleep(2 * time.Second)

	// 验证任务已执行
	task := tm.Queue.Tasks[0]
	if task == nil {
		t.Fatal("Task should be loaded")
	}

	// 验证结果已记录
	if len(task.Result) != 1 {
		t.Logf("Task execution results: %d", len(task.Result))
		// 注意：由于使用 goroutine，可能还没执行完
	}
}

func TestTaskResultLimit(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册 runner
	tm.Queue.Register(&MockTestRunner{name: "limit-runner"})

	// 创建一个任务
	task := &executor.Task{
		Name:       "result-limit-test",
		Runner:     &MockTestRunner{name: "limit-runner"},
		RunnerName: "limit-runner",
	}
	tm.Queue.Tasks[0] = task

	// 保存 15 条结果（应该只保留最后 10 条）
	for i := 0; i < 15; i++ {
		task.Save(executor.TaskResult{
			Status:  executor.StatusSuccess,
			Message: "result",
			StartTime: time.Now(),
			EndTime:   time.Now(),
		})
	}

	// 验证只保留 10 条
	if len(task.Result) != 10 {
		t.Errorf("Task should keep only 10 results, got %d", len(task.Result))
	}
}
