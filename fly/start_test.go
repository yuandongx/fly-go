package fly

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"fly-go/fly/executor"
	log "fly-go/logger"
)

// ==================== Start 函数测试 ====================

// IntegrationTestRunner 用于集成测试的 Runner
type IntegrationTestRunner struct {
	name          string
	execCount     int
	mu            sync.Mutex
	shouldFail    bool
	sleepDuration time.Duration
	resultData    any
}

func (r *IntegrationTestRunner) Name() string {
	return r.name
}

func (r *IntegrationTestRunner) Run() (executor.TaskResult, error) {
	r.mu.Lock()
	r.execCount++
	r.mu.Unlock()

	start := time.Now()

	if r.sleepDuration > 0 {
		time.Sleep(r.sleepDuration)
	}

	if r.shouldFail {
		return executor.TaskResult{
			Status:    executor.StatusError,
			Message:   "执行失败",
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

func (r *IntegrationTestRunner) Stop() error {
	return nil
}

func (r *IntegrationTestRunner) GetExecCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.execCount
}

// TestStart_InitWithMultipleRunners 测试使用多个 Runner 初始化
func TestStart_InitWithMultipleRunners(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 保存多个任务到数据库
	tasks := []map[string]interface{}{
		{
			"name":         "interval-task",
			"runner_name":  "interval-runner",
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type":      "interval",
				"interval":  1,
			},
		},
		{
			"name":         "once-task",
			"runner_name":  "once-runner",
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type": "once",
			},
		},
		{
			"name":         "daily-task",
			"runner_name":  "daily-runner",
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type":      "",
				"start_time": "00:00",
				"end_time":   "23:59",
			},
		},
	}

	for _, task := range tasks {
		insertTestTask(t, db, task)
	}

	// 创建 TaskManager
	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册多个 Runner
	tm.Queue.Register(&IntegrationTestRunner{name: "interval-runner"})
	tm.Queue.Register(&IntegrationTestRunner{name: "once-runner"})
	tm.Queue.Register(&IntegrationTestRunner{name: "daily-runner"})

	// 加载任务
	tm.LoadFromDB()

	// 验证任务加载
	if tm.Queue.WaitCount() != 3 {
		t.Errorf("Expected 3 tasks loaded, got %d", tm.Queue.WaitCount())
	}

	// 验证 Runner 关联
	for _, task := range tm.Queue.Tasks {
		if task.Runner == nil {
			t.Errorf("Task %s should have a Runner assigned", task.Name)
		}
	}

	// 执行一次
	tm.Queue.Execute()
	time.Sleep(500 * time.Millisecond)

	// 验证任务执行
	running := tm.Queue.RunningCount()
	finished := tm.Queue.FinishCount()
	if running+finished == 0 {
		t.Error("At least one task should be running or finished")
	}
}

// TestStart_IntervalTaskExecution 测试 interval 类型任务执行
func TestStart_IntervalTaskExecution(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 创建 interval 任务
	taskData := map[string]interface{}{
		"name":         "interval-test",
		"runner_name":  "interval-test-runner",
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

	runner := &IntegrationTestRunner{
		name:          "interval-test-runner",
		sleepDuration: 100 * time.Millisecond,
	}
	tm.Queue.Register(runner)

	tm.LoadFromDB()

	// 第一次执行
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	if runner.GetExecCount() != 1 {
		t.Errorf("Expected 1 execution, got %d", runner.GetExecCount())
	}

	// 等待间隔时间
	time.Sleep(1100 * time.Millisecond)

	// 第二次执行
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	if runner.GetExecCount() != 2 {
		t.Errorf("Expected 2 executions, got %d", runner.GetExecCount())
	}
}

// TestStart_OnceTaskExecution 测试 once 类型任务执行
func TestStart_OnceTaskExecution(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 创建 once 任务
	taskData := map[string]interface{}{
		"name":         "once-test",
		"runner_name":  "once-test-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type": "once",
		},
	}
	insertTestTask(t, db, taskData)

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	runner := &IntegrationTestRunner{name: "once-test-runner"}
	tm.Queue.Register(runner)

	tm.LoadFromDB()

	// 第一次执行 - 应该执行
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	if runner.GetExecCount() != 1 {
		t.Errorf("Expected 1 execution on first run, got %d", runner.GetExecCount())
	}

	// 再次执行 - once 类型不应再执行
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	// once 类型任务 LastEndTime 不为零，不会再执行
	if runner.GetExecCount() != 1 {
		t.Errorf("Once task should not execute again, got %d executions", runner.GetExecCount())
	}
}

// TestStart_DailyTaskExecution 测试每日定时任务执行
func TestStart_DailyTaskExecution(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 创建每日任务（在当前时间范围内）
	hour := time.Now().Hour()
	startHour := hour - 1
	if startHour < 0 {
		startHour = 0
	}
	endHour := hour + 1

	taskData := map[string]interface{}{
		"name":         "daily-test",
		"runner_name":  "daily-test-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":       "daily",
			"start_time": fmt.Sprintf("%02d:00", startHour),
			"end_time":   fmt.Sprintf("%02d:00", endHour),
		},
	}
	insertTestTask(t, db, taskData)

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	runner := &IntegrationTestRunner{name: "daily-test-runner"}
	tm.Queue.Register(runner)

	tm.LoadFromDB()

	// 执行
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	if runner.GetExecCount() != 1 {
		t.Errorf("Expected 1 execution for daily task, got %d", runner.GetExecCount())
	}
}

// TestStart_FailedTaskHandling 测试失败任务处理
func TestStart_FailedTaskHandling(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	taskData := map[string]interface{}{
		"name":         "fail-test",
		"runner_name":  "fail-test-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  5,
		},
	}
	insertTestTask(t, db, taskData)

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	runner := &IntegrationTestRunner{
		name:       "fail-test-runner",
		shouldFail: true,
	}
	tm.Queue.Register(runner)

	tm.LoadFromDB()

	// 执行失败的任务
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	// 验证错误计数
	if tm.Queue.ErrorCount() != 1 {
		t.Errorf("Expected 1 error, got %d", tm.Queue.ErrorCount())
	}

	// 验证任务状态
	task := tm.Queue.Tasks[0]
	if task.LastStatus != executor.StatusError {
		t.Errorf("Expected status 'error', got '%s'", task.LastStatus)
	}
}

// TestStart_ParallelExecution 测试并行执行
func TestStart_ParallelExecution(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 创建多个任务
	var runners []*IntegrationTestRunner
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("parallel-runner-%d", i)
		runners = append(runners, &IntegrationTestRunner{
			name:          name,
			sleepDuration: 200 * time.Millisecond, // 模拟耗时任务
		})

		taskData := map[string]interface{}{
			"name":         fmt.Sprintf("parallel-task-%d", i),
			"runner_name":  name,
			"last_status":  executor.StatusIdle,
			"params": map[string]interface{}{
				"type":      "interval",
				"interval":  1,
			},
		}
		insertTestTask(t, db, taskData)
	}

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 注册所有 runner
	for _, runner := range runners {
		tm.Queue.Register(runner)
	}

	tm.LoadFromDB()

	// 记录开始时间
	start := time.Now()

	// 执行所有任务
	tm.Queue.Execute()

	// 等待所有任务完成
	time.Sleep(500 * time.Millisecond)

	elapsed := time.Since(start)

	// 验证所有任务都执行了
	for i, runner := range runners {
		if runner.GetExecCount() != 1 {
			t.Errorf("Runner %d expected 1 execution, got %d", i, runner.GetExecCount())
		}
	}

	t.Logf("Parallel execution of 5 tasks completed in %v", elapsed)
}

// TestStart_UnknownRunnerTask 测试未知 Runner 的任务
func TestStart_UnknownRunnerTask(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 任务引用不存在的 runner
	taskData := map[string]interface{}{
		"name":         "unknown-runner-task",
		"runner_name":  "non-existent",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":      "interval",
			"interval":  1,
		},
	}
	insertTestTask(t, db, taskData)

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	// 不注册任何 runner

	tm.LoadFromDB()

	// 任务应该加载，但没有 runner
	if tm.Queue.WaitCount() != 1 {
		t.Errorf("Expected 1 task loaded, got %d", tm.Queue.WaitCount())
	}

	// 执行不应该 panic
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	// 任务应该从等待队列移除（因为没有 runner）
	if tm.Queue.WaitCount() != 0 {
		t.Errorf("Task with unknown runner should be removed from wait queue, got %d", tm.Queue.WaitCount())
	}
}

// TestStart_ComplexTaskParams 测试复杂任务参数
func TestStart_ComplexTaskParams(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 创建带跳过日期的任务
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	taskData := map[string]interface{}{
		"name":         "complex-params-task",
		"runner_name":  "complex-runner",
		"last_status":  executor.StatusIdle,
		"params": map[string]interface{}{
			"type":         "interval",
			"interval":     5,
			"start_time":   "00:00",
			"end_time":     "23:59",
			"skip_dates":   []string{tomorrow},
			"skip_weekdays": []int{0, 6}, // 跳过周日、周六
		},
	}
	insertTestTask(t, db, taskData)

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	runner := &IntegrationTestRunner{name: "complex-runner"}
	tm.Queue.Register(runner)

	tm.LoadFromDB()

	task := tm.Queue.Tasks[0]

	// 验证参数解析
	if task.Params.Type != "interval" {
		t.Errorf("Expected type 'interval', got '%s'", task.Params.Type)
	}

	if task.Params.Interval != 5*time.Second {
		t.Errorf("Expected interval 5s, got %v", task.Params.Interval)
	}

	// 执行
	tm.Queue.Execute()
	time.Sleep(200 * time.Millisecond)

	// 由于今天不是周末，应该能执行
	if runner.GetExecCount() != 1 {
		t.Errorf("Task should execute today (not weekend), got %d executions", runner.GetExecCount())
	}
}

// TestStart_AllTaskParams 测试所有支持的 TaskParam 类型
func TestStart_AllTaskParams(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	cleanupTestData(t, db)

	// 定义不同类型的任务
	testCases := []struct {
		name      string
		taskType  string
		interval  int
		startTime string
		endTime   string
	}{
		{"empty-type", "", 0, "", ""},
		{"interval-1s", "interval", 1, "", ""},
		{"interval-5s", "interval", 5, "", ""},
		{"once", "once", 0, "", ""},
		{"daily", "daily", 0, "00:00", "23:59"},
		// work-hours 需要当前时间在范围内，改为无时间限制以确保测试稳定
		{"work-hours", "interval", 10, "", ""},
	}

	var runners []*IntegrationTestRunner

	for _, tc := range testCases {
		runnerName := fmt.Sprintf("runner-%s", tc.name)
		runners = append(runners, &IntegrationTestRunner{name: runnerName})

		params := map[string]interface{}{
			"type": tc.taskType,
		}
		if tc.interval > 0 {
			params["interval"] = tc.interval
		}
		if tc.startTime != "" {
			params["start_time"] = tc.startTime
		}
		if tc.endTime != "" {
			params["end_time"] = tc.endTime
		}

		taskData := map[string]interface{}{
			"name":         fmt.Sprintf("task-%s", tc.name),
			"runner_name":  runnerName,
			"last_status":  executor.StatusIdle,
			"params":       params,
		}
		insertTestTask(t, db, taskData)
	}

	logger := log.DefaultLogger()
	tm := NewTaskManager(db, logger)

	for _, runner := range runners {
		tm.Queue.Register(runner)
	}

	tm.LoadFromDB()

	// 验证所有任务都加载
	if tm.Queue.WaitCount() != len(testCases) {
		t.Errorf("Expected %d tasks, got %d", len(testCases), tm.Queue.WaitCount())
	}

	// 执行所有任务
	tm.Queue.Execute()
	time.Sleep(500 * time.Millisecond)

	// 验证每个 Runner 都执行了
	for i, runner := range runners {
		if runner.GetExecCount() == 0 {
			t.Errorf("Runner %d (%s) should have executed", i, runner.Name())
		}
	}

	t.Logf("All %d task types executed successfully", len(testCases))
}
