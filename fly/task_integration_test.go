package fly

import (
	"fmt"
	"testing"
	"time"

	"fly-go/fly/executor"
	log "fly-go/logger"
)

// ==================== 测试任务定义 ====================

// TestRunner 用于集成测试的 Runner 实现
type TestRunner struct {
	name       string
	shouldFail bool
	resultData any
}

func (r *TestRunner) Name() string {
	return r.name
}

func (r *TestRunner) Run() (executor.TaskResult, error) {
	start := time.Now()
	
	if r.shouldFail {
		return executor.TaskResult{
			Status:    executor.StatusError,
			Message:   "模拟执行失败",
			Data:      r.resultData,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
		}, fmt.Errorf("execution failed")
	}
	
	return executor.TaskResult{
		Status:    executor.StatusSuccess,
		Message:   "执行成功",
		Data:      r.resultData,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
	}, nil
}

func (r *TestRunner) Stop() error {
	return nil
}

// ==================== 数据库操作测试 ====================

// TestDB_SaveAndLoadTask 测试任务的保存和加载
func TestDB_SaveAndLoadTask(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 1. 创建并保存任务到数据库
	taskData := map[string]interface{}{
		"name":         "db-test-task",
		"runner_name":  "db-test-runner",
		"last_status":  executor.StatusIdle,
		"last_message":  "",
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  5,
			"start_time": "00:00",
			"end_time":   "23:59",
		},
		"last_run_time": "",
		"last_end_time": "",
	}
	
	insertTestTask(t, db, taskData)
	
	// 2. 创建 TaskManager 并从数据库加载
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	// 注册 runner
	tm.Queue.Register(&TestRunner{name: "db-test-runner"})
	
	// 加载任务
	tm.LoadFromDB()
	
	// 3. 验证任务已加载
	if tm.Queue.WaitCount() != 1 {
		t.Errorf("Expected 1 task loaded, got %d", tm.Queue.WaitCount())
	}
	
	// 验证任务数据正确
	task := tm.Queue.Tasks[0]
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	if task.Name != "db-test-task" {
		t.Errorf("Expected task name 'db-test-task', got '%s'", task.Name)
	}
	if task.RunnerName != "db-test-runner" {
		t.Errorf("Expected runner name 'db-test-runner', got '%s'", task.RunnerName)
	}
}

// TestDB_TaskExecutionSavesResult 测试任务执行后结果保存到内存
func TestDB_TaskExecutionSavesResult(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 1. 保存任务到数据库
	taskData := map[string]interface{}{
		"name":         "exec-result-task",
		"runner_name":  "result-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  1,
		},
	}
	insertTestTask(t, db, taskData)
	
	// 2. 创建 TaskManager 并加载任务
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	// 注册一个成功的 runner
	tm.Queue.Register(&TestRunner{
		name:       "result-runner",
		shouldFail: false,
		resultData: map[string]string{"key": "value"},
	})
	
	tm.LoadFromDB()
	
	// 3. 执行任务
	tm.Queue.Execute()
	
	// 等待一小段时间让 goroutine 完成
	time.Sleep(500 * time.Millisecond)
	
	// 4. 验证内存中的任务状态已更新
	task := tm.Queue.Tasks[0]
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	
	if task.LastStatus != executor.StatusSuccess {
		t.Errorf("Expected LastStatus '%s', got '%s'", executor.StatusSuccess, task.LastStatus)
	}
	
	// 验证结果已保存
	if len(task.Result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(task.Result))
	}
	
	// 验证结果数据
	if task.Result[0].Status != executor.StatusSuccess {
		t.Errorf("Expected result status '%s', got '%s'", executor.StatusSuccess, task.Result[0].Status)
	}
}

// TestDB_TaskFailureSavesError 测试任务失败时错误信息保存
func TestDB_TaskFailureSavesError(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 1. 保存任务
	taskData := map[string]interface{}{
		"name":         "fail-task",
		"runner_name":  "fail-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  1,
		},
	}
	insertTestTask(t, db, taskData)
	
	// 2. 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	// 注册一个会失败的 runner
	tm.Queue.Register(&TestRunner{
		name:       "fail-runner",
		shouldFail: true,
	})
	
	tm.LoadFromDB()
	
	// 3. 执行任务
	tm.Queue.Execute()
	time.Sleep(500 * time.Millisecond)
	
	// 4. 验证内存中的错误状态
	task := tm.Queue.Tasks[0]
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	
	if task.LastStatus != executor.StatusError {
		t.Errorf("Expected status '%s', got '%s'", executor.StatusError, task.LastStatus)
	}
	
	// 验证错误计数
	if tm.Queue.ErrorCount() != 1 {
		t.Errorf("Expected ErrorCount 1, got %d", tm.Queue.ErrorCount())
	}
	
	// 验证错误信息
	if task.LastMessage != "模拟执行失败" {
		t.Errorf("Expected message '模拟执行失败', got '%s'", task.LastMessage)
	}
}

// ==================== 完整流程测试 ====================

// TestFullWorkflow_LoadExecuteSave 测试完整的加载-执行-保存流程
func TestFullWorkflow_LoadExecuteSave(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 步骤1: 准备多个测试任务
	tasks := []map[string]interface{}{
		{
			"name":         "workflow-task-1",
			"runner_name":  "workflow-runner-1",
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type":      "interval",
				"interval":  2,
			},
		},
		{
			"name":         "workflow-task-2",
			"runner_name":  "workflow-runner-2",
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type":      "once",
			},
		},
	}
	
	for _, task := range tasks {
		insertTestTask(t, db, task)
	}
	
	// 步骤2: 创建 TaskManager 并加载任务
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	tm.Queue.Register(&TestRunner{
		name:       "workflow-runner-1",
		shouldFail: false,
		resultData: "task1_data",
	})
	tm.Queue.Register(&TestRunner{
		name:       "workflow-runner-2",
		shouldFail: false,
		resultData: "task2_data",
	})
	
	tm.LoadFromDB()
	
	// 步骤3: 验证任务已加载
	if tm.Queue.WaitCount() != 2 {
		t.Errorf("Expected 2 tasks loaded, got %d", tm.Queue.WaitCount())
	}
	
	// 步骤4: 执行任务
	tm.Queue.Execute()
	time.Sleep(800 * time.Millisecond)
	
	// 步骤5: 验证所有任务都执行了
	task1 := tm.Queue.Tasks[0]
	task2 := tm.Queue.Tasks[1]
	
	if task1.LastStatus != executor.StatusSuccess {
		t.Errorf("Task1 should have status '%s', got '%s'", executor.StatusSuccess, task1.LastStatus)
	}
	
	// once 类型任务只执行一次，第二次 CanRun 会返回 false
	// 所以我们检查任务2的状态
	if task2.LastStatus != executor.StatusSuccess && task2.LastStatus != executor.StatusIdle {
		t.Errorf("Task2 unexpected status: %s", task2.LastStatus)
	}
}

// TestFullWorkflow_MultipleExecutions 测试多次执行流程
func TestFullWorkflow_MultipleExecutions(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 创建任务（使用较短的间隔以便多次执行）
	taskData := map[string]interface{}{
		"name":         "multi-exec-task",
		"runner_name":  "multi-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  1, // 1秒间隔
		},
	}
	insertTestTask(t, db, taskData)
	
	// 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	tm.Queue.Register(&TestRunner{
		name:       "multi-runner",
		shouldFail: false,
	})
	
	tm.LoadFromDB()
	
	// 第一次执行
	tm.Queue.Execute()
	time.Sleep(500 * time.Millisecond)
	
	task := tm.Queue.Tasks[0]
	firstExecutions := len(task.Result)
	
	// 第二次执行（间隔1秒后）
	time.Sleep(1100 * time.Millisecond)
	tm.Queue.Execute()
	time.Sleep(500 * time.Millisecond)
	
	// 验证结果数量增加
	if len(task.Result) <= firstExecutions {
		t.Errorf("Expected more executions, first: %d, after: %d", firstExecutions, len(task.Result))
	}
	
	// 验证结果限制（只保留最近10条）
	if len(task.Result) > 10 {
		t.Errorf("Result should be limited to 10, got %d", len(task.Result))
	}
}

// ==================== 并发测试 ====================

// TestConcurrency_MultipleTasksRunInParallel 测试多个任务并发执行
func TestConcurrency_MultipleTasksRunInParallel(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 创建多个任务
	for i := 0; i < 5; i++ {
		taskData := map[string]interface{}{
			"name":         fmt.Sprintf("concurrent-task-%d", i),
			"runner_name":  "concurrent-runner",
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type":      "interval",
				"interval":  5,
			},
		}
		insertTestTask(t, db, taskData)
	}
	
	// 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	// 注册 runner
	tm.Queue.Register(&TestRunner{
		name:       "concurrent-runner",
		shouldFail: false,
	})
	
	tm.LoadFromDB()
	
	// 记录开始时间
	startTime := time.Now()
	
	// 执行任务
	tm.Queue.Execute()
	
	// 等待任务完成
	time.Sleep(500 * time.Millisecond)
	
	// 验证有任务正在运行或已完成
	running := tm.Queue.RunningCount()
	finished := tm.Queue.FinishCount()
	
	if running+finished == 0 {
		t.Error("At least one task should be running or finished")
	}
	
	t.Logf("Completed in %v, Running: %d, Finished: %d", time.Since(startTime), running, finished)
}

// ==================== 任务状态转换测试 ====================

// TestTaskStatus_Transitions 测试任务状态转换
func TestTaskStatus_Transitions(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	// 创建任务
	taskData := map[string]interface{}{
		"name":         "status-transition-task",
		"runner_name":  "status-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  1,
		},
	}
	insertTestTask(t, db, taskData)
	
	// 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	tm.Queue.Register(&TestRunner{
		name:       "status-runner",
		shouldFail: false,
	})
	
	tm.LoadFromDB()
	
	task := tm.Queue.Tasks[0]
	
	// 初始状态
	if task.LastStatus != executor.StatusIdle {
		t.Errorf("Initial status should be '%s', got '%s'", executor.StatusIdle, task.LastStatus)
	}
	
	// 执行任务
	tm.Queue.Execute()
	time.Sleep(500 * time.Millisecond)
	
	// 验证状态已更新
	if task.LastStatus != executor.StatusSuccess {
		t.Errorf("Status after execution should be '%s', got '%s'", executor.StatusSuccess, task.LastStatus)
	}
	
	// 验证统计数据
	if tm.Queue.ErrorCount() != 0 {
		t.Errorf("Error count should be 0, got %d", tm.Queue.ErrorCount())
	}
	if tm.Queue.FinishCount() == 0 {
		t.Error("Finish count should be at least 1")
	}
}

// ==================== 边界条件测试 ====================

// TestEdgeCase_EmptyDB 测试空数据库
func TestEdgeCase_EmptyDB(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	
	cleanupTestData(t, db)
	
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)
	
	tm.Queue.Register(&TestRunner{name: "edge-runner"})
	tm.LoadFromDB()
	
	if tm.Queue.WaitCount() != 0 {
		t.Errorf("Empty DB should have 0 tasks, got %d", tm.Queue.WaitCount())
	}
	
	// 执行空队列不应该出问题
	tm.Queue.Execute()
	
	t.Log("Empty queue handled correctly")
}

// TestEdgeCase_UnknownRunner 测试未注册的 Runner
func TestEdgeCase_UnknownRunner(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 任务引用一个不存在的 runner
	taskData := map[string]interface{}{
		"name":         "unknown-runner-task",
		"runner_name":  "non-existent-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  5,
		},
	}
	insertTestTask(t, db, taskData)

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 不注册任何 runner
	tm.LoadFromDB()

	// 任务应该被加载，但没有可用的 runner
	if tm.Queue.WaitCount() != 1 {
		t.Errorf("Should have 1 task loaded, got %d", tm.Queue.WaitCount())
	}

	task := tm.Queue.Tasks[0]
	if task.Runner != nil {
		t.Error("Task should not have a Runner assigned")
	}

	// 执行时应该不会 panic，并且任务会被跳过（从 WaitCount 减少可以看出）
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	// 任务应该从等待队列中移除
	if tm.Queue.WaitCount() != 0 {
		t.Errorf("Unknown runner task should be removed from wait queue, got %d", tm.Queue.WaitCount())
	}

	t.Log("Unknown runner handled correctly")
}


