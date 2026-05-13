package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

// MockHandler 用于测试的 ModeHandler 实现
type MockHandler struct {
	SpiderCalled bool
	ServerCalled bool
	ServerPort   int
	SpiderErr    error
	ServerErr    error
}

func (m *MockHandler) StartSpider() {
	m.SpiderCalled = true
}

func (m *MockHandler) StartServer(port int) {
	m.ServerCalled = true
	m.ServerPort = port
}

func TestRunMode_Spider(t *testing.T) {
	handler := &MockHandler{}
	
	RunMode("spider", 8080, handler)
	
	if !handler.SpiderCalled {
		t.Error("RunMode('spider') should call StartSpider")
	}
	if handler.ServerCalled {
		t.Error("RunMode('spider') should not call StartServer")
	}
}

func TestRunMode_Server(t *testing.T) {
	handler := &MockHandler{}
	
	RunMode("server", 9090, handler)
	
	if handler.SpiderCalled {
		t.Error("RunMode('server') should not call StartSpider")
	}
	if !handler.ServerCalled {
		t.Error("RunMode('server') should call StartServer")
	}
	if handler.ServerPort != 9090 {
		t.Errorf("RunMode('server') port = %d, want %d", handler.ServerPort, 9090)
	}
}

func TestRunMode_Default(t *testing.T) {
	handler := &MockHandler{}
	
	RunMode("unknown", 8080, handler)
	
	if !handler.SpiderCalled {
		t.Error("RunMode('unknown') should call StartSpider (default behavior)")
	}
}

func TestRunMode_ServerDefaultPort(t *testing.T) {
	handler := &MockHandler{}
	
	RunMode("server", 8080, handler)
	
	if handler.ServerPort != 8080 {
		t.Errorf("RunMode('server') default port = %d, want %d", handler.ServerPort, 8080)
	}
}

func TestRunMode_DifferentPorts(t *testing.T) {
	testCases := []struct {
		name     string
		mode     string
		port     int
		wantPort int
	}{
		{"server port 3000", "server", 3000, 3000},
		{"server port 8080", "server", 8080, 8080},
		{"server port 9000", "server", 9000, 9000},
		{"server port 65535", "server", 65535, 65535},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := &MockHandler{}
			RunMode(tc.mode, tc.port, handler)
			
			if handler.ServerPort != tc.wantPort {
				t.Errorf("port = %d, want %d", handler.ServerPort, tc.wantPort)
			}
		})
	}
}

func TestRunMode_EmptyMode(t *testing.T) {
	handler := &MockHandler{}
	RunMode("", 8080, handler)
	
	// 空字符串应该触发 default 分支
	if !handler.SpiderCalled {
		t.Error("RunMode('') should call StartSpider (default behavior)")
	}
}

// ErrorHandler 模拟返回错误的 Handler
type ErrorHandler struct {
	SpiderErr error
	ServerErr error
}

func (h *ErrorHandler) StartSpider() {
	if h.SpiderErr != nil {
		panic(h.SpiderErr)
	}
}

func (h *ErrorHandler) StartServer(port int) {
	if h.ServerErr != nil {
		panic(h.ServerErr)
	}
}

func TestRunMode_SpiderError(t *testing.T) {
	handler := &ErrorHandler{
		SpiderErr: errors.New("spider init failed"),
	}
	
	// 验证 spider 模式调用了 handler
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for spider error")
		}
	}()
	
	RunMode("spider", 8080, handler)
}

func TestRunMode_ServerError(t *testing.T) {
	handler := &ErrorHandler{
		ServerErr: errors.New("server init failed"),
	}
	
	// 验证 server 模式调用了 handler
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for server error")
		}
	}()
	
	RunMode("server", 8080, handler)
}

func TestModeHandler_Interface(t *testing.T) {
	// 验证 ModeHandler 接口正确实现
	var _ ModeHandler = &MockHandler{}
	var _ ModeHandler = &ErrorHandler{}
	var _ ModeHandler = &DefaultHandler{}
}

func TestRunMode_Modes(t *testing.T) {
	modes := []string{"spider", "server", "unknown", "", "Spider", "SERVER"}
	
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			handler := &MockHandler{}
			RunMode(mode, 8080, handler)
			// 只要不 panic 就算通过
		})
	}
}

// CaptureOutput 捕获标准输出
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	f()
	
	w.Close()
	os.Stdout = old
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestRunMode_OutputServer(t *testing.T) {
	// 测试 server 模式输出
	output := captureOutput(func() {
		handler := &MockHandler{}
		RunMode("server", 8080, handler)
	})
	
	if !strings.Contains(output, "Server is starting...") {
		t.Error("Expected 'Server is starting...' in output")
	}
	if !strings.Contains(output, "8080") {
		t.Error("Expected port 8080 in output")
	}
}

func TestRunMode_OutputSpider(t *testing.T) {
	// 测试 spider 模式输出
	output := captureOutput(func() {
		handler := &MockHandler{}
		RunMode("spider", 8080, handler)
	})
	
	if !strings.Contains(output, "Spider is starting...") {
		t.Error("Expected 'Spider is starting...' in output")
	}
}

func TestRunMode_OutputDefault(t *testing.T) {
	// 测试 default 模式输出
	output := captureOutput(func() {
		handler := &MockHandler{}
		RunMode("invalid", 8080, handler)
	})
	
	if !strings.Contains(output, "Only `spider` or `server` can be selected.") {
		t.Error("Expected warning message in output")
	}
}
