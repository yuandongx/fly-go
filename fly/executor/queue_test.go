package executor

import (
	"testing"
	"time"
)

// ==================== params.go 测试 ====================

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"09:00 转换为秒", "09:00", 9 * 3600},
		{"00:00 转换为秒", "00:00", 0},
		{"18:30 转换为秒", "18:30", 18*3600 + 30*60},
		{"23:59 转换为秒", "23:59", 23*3600 + 59*60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toInt(tt.input)
			if result != tt.expected {
				t.Errorf("toInt(%s) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToIntDay(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectedLen int
	}{
		{"单日期转换", []string{"2026-01-01"}, 1},
		{"多日期转换", []string{"2026-01-01", "2026-06-15"}, 2},
		{"无效日期被忽略", []string{"invalid"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toIntDay(tt.input)
			if len(result) != tt.expectedLen {
				t.Errorf("toIntDay(%v) length = %d, expected %d", tt.input, len(result), tt.expectedLen)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"标准日期格式", "2026-01-01", false},
		{"标准日期格式2", "2026-12-31", false},
		{"无效格式", "2026/01/01", true},
		{"无效格式2", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestToAnyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected []string
	}{
		{"全是字符串", []any{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"包含非字符串", []any{"a", 1, "b", nil}, []string{"a", "b"}},
		{"空数组", []any{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toAnyStrings(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("toAnyStrings() = %v, expected %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("toAnyStrings()[%d] = %s, expected %s", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestToInts(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected []int
	}{
		{"int类型", []any{int(1), int(2), int(3)}, []int{1, 2, 3}},
		{"float64类型", []any{float64(1.0), float64(2.5), float64(3.9)}, []int{1, 2, 3}},
		{"混合类型", []any{int(1), float64(2.5), int(3)}, []int{1, 2, 3}},
		{"空数组", []any{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toInts(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("toInts() = %v, expected %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("toInts()[%d] = %d, expected %d", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestTaskParamInit(t *testing.T) {
	tests := []struct {
		name            string
		params          map[string]any
		expectStartTime int
		expectEndTime   int
		expectType      string
	}{
		{
			"完整参数初始化",
			map[string]any{
				"start_time":    "09:00",
				"end_time":      "18:00",
				"type":          "interval",
				"interval":      60,
				"start_date":    "2026-01-01",
				"end_date":      "2026-12-31",
				"skip_dates":    []any{"2026-06-01"},
				"skip_weekdays": []any{int(1)},
			},
			9 * 3600,
			18 * 3600,
			"interval",
		},
		{
			"空参数",
			map[string]any{},
			0,
			0,
			"",
		},
		{
			"部分参数",
			map[string]any{
				"start_time": "10:00",
				"type":       "once",
			},
			10 * 3600,
			0,
			"once",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TaskParam{}
			tp.Init(tt.params)

			if tp.StartTime != tt.expectStartTime {
				t.Errorf("StartTime = %d, expected %d", tp.StartTime, tt.expectStartTime)
			}
			if tp.EndTime != tt.expectEndTime {
				t.Errorf("EndTime = %d, expected %d", tp.EndTime, tt.expectEndTime)
			}
			if tp.Type != tt.expectType {
				t.Errorf("Type = %s, expected %s", tp.Type, tt.expectType)
			}
		})
	}
}

func TestTaskParamActive_DateRange(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		startDate    time.Time
		endDate      time.Time
		expectActive bool
	}{
		{"无日期限制", time.Time{}, time.Time{}, true},
		{"在开始日期之前", now.Add(24 * time.Hour), time.Time{}, false},
		{"在结束日期之后", time.Time{}, now.Add(-24 * time.Hour), false},
		{"在有效期内", now.Add(-24 * time.Hour), now.Add(24 * time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TaskParam{
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}
			result := tp.Active()
			if result != tt.expectActive {
				t.Errorf("Active() = %v, expected %v", result, tt.expectActive)
			}
		})
	}
}

func TestTaskParamActive_TimeRange(t *testing.T) {
	tests := []struct {
		name         string
		startTime    int
		endTime      int
		expectActive bool
	}{
		{"无时间限制", 0, 0, true},
		{"当前时间在范围内", 0, 86399, true}, // 全天
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TaskParam{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}
			result := tp.Active()
			if result != tt.expectActive {
				t.Errorf("Active() = %v, expected %v", result, tt.expectActive)
			}
		})
	}
}

// ==================== task.go 测试 ====================

// MockRunner 用于测试的 Runner 实现
type MockRunner struct {
	name     string
	runFunc  func() (TaskResult, error)
	stopFunc func() error
}

func (m *MockRunner) Name() string {
	return m.name
}

func (m *MockRunner) Run() (TaskResult, error) {
	return m.runFunc()
}

func (m *MockRunner) Stop() error {
	return m.stopFunc()
}

func TestTaskSave(t *testing.T) {
	t.Run("保存结果并更新状态", func(t *testing.T) {
		task := &Task{
			Name: "test-task",
		}
		result := TaskResult{
			Status:    StatusSuccess,
			Message:   "success",
			StartTime: time.Now().Add(-time.Second),
			EndTime:   time.Now(),
		}

		task.Save(result)

		if len(task.Result) != 1 {
			t.Errorf("Result length = %d, expected 1", len(task.Result))
		}
		if task.LastStatus != StatusSuccess {
			t.Errorf("LastStatus = %s, expected %s", task.LastStatus, StatusSuccess)
		}
	})

	t.Run("只保留最近10条结果", func(t *testing.T) {
		task := &Task{
			Name: "test-task",
		}

		// 添加15条结果
		for i := 0; i < 15; i++ {
			task.Save(TaskResult{
				Status:  StatusSuccess,
				Message: "result",
			})
		}

		if len(task.Result) != 10 {
			t.Errorf("Result length = %d, expected 10", len(task.Result))
		}
		// 第一条应该是第6条（索引5），因为前5条被移除了
	})
}

func TestTaskCanRun(t *testing.T) {
	t.Run("无参数限制可执行", func(t *testing.T) {
		task := &Task{
			Name:   "test-task",
			Params: TaskParam{},
		}
		if !task.CanRun() {
			t.Error("CanRun() should return true when no params limit")
		}
	})

	t.Run("once类型已执行过不可再次执行", func(t *testing.T) {
		task := &Task{
			Name: "test-task",
			Params: TaskParam{
				Type: TriggerOnce,
			},
			LastEndTime: time.Now().Add(-time.Minute), // 已执行过
		}
		if task.CanRun() {
			t.Error("CanRun() should return false for once type already executed")
		}
	})

	t.Run("interval类型间隔未到不可执行", func(t *testing.T) {
		task := &Task{
			Name: "test-task",
			Params: TaskParam{
				Type:     TriggerInterval,
				interval: 10 * time.Minute,
			},
			LastEndTime: time.Now().Add(-5 * time.Minute), // 只过了5分钟
		}
		if task.CanRun() {
			t.Error("CanRun() should return false when interval not reached")
		}
	})

	t.Run("interval类型间隔已到可以执行", func(t *testing.T) {
		task := &Task{
			Name: "test-task",
			Params: TaskParam{
				Type:     TriggerInterval,
				interval: 1 * time.Second, // 1秒间隔用于测试
			},
			LastEndTime: time.Now().Add(-2 * time.Second), // 过了2秒
		}
		if !task.CanRun() {
			t.Error("CanRun() should return true when interval reached")
		}
	})
}

func TestTaskRun(t *testing.T) {
	t.Run("正常执行Runner", func(t *testing.T) {
		expectedResult := TaskResult{
			Status:  StatusSuccess,
			Message: "executed",
		}
		runner := &MockRunner{
			name: "test-runner",
			runFunc: func() (TaskResult, error) {
				return expectedResult, nil
			},
		}

		task := &Task{
			Name:   "test-task",
			Runner: runner,
		}

		result, err := task.Run()
		if err != nil {
			t.Errorf("Run() error = %v", err)
		}
		if result.Status != expectedResult.Status {
			t.Errorf("Run() Status = %s, expected %s", result.Status, expectedResult.Status)
		}
	})

	t.Run("Runner返回错误", func(t *testing.T) {
		expectedErr := &MockError{msg: "test error"}
		runner := &MockRunner{
			name: "error-runner",
			runFunc: func() (TaskResult, error) {
				return TaskResult{}, expectedErr
			},
		}

		task := &Task{
			Name:   "test-task",
			Runner: runner,
		}

		_, err := task.Run()
		if err == nil {
			t.Error("Run() should return error")
		}
	})
}

func TestTaskStop(t *testing.T) {
	t.Run("正常停止Runner", func(t *testing.T) {
		stopCalled := false
		runner := &MockRunner{
			name: "test-runner",
			runFunc: func() (TaskResult, error) {
				return TaskResult{}, nil
			},
			stopFunc: func() error {
				stopCalled = true
				return nil
			},
		}

		task := &Task{
			Name:   "test-task",
			Runner: runner,
		}

		err := task.Stop()
		if err != nil {
			t.Errorf("Stop() error = %v", err)
		}
		if !stopCalled {
			t.Error("Stop() should call runner's Stop()")
		}
	})
}

func TestTaskIsTimeInRange(t *testing.T) {
	t.Run("无时间限制", func(t *testing.T) {
		task := &Task{
			Params: TaskParam{},
		}
		if !task.isTimeInRange(time.Now()) {
			t.Error("isTimeInRange() should return true when no time limit")
		}
	})

	t.Run("当前时间在范围内", func(t *testing.T) {
		now := time.Now()
		task := &Task{
			Params: TaskParam{
				StartTime: 0,
				EndTime:   86399, // 全天
			},
		}
		if !task.isTimeInRange(now) {
			t.Error("isTimeInRange() should return true for time within range")
		}
	})
}

// MockError 用于测试的错误类型
type MockError struct {
	msg string
}

func (e *MockError) Error() string {
	return e.msg
}

// ==================== queue.go 测试 ====================

func TestTaskQueueNew(t *testing.T) {
	t.Run("创建TaskQueue初始化所有字段", func(t *testing.T) {
		tq := New(nil)

		if tq == nil {
			t.Fatal("New() should not return nil")
		}
		if tq.Tasks == nil {
			t.Error("Tasks should be initialized")
		}
		if tq.Runners == nil {
			t.Error("runner should be initialized")
		}
		if tq.runningCount != 0 {
			t.Errorf("runningCount = %d, expected 0", tq.runningCount)
		}
		if tq.waitCount != 0 {
			t.Errorf("waitCount = %d, expected 0", tq.waitCount)
		}
		if tq.finishCount != 0 {
			t.Errorf("finishCount = %d, expected 0", tq.finishCount)
		}
		if tq.errorCount != 0 {
			t.Errorf("errorCount = %d, expected 0", tq.errorCount)
		}
	})
}

func TestTaskQueueRegister(t *testing.T) {
	t.Run("注册Runner", func(t *testing.T) {
		tq := New(nil)

		runner := &MockRunner{
			name: "test-runner",
			runFunc: func() (TaskResult, error) {
				return TaskResult{Status: StatusSuccess}, nil
			},
		}

		tq.Register(runner)

		if len(tq.Runners) != 1 {
			t.Errorf("runner count = %d, expected 1", len(tq.Runners))
		}
		if _, ok := tq.Runners["test-runner"]; !ok {
			t.Error("runner 'test-runner' should be registered")
		}
	})

	t.Run("重复注册覆盖", func(t *testing.T) {
		tq := New(nil)

		runner1 := &MockRunner{
			name: "test-runner",
			runFunc: func() (TaskResult, error) {
				return TaskResult{Status: StatusSuccess}, nil
			},
		}
		runner2 := &MockRunner{
			name: "test-runner",
			runFunc: func() (TaskResult, error) {
				return TaskResult{Status: StatusError}, nil
			},
		}

		tq.Register(runner1)
		tq.Register(runner2)

		if len(tq.Runners) != 1 {
			t.Errorf("runner count = %d, expected 1", len(tq.Runners))
		}
	})
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectEmpty bool
	}{
		{"RFC3339格式有效", "2026-01-01T10:00:00Z", false},
		{"无效格式返回空时间", "invalid", true},
		{"空字符串返回空时间", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTime(tt.input)
			if tt.expectEmpty && !result.IsZero() {
				t.Errorf("parseTime(%s) should return zero time for invalid input", tt.input)
			}
			if !tt.expectEmpty && result.IsZero() {
				t.Errorf("parseTime(%s) should not return zero time for valid input", tt.input)
			}
		})
	}
}
