package fly_test

import (
	"fly-go/fly"
	"os"
	"testing"
)

func Setup() {
	// Initialize database connection
	// Configure settings
	err := os.Setenv("FLY_CONFIG", "D:\\code\\web-app\\fly-go\\config.yaml")
	if err != nil {
		panic(err)
	}
}

func TestInit(t *testing.T) {
	Setup()

	tests := []struct {
		name    string // description of this test case
		want    *fly.TaskManager
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Test successful initialization", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := fly.Init()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Init() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Init() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
